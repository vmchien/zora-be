package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Wallet struct {
	UserId int64
}

type WalletRepo interface {
	GetBalance(context.Context, *Wallet) (*Wallet, error)
}

type WalletUseCase struct {
	repo WalletRepo
}

func NewWalletUseCase(repo WalletRepo) *WalletUseCase {
	return &WalletUseCase{repo: repo}
}

// GetBalance creates a Booking, and returns the new Booking.
func (uc *WalletUseCase) GetBalance(ctx context.Context, g *Wallet) (*Wallet, error) {
	log.Infof("GetBalance: %v", g.UserId)
	_, err := uc.repo.GetBalance(ctx, g)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
