-- Migration: Add location and tutorial fields
-- Run this on production database

BEGIN;

-- Add location columns if they don't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'latitude'
    ) THEN
        ALTER TABLE users ADD COLUMN latitude DECIMAL(10,8);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'longitude'
    ) THEN
        ALTER TABLE users ADD COLUMN longitude DECIMAL(11,8);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'has_seen_swipe_tutorial'
    ) THEN
        ALTER TABLE users ADD COLUMN has_seen_swipe_tutorial BOOLEAN DEFAULT FALSE;
    END IF;
END $$;

-- Create index for location-based queries
CREATE INDEX IF NOT EXISTS idx_users_location ON users(latitude, longitude) 
WHERE latitude IS NOT NULL AND longitude IS NOT NULL;

-- Create index for tutorial flag
CREATE INDEX IF NOT EXISTS idx_users_tutorial ON users(has_seen_swipe_tutorial) 
WHERE has_seen_swipe_tutorial = FALSE;

COMMIT;

-- Verify migration
SELECT 
    column_name, 
    data_type, 
    column_default,
    is_nullable
FROM information_schema.columns 
WHERE table_name = 'users' 
AND column_name IN ('latitude', 'longitude', 'has_seen_swipe_tutorial')
ORDER BY column_name;
