-- Gift System Migration
-- Adds wallet tracking and updates gift system for luxury virtual gifts

-- Add wallet tracking fields to users (if not exists)
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'total_spent'
    ) THEN
        ALTER TABLE users ADD COLUMN total_spent INTEGER DEFAULT 0 CHECK (total_spent >= 0);
    END IF;
    
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'total_earned'
    ) THEN
        ALTER TABLE users ADD COLUMN total_earned INTEGER DEFAULT 0 CHECK (total_earned >= 0);
    END IF;
END $$;

-- Add gift_type to gift_transactions (for easier querying)
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'gift_transactions' AND column_name = 'gift_type'
    ) THEN
        ALTER TABLE gift_transactions ADD COLUMN gift_type VARCHAR(50);
    END IF;
END $$;

-- Add coins field to cashout (for coin-based cashout)
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'payouts' AND column_name = 'coins'
    ) THEN
        ALTER TABLE payouts ADD COLUMN coins INTEGER;
    END IF;
END $$;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_gift_transactions_gift_type ON gift_transactions(gift_type);
CREATE INDEX IF NOT EXISTS idx_gift_transactions_receiver_created ON gift_transactions(receiver_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_payouts_coins ON payouts(coins) WHERE coins IS NOT NULL;

