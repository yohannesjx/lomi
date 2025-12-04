-- ============================================================================
-- APP SETTINGS
-- User preferences for app language, theme, and cache
-- ============================================================================

-- Add app settings columns to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS app_language VARCHAR(10) DEFAULT 'en';
ALTER TABLE users ADD COLUMN IF NOT EXISTS app_theme VARCHAR(20) DEFAULT 'light';
ALTER TABLE users ADD COLUMN IF NOT EXISTS cache_cleared_at TIMESTAMP WITH TIME ZONE;

CREATE INDEX IF NOT EXISTS idx_users_app_language ON users(app_language);

COMMENT ON COLUMN users.app_language IS 'User preferred app language (en, am, om, etc)';
COMMENT ON COLUMN users.app_theme IS 'User preferred theme (light, dark, auto)';
COMMENT ON COLUMN users.cache_cleared_at IS 'Last time user cleared cache';
