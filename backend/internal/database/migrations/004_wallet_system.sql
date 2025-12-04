-- Wallet Management Schema
-- Production-grade database schema for TikTok-style app wallet system

-- ============================================
-- 1. WALLET TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    balance DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    total_earned DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    total_spent DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    total_withdrawn DECIMAL(15, 2) NOT NULL DEFAULT 0.00,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT positive_balance CHECK (balance >= 0),
    CONSTRAINT positive_earned CHECK (total_earned >= 0),
    CONSTRAINT positive_spent CHECK (total_spent >= 0),
    CONSTRAINT positive_withdrawn CHECK (total_withdrawn >= 0)
);

CREATE INDEX idx_wallets_user_id ON wallets(user_id);
CREATE INDEX idx_wallets_balance ON wallets(balance);

-- ============================================
-- 2. TRANSACTIONS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS wallet_transactions (
    id SERIAL PRIMARY KEY,
    wallet_id INTEGER NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_type VARCHAR(50) NOT NULL, -- 'credit', 'debit', 'purchase', 'gift_sent', 'gift_received', 'withdrawal', 'refund'
    amount DECIMAL(15, 2) NOT NULL,
    balance_before DECIMAL(15, 2) NOT NULL,
    balance_after DECIMAL(15, 2) NOT NULL,
    description TEXT,
    reference_id VARCHAR(100), -- External reference (payment gateway, gift ID, etc.)
    reference_type VARCHAR(50), -- 'gift', 'video', 'live_stream', 'purchase', 'withdrawal'
    status VARCHAR(20) NOT NULL DEFAULT 'completed', -- 'pending', 'completed', 'failed', 'cancelled'
    metadata JSONB, -- Additional data (payment method, recipient, etc.)
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT positive_amount CHECK (amount > 0),
    CONSTRAINT valid_status CHECK (status IN ('pending', 'completed', 'failed', 'cancelled'))
);

CREATE INDEX idx_transactions_wallet_id ON wallet_transactions(wallet_id);
CREATE INDEX idx_transactions_user_id ON wallet_transactions(user_id);
CREATE INDEX idx_transactions_type ON wallet_transactions(transaction_type);
CREATE INDEX idx_transactions_status ON wallet_transactions(status);
CREATE INDEX idx_transactions_created_at ON wallet_transactions(created_at DESC);
CREATE INDEX idx_transactions_reference ON wallet_transactions(reference_id, reference_type);

