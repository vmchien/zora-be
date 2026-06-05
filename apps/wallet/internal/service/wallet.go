package service

import (
	"context"
	"time"

	v1 "vn.vato.zora.be.api/api/wallet/v1"
	"vn.vato.zora.be.api/apps/wallet/internal/biz"
)

type WalletService struct {
	v1.UnimplementedWalletServiceServer
	billing *biz.BillingUseCase
}

func NewWalletService(billing *biz.BillingUseCase) *WalletService {
	return &WalletService{billing: billing}
}

// GetWallets – GET /billing/wallets
func (s *WalletService) GetWallets(ctx context.Context, request *v1.GetWalletsRequest) (*v1.GetWalletsReply, error) {
	userID := extractUserIDFromCtx(ctx)
	balance, err := s.billing.GetBalance(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &v1.GetWalletsReply{
		Balance: &v1.WalletBalance{
			BalanceZatify: balance.BalanceZatify,
			BalanceZns:    balance.BalanceZNS,
			Currency:      balance.Currency,
		},
	}, nil
}

// InitTopup – POST /billing/topup/init
func (s *WalletService) InitTopup(ctx context.Context, in *v1.InitTopupRequest) (*v1.InitTopupReply, error) {
	userID := extractUserIDFromCtx(ctx)
	qrCode, expiredAt, err := s.billing.InitTopUp(ctx, userID, in.Amount)
	if err != nil {
		return nil, err
	}
	return &v1.InitTopupReply{
		QrUrl:     *qrCode,
		Amount:    in.Amount,
		ExpiredAt: expiredAt.Format(time.RFC3339),
	}, nil
}

// GetTopupStatus – GET /billing/topup/{order_id}
func (s *WalletService) GetTopupStatus(ctx context.Context, in *v1.GetTopupStatusRequest) (*v1.GetTopupStatusReply, error) {
	order, err := s.billing.GetTopUpStatus(ctx, in.OrderId)
	if err != nil {
		return nil, err
	}
	return &v1.GetTopupStatusReply{Status: order.Status}, nil
}

// TransferToZNS – POST /billing/transfer-to-zns
func (s *WalletService) TransferToZNS(ctx context.Context, in *v1.TransferToZNSRequest) (*v1.GetWalletsReply, error) {
	userID := extractUserIDFromCtx(ctx)
	balance, err := s.billing.TransferToZNS(ctx, userID, in.Amount, in.PromotionCode)
	if err != nil {
		return nil, err
	}
	return &v1.GetWalletsReply{
		Balance: &v1.WalletBalance{
			BalanceZatify: balance.BalanceZatify,
			BalanceZns:    balance.BalanceZNS,
			Currency:      balance.Currency,
		},
	}, nil
}

// GetTransactions – GET /billing/transactions
func (s *WalletService) GetTransactions(ctx context.Context, in *v1.GetTransactionsRequest) (*v1.GetTransactionsReply, error) {
	userID := extractUserIDFromCtx(ctx)
	filter := &biz.TransactionFilter{
		Type: in.Type,
		From: in.From,
		To:   in.To,
		Page: int(in.Page),
	}
	txs, total, err := s.billing.GetTransactions(ctx, userID, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.Transaction, 0, len(txs))
	for _, tx := range txs {
		items = append(items, &v1.Transaction{
			Id:        tx.ID,
			Type:      tx.Type,
			Amount:    tx.Amount,
			Status:    tx.Status,
			RefCode:   tx.RefCode,
			CreatedAt: tx.CreatedAt.Format(time.RFC3339),
		})
	}
	return &v1.GetTransactionsReply{
		Data: items,
		Pagination: &v1.Pagination{
			Page:    int32(filter.Page),
			PerPage: 20,
			Total:   int32(total),
		},
	}, nil
}

// SepayWebhook – POST /webhook/sepay
func (s *WalletService) SepayWebhook(ctx context.Context, in *v1.SepayWebhookRequest) (*v1.EmptyReply, error) {
	payload := &biz.SepayWebhookPayload{
		OrderID:   in.OrderId,
		Amount:    in.Amount,
		Status:    in.Status,
		Signature: in.Signature,
	}
	if err := s.billing.IPNTopUp(ctx, payload); err != nil {
		return nil, err
	}
	return &v1.EmptyReply{}, nil
}

// extractUserIDFromCtx extracts the authenticated user ID from gRPC metadata context.
// Returns 0 for unauthenticated requests (e.g. Sepay webhook).
func extractUserIDFromCtx(ctx context.Context) int64 {
	type ctxKey string
	v := ctx.Value(ctxKey("user_id"))
	if v == nil {
		return 0
	}
	if id, ok := v.(int64); ok {
		return id
	}
	return 0
}
