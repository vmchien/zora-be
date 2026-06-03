package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Wallet struct {
	UserId int64
}

// WalletRepo is a Greater repo.
type WalletRepo interface {
	GetBalance(context.Context, *Wallet) (*Wallet, error)
}

type WalletUseCase struct {
	walletRepo WalletRepo
}

func NewWalletUseCase(repo WalletRepo) *WalletUseCase {
	return &WalletUseCase{walletRepo: repo}
}

func (uc *WalletUseCase) GetBalance(ctx context.Context, g *Wallet) (*Wallet, error) {
	log.Infof("Balance: %v", g.UserId)
	payMsg, err := uc.walletRepo.GetBalance(ctx, g)
	if err != nil {
		log.Errorf("failed to call payment service: %v", err)
		return nil, err
	}
	return payMsg, nil
}
