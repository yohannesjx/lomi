-- Migration: Add daily free reveal tracking to users table
-- Run this migration to add the "Who Likes You" feature

ALTER TABLE users 
ADD COLUMN IF NOT EXISTS daily_free_reveal_used BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS last_reveal_date DATE;

-- Index for efficient queries
CREATE INDEX IF NOT EXISTS idx_users_last_reveal_date ON users(last_reveal_date);

COMMENT ON COLUMN users.daily_free_reveal_used IS 'Whether user has used their free daily reveal (resets at midnight Addis time)';
COMMENT ON COLUMN users.last_reveal_date IS 'Last date when user used a reveal (for daily reset logic)';

