-- ============================================================================
-- SOCIAL FEATURES
-- Profile shares tracking
-- ============================================================================

CREATE TABLE IF NOT EXISTS profile_shares (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    shared_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    platform VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_profile_shares_user_id ON profile_shares(user_id);
CREATE INDEX IF NOT EXISTS idx_profile_shares_shared_by ON profile_shares(shared_by);
CREATE INDEX IF NOT EXISTS idx_profile_shares_created_at ON profile_shares(created_at DESC);

COMMENT ON TABLE profile_shares IS 'Track when users share profiles on social platforms';
