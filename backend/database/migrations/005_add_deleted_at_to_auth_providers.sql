-- Migration: Add deleted_at column to auth_providers (if missing)
-- Date: 2025-11-24
-- This fixes the error: column auth_providers.deleted_at does not exist

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'auth_providers' AND column_name = 'deleted_at'
    ) THEN
        ALTER TABLE auth_providers ADD COLUMN deleted_at TIMESTAMPTZ;
        CREATE INDEX IF NOT EXISTS idx_auth_providers_deleted_at
            ON auth_providers (deleted_at);
        RAISE NOTICE 'Added deleted_at column to auth_providers';
    ELSE
        RAISE NOTICE 'deleted_at column already exists in auth_providers';
    END IF;
END $$;

