-- ============================================================================
-- PROFILE & SETTINGS MIGRATION
-- Adds missing tables and columns for profile management
-- ============================================================================

-- Add missing columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS website VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_private BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_code VARCHAR(20) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS followers_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS following_count INTEGER DEFAULT 0;

-- Create index for referral code
CREATE INDEX IF NOT EXISTS idx_users_referral_code ON users(referral_code) WHERE referral_code IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_users_is_private ON users(is_private);

-- ============================================================================
-- FOLLOWS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS follows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX IF NOT EXISTS idx_follows_follower_id ON follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_follows_following_id ON follows(following_id);
CREATE INDEX IF NOT EXISTS idx_follows_created_at ON follows(created_at);

-- ============================================================================
-- PRIVACY SETTINGS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS privacy_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    account_privacy VARCHAR(20) DEFAULT 'public' CHECK (account_privacy IN ('public', 'private')),
    who_can_comment VARCHAR(20) DEFAULT 'everyone' CHECK (who_can_comment IN ('everyone', 'followers', 'nobody')),
    who_can_duet VARCHAR(20) DEFAULT 'everyone' CHECK (who_can_duet IN ('everyone', 'followers', 'nobody')),
    who_can_stitch VARCHAR(20) DEFAULT 'everyone' CHECK (who_can_stitch IN ('everyone', 'followers', 'nobody')),
    who_can_message VARCHAR(20) DEFAULT 'everyone' CHECK (who_can_message IN ('everyone', 'followers', 'nobody')),
    show_liked_videos BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_privacy_settings_user_id ON privacy_settings(user_id);

-- ============================================================================
-- NOTIFICATION SETTINGS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS notification_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    likes BOOLEAN DEFAULT TRUE,
    comments BOOLEAN DEFAULT TRUE,
    new_followers BOOLEAN DEFAULT TRUE,
    mentions BOOLEAN DEFAULT TRUE,
    live_streams BOOLEAN DEFAULT TRUE,
    direct_messages BOOLEAN DEFAULT TRUE,
    video_updates BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notification_settings_user_id ON notification_settings(user_id);

-- ============================================================================
-- REFERRALS TABLE
-- ============================================================================

CREATE TABLE IF NOT EXISTS referrals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    referrer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referred_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    referral_code VARCHAR(20) NOT NULL,
    reward_coins INTEGER DEFAULT 0,
    is_rewarded BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(referred_id)
);

CREATE INDEX IF NOT EXISTS idx_referrals_referrer_id ON referrals(referrer_id);
CREATE INDEX IF NOT EXISTS idx_referrals_referred_id ON referrals(referred_id);
CREATE INDEX IF NOT EXISTS idx_referrals_code ON referrals(referral_code);

-- ============================================================================
-- FUNCTIONS & TRIGGERS
-- ============================================================================

-- Function to update followers/following counts
CREATE OR REPLACE FUNCTION update_follow_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Increment follower count for the followed user
        UPDATE users SET followers_count = followers_count + 1 WHERE id = NEW.following_id;
        -- Increment following count for the follower
        UPDATE users SET following_count = following_count + 1 WHERE id = NEW.follower_id;
    ELSIF TG_OP = 'DELETE' THEN
        -- Decrement follower count for the unfollowed user
        UPDATE users SET followers_count = GREATEST(0, followers_count - 1) WHERE id = OLD.following_id;
        -- Decrement following count for the unfollower
        UPDATE users SET following_count = GREATEST(0, following_count - 1) WHERE id = OLD.follower_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Trigger for follow counts
DROP TRIGGER IF EXISTS trigger_update_follow_counts ON follows;
CREATE TRIGGER trigger_update_follow_counts
AFTER INSERT OR DELETE ON follows
FOR EACH ROW EXECUTE FUNCTION update_follow_counts();

-- Function to generate unique referral code
CREATE OR REPLACE FUNCTION generate_referral_code()
RETURNS TRIGGER AS $$
DECLARE
    new_code VARCHAR(20);
    code_exists BOOLEAN;
BEGIN
    IF NEW.referral_code IS NULL THEN
        LOOP
            -- Generate 8-character alphanumeric code
            new_code := UPPER(SUBSTRING(MD5(RANDOM()::TEXT || CLOCK_TIMESTAMP()::TEXT) FROM 1 FOR 8));
            
            -- Check if code already exists
            SELECT EXISTS(SELECT 1 FROM users WHERE referral_code = new_code) INTO code_exists;
            
            EXIT WHEN NOT code_exists;
        END LOOP;
        
        NEW.referral_code := new_code;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-generate referral code on user creation
DROP TRIGGER IF EXISTS trigger_generate_referral_code ON users;
CREATE TRIGGER trigger_generate_referral_code
BEFORE INSERT ON users
FOR EACH ROW EXECUTE FUNCTION generate_referral_code();

-- Apply updated_at triggers to new tables
DROP TRIGGER IF EXISTS update_privacy_settings_updated_at ON privacy_settings;
CREATE TRIGGER update_privacy_settings_updated_at 
BEFORE UPDATE ON privacy_settings 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_notification_settings_updated_at ON notification_settings;
CREATE TRIGGER update_notification_settings_updated_at 
BEFORE UPDATE ON notification_settings 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- DEFAULT SETTINGS FOR EXISTING USERS
-- ============================================================================

-- Insert default privacy settings for existing users
INSERT INTO privacy_settings (user_id)
SELECT id FROM users
WHERE NOT EXISTS (
    SELECT 1 FROM privacy_settings WHERE privacy_settings.user_id = users.id
);

-- Insert default notification settings for existing users
INSERT INTO notification_settings (user_id)
SELECT id FROM users
WHERE NOT EXISTS (
    SELECT 1 FROM notification_settings WHERE notification_settings.user_id = users.id
);

-- ============================================================================
-- COMMENTS
-- ============================================================================

COMMENT ON TABLE follows IS 'User follow relationships';
COMMENT ON TABLE privacy_settings IS 'User privacy preferences';
COMMENT ON TABLE notification_settings IS 'User notification preferences';
COMMENT ON TABLE referrals IS 'User referral tracking and rewards';
