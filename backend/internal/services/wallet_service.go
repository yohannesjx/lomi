package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lomi-backend/internal/models"
	"lomi-backend/internal/repositories"
)

// WalletService handles wallet business logic
type WalletService struct {
	walletRepo *repositories.WalletRepository
}

// NewWalletService creates a new wallet service
func NewWalletService(walletRepo *repositories.WalletRepository) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
	}
}

// ============================================
// WALLET OPERATIONS
// ============================================

// GetWalletBalance gets the wallet balance for a user
func (s *WalletService) GetWalletBalance(ctx context.Context, userID string) (*models.GetWalletBalanceResponse, error) {
	wallet, err := s.walletRepo.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &models.GetWalletBalanceResponse{
		Balance:        wallet.Balance,
		TotalEarned:    wallet.TotalEarned,
		TotalSpent:     wallet.TotalSpent,
		TotalWithdrawn: wallet.TotalWithdrawn,
		Currency:       wallet.Currency,
	}, nil
}

// ============================================
// COIN PURCHASE OPERATIONS
// ============================================

// PurchaseCoins handles coin purchase
func (s *WalletService) PurchaseCoins(ctx context.Context, userID string, req *models.PurchaseCoinsRequest) (*models.WalletTransaction, error) {
	// Validate request
	if req.Coins <= 0 {
		return nil, errors.New("coins must be greater than 0")
	}
	if req.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Get or create wallet
	wallet, err := s.walletRepo.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Start transaction
	tx, err := s.walletRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Create coin purchase record
	purchase := &models.CoinPurchase{
		UserID:           userID,
		PackageID:        req.PackageID,
		CoinsPurchased:   req.Coins,
		AmountPaid:       req.Amount,
		Currency:         "USD",
		PaymentMethod:    req.PaymentMethod,
		PaymentReference: &req.PaymentReference,
		Status:           models.PurchaseStatusCompleted,
	}

	err = s.walletRepo.CreateCoinPurchase(ctx, tx, purchase)
	if err != nil {
		return nil, fmt.Errorf("failed to create purchase record: %w", err)
	}

	// Calculate new balance
	coinsAsBalance := float64(req.Coins) // 1 coin = 1 unit in balance
	newBalance := wallet.Balance + coinsAsBalance

	// Update wallet balance
	err = s.walletRepo.UpdateWalletBalance(ctx, tx, wallet.ID, newBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Increment total earned
	err = s.walletRepo.IncrementTotalEarned(ctx, tx, wallet.ID, coinsAsBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to increment total earned: %w", err)
	}

	// Create transaction record
	transaction := &models.WalletTransaction{
		WalletID:        wallet.ID,
		UserID:          userID,
		TransactionType: models.TransactionTypeCredit,
		Amount:          coinsAsBalance,
		BalanceBefore:   wallet.Balance,
		BalanceAfter:    newBalance,
		Description:     fmt.Sprintf("Purchased %d coins", req.Coins),
		ReferenceID:     &req.PaymentReference,
		ReferenceType:   stringPtr("coin_purchase"),
		Status:          models.TransactionStatusCompleted,
		Metadata: models.JSONB{
			"coins":          req.Coins,
			"amount_paid":    req.Amount,
			"payment_method": req.PaymentMethod,
		},
	}

	err = s.walletRepo.CreateTransaction(ctx, tx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// GetCoinPackages gets all available coin packages
func (s *WalletService) GetCoinPackages(ctx context.Context) ([]models.CoinPackage, error) {
	return s.walletRepo.GetActiveCoinPackages(ctx)
}

// ============================================
// WITHDRAWAL OPERATIONS
// ============================================

// RequestWithdrawal creates a withdrawal request
func (s *WalletService) RequestWithdrawal(ctx context.Context, userID string, req *models.WithdrawRequest) (*models.WithdrawalRequest, error) {
	// Validate request
	if req.Amount <= 0 {
		return nil, errors.New("withdrawal amount must be greater than 0")
	}

	// Get wallet
	wallet, err := s.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Check if user has sufficient balance
	if wallet.Balance < req.Amount {
		return nil, errors.New("insufficient balance")
	}

	// Minimum withdrawal amount check (e.g., $10)
	const minWithdrawal = 10.0
	if req.Amount < minWithdrawal {
		return nil, fmt.Errorf("minimum withdrawal amount is %.2f", minWithdrawal)
	}

	// Start transaction
	tx, err := s.walletRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Deduct amount from wallet (hold it)
	newBalance := wallet.Balance - req.Amount
	err = s.walletRepo.UpdateWalletBalance(ctx, tx, wallet.ID, newBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Create withdrawal transaction
	transaction := &models.WalletTransaction{
		WalletID:        wallet.ID,
		UserID:          userID,
		TransactionType: models.TransactionTypeWithdrawal,
		Amount:          req.Amount,
		BalanceBefore:   wallet.Balance,
		BalanceAfter:    newBalance,
		Description:     fmt.Sprintf("Withdrawal request - %s", req.WithdrawalMethod),
		Status:          models.TransactionStatusPending,
		Metadata: models.JSONB{
			"withdrawal_method": req.WithdrawalMethod,
			"account_details":   req.AccountDetails,
		},
	}

	err = s.walletRepo.CreateTransaction(ctx, tx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Create withdrawal request
	withdrawalReq := &models.WithdrawalRequest{
		UserID:           userID,
		WalletID:         wallet.ID,
		Amount:           req.Amount,
		WithdrawalMethod: req.WithdrawalMethod,
		AccountDetails:   req.AccountDetails,
		Status:           models.WithdrawalStatusPending,
		TransactionID:    &transaction.ID,
	}

	err = s.walletRepo.CreateWithdrawalRequest(ctx, tx, withdrawalReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create withdrawal request: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return withdrawalReq, nil
}

// GetWithdrawalHistory gets withdrawal history for a user
func (s *WalletService) GetWithdrawalHistory(ctx context.Context, userID string, page, pageSize int) ([]models.WithdrawalRequest, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.walletRepo.GetWithdrawalHistory(ctx, userID, pageSize, offset)
}

// ============================================
// TRANSACTION OPERATIONS
// ============================================

// GetTransactionHistory gets transaction history for a user
func (s *WalletService) GetTransactionHistory(ctx context.Context, userID string, page, pageSize int) ([]models.WalletTransaction, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.walletRepo.GetTransactionHistory(ctx, userID, pageSize, offset)
}

// DebitWallet debits amount from user's wallet (for purchases, gifts, etc.)
func (s *WalletService) DebitWallet(ctx context.Context, userID string, amount float64, transactionType, description string, metadata models.JSONB) (*models.WalletTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Get wallet
	wallet, err := s.walletRepo.GetWalletByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Check balance
	if wallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Start transaction
	tx, err := s.walletRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Calculate new balance
	newBalance := wallet.Balance - amount

	// Update wallet balance
	err = s.walletRepo.UpdateWalletBalance(ctx, tx, wallet.ID, newBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Increment total spent
	err = s.walletRepo.IncrementTotalSpent(ctx, tx, wallet.ID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to increment total spent: %w", err)
	}

	// Create transaction record
	transaction := &models.WalletTransaction{
		WalletID:        wallet.ID,
		UserID:          userID,
		TransactionType: transactionType,
		Amount:          amount,
		BalanceBefore:   wallet.Balance,
		BalanceAfter:    newBalance,
		Description:     description,
		Status:          models.TransactionStatusCompleted,
		Metadata:        metadata,
	}

	err = s.walletRepo.CreateTransaction(ctx, tx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// CreditWallet credits amount to user's wallet (for gifts received, earnings, etc.)
func (s *WalletService) CreditWallet(ctx context.Context, userID string, amount float64, transactionType, description string, metadata models.JSONB) (*models.WalletTransaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	// Get or create wallet
	wallet, err := s.walletRepo.GetOrCreateWallet(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	// Start transaction
	tx, err := s.walletRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Calculate new balance
	newBalance := wallet.Balance + amount

	// Update wallet balance
	err = s.walletRepo.UpdateWalletBalance(ctx, tx, wallet.ID, newBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Increment total earned
	err = s.walletRepo.IncrementTotalEarned(ctx, tx, wallet.ID, amount)
	if err != nil {
		return nil, fmt.Errorf("failed to increment total earned: %w", err)
	}

	// Create transaction record
	transaction := &models.WalletTransaction{
		WalletID:        wallet.ID,
		UserID:          userID,
		TransactionType: transactionType,
		Amount:          amount,
		BalanceBefore:   wallet.Balance,
		BalanceAfter:    newBalance,
		Description:     description,
		Status:          models.TransactionStatusCompleted,
		Metadata:        metadata,
	}

	err = s.walletRepo.CreateTransaction(ctx, tx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// ============================================
// PAYOUT METHOD OPERATIONS
// ============================================

// AddPayoutMethod adds a new payout method
func (s *WalletService) AddPayoutMethod(ctx context.Context, userID string, req *models.AddPayoutMethodRequest) (*models.PayoutMethod, error) {
	// Validate method type
	validMethods := map[string]bool{
		models.PayoutMethodBankAccount: true,
		models.PayoutMethodMobileMoney: true,
		models.PayoutMethodPayPal:      true,
	}

	if !validMethods[req.MethodType] {
		return nil, errors.New("invalid payout method type")
	}

	method := &models.PayoutMethod{
		UserID:         userID,
		MethodType:     req.MethodType,
		AccountName:    req.AccountName,
		AccountDetails: req.AccountDetails,
		IsDefault:      req.IsDefault,
		IsVerified:     false, // Requires admin verification
	}

	err := s.walletRepo.CreatePayoutMethod(ctx, method)
	if err != nil {
		return nil, fmt.Errorf("failed to create payout method: %w", err)
	}

	return method, nil
}

// GetPayoutMethods gets all payout methods for a user
func (s *WalletService) GetPayoutMethods(ctx context.Context, userID string) ([]models.PayoutMethod, error) {
	return s.walletRepo.GetPayoutMethods(ctx, userID)
}

// DeletePayoutMethod deletes a payout method
func (s *WalletService) DeletePayoutMethod(ctx context.Context, methodID, userID string) error {
	return s.walletRepo.DeletePayoutMethod(ctx, methodID, userID)
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func timePtr(t time.Time) *time.Time {
	return &t
}
