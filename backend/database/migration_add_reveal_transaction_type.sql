-- Migration: Add 'reveal' transaction type to existing enum
-- Run this if you have an existing database

-- Add new value to enum (PostgreSQL 9.1+)
ALTER TYPE transaction_type ADD VALUE IF NOT EXISTS 'reveal';

-- Note: If the above doesn't work (older PostgreSQL), you may need to:
-- 1. Create new enum with all values
-- 2. Alter table to use new enum
-- 3. Drop old enum
-- This is more complex and requires downtime, so IF NOT EXISTS should work for most cases.

