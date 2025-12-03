-- Add missing columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS username VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS profile_pic VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS profile_pic_small VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';
ALTER TABLE users ADD COLUMN IF NOT EXISTS device_token VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS auth_token VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS social VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS version VARCHAR(20);
ALTER TABLE users ADD COLUMN IF NOT EXISTS device VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS ip VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS country_id INTEGER;
ALTER TABLE users ADD COLUMN IF NOT EXISTS wallet DECIMAL(10, 2) DEFAULT 0.00;
ALTER TABLE users ADD COLUMN IF NOT EXISTS paypal VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS private INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS profile_view INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS reset_wallet_datetime TIMESTAMP WITH TIME ZONE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS referral_code VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS business INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS parent INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS comission_earned DECIMAL(10, 2) DEFAULT 0.00;
ALTER TABLE users ADD COLUMN IF NOT EXISTS followers_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS following_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS likes_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS video_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS block INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS sold_items_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS tagged_products_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS button VARCHAR(20) DEFAULT 'follow';
ALTER TABLE users ADD COLUMN IF NOT EXISTS notification INTEGER DEFAULT 1;
ALTER TABLE users ADD COLUMN IF NOT EXISTS unread_notification INTEGER DEFAULT 0;

-- Create privacy_settings table
CREATE TABLE IF NOT EXISTS privacy_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    videos_download INTEGER DEFAULT 1,
    direct_message INTEGER DEFAULT 1,
    duet INTEGER DEFAULT 1,
    liked_videos INTEGER DEFAULT 1,
    video_comment INTEGER DEFAULT 1,
    order_history INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Create push_notifications table
CREATE TABLE IF NOT EXISTS push_notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    likes INTEGER DEFAULT 1,
    comments INTEGER DEFAULT 1,
    new_followers INTEGER DEFAULT 1,
    mentions INTEGER DEFAULT 1,
    direct_messages INTEGER DEFAULT 1,
    video_updates INTEGER DEFAULT 1,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_privacy_settings_user_id ON privacy_settings(user_id);
CREATE INDEX IF NOT EXISTS idx_push_notifications_user_id ON push_notifications(user_id);

-- Function to update updated_at timestamp (ensure it exists)
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
DROP TRIGGER IF EXISTS update_privacy_settings_updated_at ON privacy_settings;
CREATE TRIGGER update_privacy_settings_updated_at BEFORE UPDATE ON privacy_settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_push_notifications_updated_at ON push_notifications;
CREATE TRIGGER update_push_notifications_updated_at BEFORE UPDATE ON push_notifications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