-- ============================================
-- 3. WITHDRAWAL REQUESTS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS withdrawal_requests (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wallet_id INTEGER NOT NULL REFERENCES wallets(id) ON DELETE CASCADE,
    amount DECIMAL(15, 2) NOT NULL,
    withdrawal_method VARCHAR(50) NOT NULL, -- 'bank_transfer', 'mobile_money', 'paypal', etc.
    account_details JSONB NOT NULL, -- Bank account, mobile number, etc.
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'rejected', 'cancelled'
    rejection_reason TEXT,
    processed_by INTEGER REFERENCES users(id), -- Admin who processed
    processed_at TIMESTAMP,
    transaction_id INTEGER REFERENCES wallet_transactions(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT positive_withdrawal CHECK (amount > 0),
    CONSTRAINT valid_withdrawal_status CHECK (status IN ('pending', 'processing', 'completed', 'rejected', 'cancelled'))
);

CREATE INDEX idx_withdrawals_user_id ON withdrawal_requests(user_id);
CREATE INDEX idx_withdrawals_status ON withdrawal_requests(status);
CREATE INDEX idx_withdrawals_created_at ON withdrawal_requests(created_at DESC);

-- ============================================
-- 4. PAYOUT METHODS TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS payout_methods (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    method_type VARCHAR(50) NOT NULL, -- 'bank_account', 'mobile_money', 'paypal'
    account_name VARCHAR(255) NOT NULL,
    account_details JSONB NOT NULL, -- Account number, bank name, mobile number, etc.
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payout_methods_user_id ON payout_methods(user_id);
CREATE INDEX idx_payout_methods_default ON payout_methods(user_id, is_default) WHERE is_default = TRUE;

-- ============================================
-- 5. COIN PACKAGES TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS coin_packages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    coins INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    bonus_coins INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT positive_coins CHECK (coins > 0),
    CONSTRAINT positive_price CHECK (price > 0)
);

CREATE INDEX idx_coin_packages_active ON coin_packages(is_active, display_order);

-- ============================================
-- 6. PURCHASE HISTORY TABLE
-- ============================================
CREATE TABLE IF NOT EXISTS coin_purchases (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    package_id INTEGER REFERENCES coin_packages(id),
    coins_purchased INTEGER NOT NULL,
    amount_paid DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    payment_method VARCHAR(50) NOT NULL, -- 'local_wallet', 'stripe', 'paypal', etc.
    payment_reference VARCHAR(255), -- External payment ID
    transaction_id INTEGER REFERENCES wallet_transactions(id),
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'completed', 'failed', 'refunded'
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    
    CONSTRAINT positive_coins_purchased CHECK (coins_purchased > 0),
    CONSTRAINT positive_amount_paid CHECK (amount_paid > 0)
);

CREATE INDEX idx_purchases_user_id ON coin_purchases(user_id);
CREATE INDEX idx_purchases_status ON coin_purchases(status);
CREATE INDEX idx_purchases_created_at ON coin_purchases(created_at DESC);

-- ============================================
-- 7. TRIGGERS FOR UPDATED_AT
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_wallets_updated_at BEFORE UPDATE ON wallets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_withdrawal_requests_updated_at BEFORE UPDATE ON withdrawal_requests
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payout_methods_updated_at BEFORE UPDATE ON payout_methods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_coin_packages_updated_at BEFORE UPDATE ON coin_packages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 8. SEED DATA - DEFAULT COIN PACKAGES
-- ============================================
INSERT INTO coin_packages (name, coins, price, currency, bonus_coins, display_order) VALUES
    ('Starter Pack', 100, 0.99, 'USD', 0, 1),
    ('Popular Pack', 500, 4.99, 'USD', 50, 2),
    ('Best Value', 1000, 9.99, 'USD', 150, 3),
    ('Premium Pack', 5000, 49.99, 'USD', 1000, 4),
    ('Ultimate Pack', 10000, 99.99, 'USD', 2500, 5)
ON CONFLICT DO NOTHING;

-- ============================================
-- 9. VIEWS FOR ANALYTICS
-- ============================================
CREATE OR REPLACE VIEW wallet_summary AS
SELECT 
    w.user_id,
    u.username,
    w.balance,
    w.total_earned,
    w.total_spent,
    w.total_withdrawn,
    COUNT(DISTINCT wt.id) as total_transactions,
    COUNT(DISTINCT CASE WHEN wt.transaction_type = 'credit' THEN wt.id END) as credit_count,
    COUNT(DISTINCT CASE WHEN wt.transaction_type = 'debit' THEN wt.id END) as debit_count,
    w.created_at as wallet_created_at
FROM wallets w
JOIN users u ON w.user_id = u.id
LEFT JOIN wallet_transactions wt ON w.id = wt.wallet_id
GROUP BY w.id, u.username;

-- ============================================
-- 10. FUNCTIONS FOR WALLET OPERATIONS
-- ============================================

-- Function to create wallet for new user
CREATE OR REPLACE FUNCTION create_user_wallet(p_user_id INTEGER)
RETURNS INTEGER AS $$
DECLARE
    v_wallet_id INTEGER;
BEGIN
    INSERT INTO wallets (user_id, balance, currency)
    VALUES (p_user_id, 0.00, 'USD')
    ON CONFLICT (user_id) DO NOTHING
    RETURNING id INTO v_wallet_id;
    
    RETURN v_wallet_id;
END;
$$ LANGUAGE plpgsql;

-- Function to get or create wallet
CREATE OR REPLACE FUNCTION get_or_create_wallet(p_user_id INTEGER)
RETURNS INTEGER AS $$
DECLARE
    v_wallet_id INTEGER;
BEGIN
    SELECT id INTO v_wallet_id FROM wallets WHERE user_id = p_user_id;
    
    IF v_wallet_id IS NULL THEN
        v_wallet_id := create_user_wallet(p_user_id);
    END IF;
    
    RETURN v_wallet_id;
END;
$$ LANGUAGE plpgsql;
