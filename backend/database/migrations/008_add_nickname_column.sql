-- Add nickname column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS nickname VARCHAR(50);
