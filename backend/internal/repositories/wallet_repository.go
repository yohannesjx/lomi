package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"lomi-backend/internal/models"

	"github.com/jmoiron/sqlx"
)

// WalletRepository handles wallet database operations
type WalletRepository struct {
	db *sqlx.DB
}

// NewWalletRepository creates a new wallet repository
func NewWalletRepository(db *sqlx.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

// ============================================
// WALLET OPERATIONS
// ============================================

// GetOrCreateWallet gets or creates a wallet for a user
func (r *WalletRepository) GetOrCreateWallet(ctx context.Context, userID string) (*models.Wallet, error) {
	var wallet models.Wallet

	// Try to get existing wallet
	err := r.db.GetContext(ctx, &wallet, `
		SELECT * FROM wallets WHERE user_id = $1
	`, userID)

	if err == nil {
		return &wallet, nil
	}

	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Create new wallet
	err = r.db.GetContext(ctx, &wallet, `
		INSERT INTO wallets (user_id, balance, currency)
		VALUES ($1, 0.00, 'USD')
		RETURNING *
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return &wallet, nil
}

// GetWalletByUserID gets a wallet by user ID
func (r *WalletRepository) GetWalletByUserID(ctx context.Context, userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.GetContext(ctx, &wallet, `
		SELECT * FROM wallets WHERE user_id = $1
	`, userID)

	if err == sql.ErrNoRows {
		return nil, errors.New("wallet not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

// UpdateWalletBalance updates wallet balance (use within transaction)
func (r *WalletRepository) UpdateWalletBalance(ctx context.Context, tx *sqlx.Tx, walletID string, newBalance float64) error {
	result, err := tx.ExecContext(ctx, `
		UPDATE wallets 
		SET balance = $1, updated_at = NOW()
		WHERE id = $2
	`, newBalance, walletID)

	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("wallet not found")
	}

	return nil
}

// IncrementTotalEarned increments total earned amount
func (r *WalletRepository) IncrementTotalEarned(ctx context.Context, tx *sqlx.Tx, walletID string, amount float64) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE wallets 
		SET total_earned = total_earned + $1, updated_at = NOW()
		WHERE id = $2
	`, amount, walletID)

	return err
}

// IncrementTotalSpent increments total spent amount
func (r *WalletRepository) IncrementTotalSpent(ctx context.Context, tx *sqlx.Tx, walletID string, amount float64) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE wallets 
		SET total_spent = total_spent + $1, updated_at = NOW()
		WHERE id = $2
	`, amount, walletID)

	return err
}

// IncrementTotalWithdrawn increments total withdrawn amount
func (r *WalletRepository) IncrementTotalWithdrawn(ctx context.Context, tx *sqlx.Tx, walletID string, amount float64) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE wallets 
		SET total_withdrawn = total_withdrawn + $1, updated_at = NOW()
		WHERE id = $2
	`, amount, walletID)

	return err
}

// ============================================
// TRANSACTION OPERATIONS
// ============================================

// CreateTransaction creates a new wallet transaction
func (r *WalletRepository) CreateTransaction(ctx context.Context, tx *sqlx.Tx, transaction *models.WalletTransaction) error {
	query := `
		INSERT INTO wallet_transactions (
			wallet_id, user_id, transaction_type, amount, 
			balance_before, balance_after, description, 
			reference_id, reference_type, status, metadata
		) VALUES (
			:wallet_id, :user_id, :transaction_type, :amount,
			:balance_before, :balance_after, :description,
			:reference_id, :reference_type, :status, :metadata
		) RETURNING id, created_at
	`

	rows, err := tx.NamedQuery(query, transaction)
	if err != nil {
		return fmt.Errorf("failed to create transaction: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&transaction.ID, &transaction.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan transaction: %w", err)
		}
	}

	return nil
}

// GetTransactionHistory gets transaction history for a user
func (r *WalletRepository) GetTransactionHistory(ctx context.Context, userID string, limit, offset int) ([]models.WalletTransaction, error) {
	var transactions []models.WalletTransaction

	err := r.db.SelectContext(ctx, &transactions, `
		SELECT * FROM wallet_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction history: %w", err)
	}

	return transactions, nil
}

// GetTransactionByID gets a transaction by ID
func (r *WalletRepository) GetTransactionByID(ctx context.Context, id string) (*models.WalletTransaction, error) {
	var transaction models.WalletTransaction

	err := r.db.GetContext(ctx, &transaction, `
		SELECT * FROM wallet_transactions WHERE id = $1
	`, id)

	if err == sql.ErrNoRows {
		return nil, errors.New("transaction not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &transaction, nil
}

// ============================================
// WITHDRAWAL OPERATIONS
// ============================================

// CreateWithdrawalRequest creates a new withdrawal request
func (r *WalletRepository) CreateWithdrawalRequest(ctx context.Context, tx *sqlx.Tx, request *models.WithdrawalRequest) error {
	query := `
		INSERT INTO withdrawal_requests (
			user_id, wallet_id, amount, withdrawal_method, 
			account_details, status
		) VALUES (
			:user_id, :wallet_id, :amount, :withdrawal_method,
			:account_details, :status
		) RETURNING id, created_at, updated_at
	`

	rows, err := tx.NamedQuery(query, request)
	if err != nil {
		return fmt.Errorf("failed to create withdrawal request: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&request.ID, &request.CreatedAt, &request.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan withdrawal request: %w", err)
		}
	}

	return nil
}

// GetWithdrawalHistory gets withdrawal history for a user
func (r *WalletRepository) GetWithdrawalHistory(ctx context.Context, userID string, limit, offset int) ([]models.WithdrawalRequest, error) {
	var requests []models.WithdrawalRequest

	err := r.db.SelectContext(ctx, &requests, `
		SELECT * FROM withdrawal_requests
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get withdrawal history: %w", err)
	}

	return requests, nil
}

