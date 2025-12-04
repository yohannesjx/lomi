package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONB is a custom type for PostgreSQL JSONB
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// ============================================
// WALLET MODELS
// ============================================

// Wallet represents a user's wallet
type Wallet struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	Balance        float64   `json:"balance" db:"balance"`
	TotalEarned    float64   `json:"total_earned" db:"total_earned"`
	TotalSpent     float64   `json:"total_spent" db:"total_spent"`
	TotalWithdrawn float64   `json:"total_withdrawn" db:"total_withdrawn"`
	Currency       string    `json:"currency" db:"currency"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// WalletTransaction represents a wallet transaction
type WalletTransaction struct {
	ID              int64     `json:"id" db:"id"`
	WalletID        int64     `json:"wallet_id" db:"wallet_id"`
	UserID          int64     `json:"user_id" db:"user_id"`
	TransactionType string    `json:"transaction_type" db:"transaction_type"`
	Amount          float64   `json:"amount" db:"amount"`
	BalanceBefore   float64   `json:"balance_before" db:"balance_before"`
	BalanceAfter    float64   `json:"balance_after" db:"balance_after"`
	Description     string    `json:"description" db:"description"`
	ReferenceID     *string   `json:"reference_id,omitempty" db:"reference_id"`
	ReferenceType   *string   `json:"reference_type,omitempty" db:"reference_type"`
	Status          string    `json:"status" db:"status"`
	Metadata        JSONB     `json:"metadata,omitempty" db:"metadata"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// TransactionType constants
const (
	TransactionTypeCredit       = "credit"
	TransactionTypeDebit        = "debit"
	TransactionTypePurchase     = "purchase"
	TransactionTypeGiftSent     = "gift_sent"
	TransactionTypeGiftReceived = "gift_received"
	TransactionTypeWithdrawal   = "withdrawal"
	TransactionTypeRefund       = "refund"
)

// TransactionStatus constants
const (
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
	TransactionStatusCancelled = "cancelled"
)

// WithdrawalRequest represents a withdrawal request
type WithdrawalRequest struct {
	ID               int64      `json:"id" db:"id"`
	UserID           int64      `json:"user_id" db:"user_id"`
	WalletID         int64      `json:"wallet_id" db:"wallet_id"`
	Amount           float64    `json:"amount" db:"amount"`
	WithdrawalMethod string     `json:"withdrawal_method" db:"withdrawal_method"`
	AccountDetails   JSONB      `json:"account_details" db:"account_details"`
	Status           string     `json:"status" db:"status"`
	RejectionReason  *string    `json:"rejection_reason,omitempty" db:"rejection_reason"`
	ProcessedBy      *int64     `json:"processed_by,omitempty" db:"processed_by"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty" db:"processed_at"`
	TransactionID    *int64     `json:"transaction_id,omitempty" db:"transaction_id"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

// WithdrawalStatus constants
const (
	WithdrawalStatusPending    = "pending"
	WithdrawalStatusProcessing = "processing"
	WithdrawalStatusCompleted  = "completed"
	WithdrawalStatusRejected   = "rejected"
	WithdrawalStatusCancelled  = "cancelled"
)

// PayoutMethod represents a user's payout method
type PayoutMethod struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	MethodType     string    `json:"method_type" db:"method_type"`
	AccountName    string    `json:"account_name" db:"account_name"`
	AccountDetails JSONB     `json:"account_details" db:"account_details"`
	IsDefault      bool      `json:"is_default" db:"is_default"`
	IsVerified     bool      `json:"is_verified" db:"is_verified"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// PayoutMethodType constants
const (
	PayoutMethodBankAccount = "bank_account"
	PayoutMethodMobileMoney = "mobile_money"
	PayoutMethodPayPal      = "paypal"
)

// CoinPackage represents a coin package for purchase
type CoinPackage struct {
	ID           int64     `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Coins        int       `json:"coins" db:"coins"`
	Price        float64   `json:"price" db:"price"`
	Currency     string    `json:"currency" db:"currency"`
	BonusCoins   int       `json:"bonus_coins" db:"bonus_coins"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	DisplayOrder int       `json:"display_order" db:"display_order"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CoinPurchase represents a coin purchase transaction
type CoinPurchase struct {
	ID               int64     `json:"id" db:"id"`
	UserID           int64     `json:"user_id" db:"user_id"`
	PackageID        *int64    `json:"package_id,omitempty" db:"package_id"`
	CoinsPurchased   int       `json:"coins_purchased" db:"coins_purchased"`
	AmountPaid       float64   `json:"amount_paid" db:"amount_paid"`
	Currency         string    `json:"currency" db:"currency"`
	PaymentMethod    string    `json:"payment_method" db:"payment_method"`
	PaymentReference *string   `json:"payment_reference,omitempty" db:"payment_reference"`
	TransactionID    *int64    `json:"transaction_id,omitempty" db:"transaction_id"`
	Status           string    `json:"status" db:"status"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// PurchaseStatus constants
const (
	PurchaseStatusPending   = "pending"
	PurchaseStatusCompleted = "completed"
	PurchaseStatusFailed    = "failed"
	PurchaseStatusRefunded  = "refunded"
)

// ============================================
// REQUEST/RESPONSE DTOs
// ============================================

// GetWalletBalanceResponse represents the wallet balance response
type GetWalletBalanceResponse struct {
	Balance        float64 `json:"balance"`
	TotalEarned    float64 `json:"total_earned"`
	TotalSpent     float64 `json:"total_spent"`
	TotalWithdrawn float64 `json:"total_withdrawn"`
	Currency       string  `json:"currency"`
}

// PurchaseCoinsRequest represents a coin purchase request
type PurchaseCoinsRequest struct {
	PackageID        *int64  `json:"package_id,omitempty"`
	Coins            int     `json:"coins" validate:"required,min=1"`
	Amount           float64 `json:"amount" validate:"required,min=0.01"`
	PaymentMethod    string  `json:"payment_method" validate:"required"`
	PaymentReference string  `json:"payment_reference,omitempty"`
}

// WithdrawRequest represents a withdrawal request
type WithdrawRequest struct {
	Amount           float64 `json:"amount" validate:"required,min=1"`
	WithdrawalMethod string  `json:"withdrawal_method" validate:"required"`
	AccountDetails   JSONB   `json:"account_details" validate:"required"`
}

// AddPayoutMethodRequest represents a request to add a payout method
type AddPayoutMethodRequest struct {
	MethodType     string `json:"method_type" validate:"required"`
	AccountName    string `json:"account_name" validate:"required"`
	AccountDetails JSONB  `json:"account_details" validate:"required"`
	IsDefault      bool   `json:"is_default"`
}

// TransactionHistoryResponse represents a transaction history item
type TransactionHistoryResponse struct {
	ID              int64     `json:"id"`
	TransactionType string    `json:"transaction_type"`
	Amount          float64   `json:"amount"`
	BalanceBefore   float64   `json:"balance_before"`
	BalanceAfter    float64   `json:"balance_after"`
	Description     string    `json:"description"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// WithdrawalHistoryResponse represents a withdrawal history item
type WithdrawalHistoryResponse struct {
	ID               int64      `json:"id"`
	Amount           float64    `json:"amount"`
	WithdrawalMethod string     `json:"withdrawal_method"`
	Status           string     `json:"status"`
	RejectionReason  *string    `json:"rejection_reason,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	ProcessedAt      *time.Time `json:"processed_at,omitempty"`
}
