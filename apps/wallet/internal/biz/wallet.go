package biz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// ──────────────────────────────────────────────
// Domain models
// ──────────────────────────────────────────────

type WalletBalance struct {
	BalanceZatify float64
	BalanceZNS    float64
	Currency      string
}

type TopupOrder struct {
	ID        string
	UserID    int64
	Amount    int64
	Status    string // pending | completed | expired
	ExpiredAt time.Time
}

type Transaction struct {
	ID        string
	Type      string // topup | transfer_zns | service_fee | oa_package | zns_cost
	Amount    float64
	Status    string // success | failed | pending
	RefCode   string
	CreatedAt time.Time
}

type TransactionFilter struct {
	Type string
	From string
	To   string
	Page int
}

type SepayWebhookPayload struct {
	OrderID   string
	Amount    int64
	Status    string
	Signature string
}

// ──────────────────────────────────────────────
// Repository / gateway interfaces (implemented by data layer)
// ──────────────────────────────────────────────

type BillingRepo interface {
	GetBalance(ctx context.Context, userID int64) (*WalletBalance, error)
	GetTransactions(ctx context.Context, userID int64, filter *TransactionFilter) ([]*Transaction, int, error)
	DepositZatify(ctx context.Context, userID int64, amount int64, refCode string) error
	Withdraw(ctx context.Context, userID int64, amount int64, txType string) (*WalletBalance, error)
	TransferZatifyToZNS(ctx context.Context, userID int64, amount int64) (*WalletBalance, error)
}

type TopUpRepo interface {
	CreateOrder(context.Context, *TopupOrder) (*TopupOrder, error)
	UpdateOrder(context.Context, *TopupOrder) (*TopupOrder, error)
	FindOrderByID(context.Context, string) (*TopupOrder, error)
	ListOrderByStatus(context.Context, string) ([]*TopupOrder, error)
	ListAllOrders(context.Context) ([]*TopupOrder, error)
}

// PaymentGateway abstracts the external payment provider.
// Defined here in biz so data/client layer depends on biz, not vice versa.
type PaymentGateway interface {
	CreateQR(ctx context.Context, orderID string, amount int64) (qrURL string, expiredAt time.Time, err error)
	VerifyWebhookSignature(ctx context.Context, payload *SepayWebhookPayload) error
}

// ──────────────────────────────────────────────
// Use case
// ──────────────────────────────────────────────

type BillingUseCase struct {
	billingRepo BillingRepo
	topUpRepo   TopUpRepo
	sePayRepo   PaymentGateway
	log         *log.Helper
}

func NewBillingUseCase(billingRepo BillingRepo, topUpRepo TopUpRepo, gateway PaymentGateway, logger log.Logger) *BillingUseCase {
	return &BillingUseCase{
		billingRepo: billingRepo,
		topUpRepo:   topUpRepo,
		sePayRepo:   gateway,
		log:         log.NewHelper(logger),
	}
}

func (uc *BillingUseCase) GetBalance(ctx context.Context, userID int64) (*WalletBalance, error) {
	return uc.billingRepo.GetBalance(ctx, userID)
}

func (uc *BillingUseCase) InitTopUp(ctx context.Context, userID int64, amount int64) (*string, *time.Time, error) {
	topUp, err := uc.topUpRepo.CreateOrder(ctx, &TopupOrder{UserID: userID, Amount: amount})
	if err != nil {
		return nil, nil, fmt.Errorf("init QR failed: %w", err)
	}
	qrCode, expiredAt, err := uc.sePayRepo.CreateQR(ctx, topUp.ID, topUp.Amount)
	if err != nil {
		return nil, nil, fmt.Errorf("create QR failed: %w", err)
	}
	return &qrCode, &expiredAt, nil
}

func (uc *BillingUseCase) IPNTopUp(ctx context.Context, payload *SepayWebhookPayload) error {
	if err := uc.sePayRepo.VerifyWebhookSignature(ctx, payload); err != nil {
		return fmt.Errorf("invalid webhook signature: %w", err)
	}

	order, err := uc.topUpRepo.FindOrderByID(ctx, payload.OrderID)
	if err != nil {
		return err
	}
	if order.Status != "pending" {
		return nil // idempotent
	}

	return uc.billingRepo.DepositZatify(ctx, order.UserID, payload.Amount, "")
}

func (uc *BillingUseCase) GetTopUpStatus(ctx context.Context, orderID string) (*TopupOrder, error) {
	return uc.topUpRepo.FindOrderByID(ctx, orderID)
}

func (uc *BillingUseCase) TransferToZNS(ctx context.Context, userID int64, amount int64, promotionCode string) (*WalletBalance, error) {
	if amount < 1000 {
		return nil, errors.New("amount must be at least 1000 VND")
	}
	// TODO: apply promotionCode discount
	return uc.billingRepo.TransferZatifyToZNS(ctx, userID, amount)
}

func (uc *BillingUseCase) GetTransactions(ctx context.Context, userID int64, filter *TransactionFilter) ([]*Transaction, int, error) {
	return uc.billingRepo.GetTransactions(ctx, userID, filter)
}