// GetPendingWithdrawals gets pending withdrawal requests
func (r *WalletRepository) GetPendingWithdrawals(ctx context.Context, limit, offset int) ([]models.WithdrawalRequest, error) {
	var requests []models.WithdrawalRequest

	err := r.db.SelectContext(ctx, &requests, `
		SELECT * FROM withdrawal_requests
		WHERE status = $1
		ORDER BY created_at ASC
		LIMIT $2 OFFSET $3
	`, models.WithdrawalStatusPending, limit, offset)

	if err != nil {
		return nil, fmt.Errorf("failed to get pending withdrawals: %w", err)
	}

	return requests, nil
}

// UpdateWithdrawalStatus updates withdrawal request status
func (r *WalletRepository) UpdateWithdrawalStatus(ctx context.Context, tx *sqlx.Tx, id string, status string, processedBy *int64, rejectionReason *string) error {
	_, err := tx.ExecContext(ctx, `
		UPDATE withdrawal_requests
		SET status = $1, processed_by = $2, rejection_reason = $3, 
		    processed_at = NOW(), updated_at = NOW()
		WHERE id = $4
	`, status, processedBy, rejectionReason, id)

	return err
}

// ============================================
// PAYOUT METHOD OPERATIONS
// ============================================

// CreatePayoutMethod creates a new payout method
func (r *WalletRepository) CreatePayoutMethod(ctx context.Context, method *models.PayoutMethod) error {
	// If this is set as default, unset other defaults first
	if method.IsDefault {
		_, err := r.db.ExecContext(ctx, `
			UPDATE payout_methods SET is_default = FALSE WHERE user_id = $1
		`, method.UserID)
		if err != nil {
			return fmt.Errorf("failed to unset default payout methods: %w", err)
		}
	}

	query := `
		INSERT INTO payout_methods (
			user_id, method_type, account_name, account_details, is_default
		) VALUES (
			:user_id, :method_type, :account_name, :account_details, :is_default
		) RETURNING id, created_at, updated_at
	`

	rows, err := r.db.NamedQuery(query, method)
	if err != nil {
		return fmt.Errorf("failed to create payout method: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&method.ID, &method.CreatedAt, &method.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan payout method: %w", err)
		}
	}

	return nil
}

// GetPayoutMethods gets all payout methods for a user
func (r *WalletRepository) GetPayoutMethods(ctx context.Context, userID string) ([]models.PayoutMethod, error) {
	var methods []models.PayoutMethod

	err := r.db.SelectContext(ctx, &methods, `
		SELECT * FROM payout_methods
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`, userID)

	if err != nil {
		return nil, fmt.Errorf("failed to get payout methods: %w", err)
	}

	return methods, nil
}

// DeletePayoutMethod deletes a payout method
func (r *WalletRepository) DeletePayoutMethod(ctx context.Context, id, userID string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM payout_methods WHERE id = $1 AND user_id = $2
	`, id, userID)

	if err != nil {
		return fmt.Errorf("failed to delete payout method: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return errors.New("payout method not found")
	}

	return nil
}

// ============================================
// COIN PACKAGE OPERATIONS
// ============================================

// GetActiveCoinPackages gets all active coin packages
func (r *WalletRepository) GetActiveCoinPackages(ctx context.Context) ([]models.CoinPackage, error) {
	var packages []models.CoinPackage

	err := r.db.SelectContext(ctx, &packages, `
		SELECT * FROM coin_packages
		WHERE is_active = TRUE
		ORDER BY display_order ASC
	`)

	if err != nil {
		return nil, fmt.Errorf("failed to get coin packages: %w", err)
	}

	return packages, nil
}

// GetCoinPackageByID gets a coin package by ID
func (r *WalletRepository) GetCoinPackageByID(ctx context.Context, id string) (*models.CoinPackage, error) {
	var pkg models.CoinPackage

	err := r.db.GetContext(ctx, &pkg, `
		SELECT * FROM coin_packages WHERE id = $1 AND is_active = TRUE
	`, id)

	if err == sql.ErrNoRows {
		return nil, errors.New("coin package not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get coin package: %w", err)
	}

	return &pkg, nil
}

// CreateCoinPurchase creates a new coin purchase record
func (r *WalletRepository) CreateCoinPurchase(ctx context.Context, tx *sqlx.Tx, purchase *models.CoinPurchase) error {
	query := `
		INSERT INTO coin_purchases (
			user_id, package_id, coins_purchased, amount_paid,
			currency, payment_method, payment_reference, status
		) VALUES (
			:user_id, :package_id, :coins_purchased, :amount_paid,
			:currency, :payment_method, :payment_reference, :status
		) RETURNING id, created_at
	`

	rows, err := tx.NamedQuery(query, purchase)
	if err != nil {
		return fmt.Errorf("failed to create coin purchase: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&purchase.ID, &purchase.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to scan coin purchase: %w", err)
		}
	}

	return nil
}

// BeginTx starts a new database transaction
func (r *WalletRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}
