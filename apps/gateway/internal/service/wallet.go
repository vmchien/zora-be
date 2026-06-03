package service

import (
	"context"

	"vn.vato.zora.be.api/apps/gateway/internal/biz"

	v1 "vn.vato.zora.be.api/api/gateway/v1"
)

// WalletService is a greeter service.
type WalletService struct {
	v1.UnimplementedGreeterServer

	uc *biz.WalletUseCase
}

func NewWalletService(uc *biz.WalletUseCase) *WalletService {
	return &WalletService{uc: uc}
}

func (s *WalletService) GetBalance(ctx context.Context, in *v1.BalanceRequest) (*v1.BalanceReply, error) {
	_, err := s.uc.GetBalance(ctx, &biz.Wallet{UserId: int64(1)})
	if err != nil {
		return nil, err
	}
	return &v1.BalanceReply{}, nil
}
