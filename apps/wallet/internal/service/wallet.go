package service

import (
	"context"

	v1 "vn.vato.zora.be.api/api/wallet/v1"
	"vn.vato.zora.be.api/apps/wallet/internal/biz"
)

type WalletService struct {
	v1.UnimplementedWalletServiceServer
	uc *biz.WalletUseCase
}

func NewWalletService(uc *biz.WalletUseCase) *WalletService {
	return &WalletService{uc: uc}
}

func (s *WalletService) GetBalance(ctx context.Context, in *v1.BalanceRequest) (*v1.BalanceReply, error) {
	_, err := s.uc.GetBalance(ctx, &biz.Wallet{UserId: in.UserId})
	if err != nil {
		return nil, err
	}
	return &v1.BalanceReply{}, nil
}

func (s *WalletService) Transfer(ctx context.Context, request *v1.BalanceRequest) (*v1.BalanceReply, error) {
	// TODO implement me
	panic("implement me")
}
