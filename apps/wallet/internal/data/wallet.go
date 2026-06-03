package data

import (
	"context"

	"vn.vato.zora.be.api/apps/wallet/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type walletRepo struct {
	data *Data
	log  *log.Helper
}

// NewWalletRepo .
func NewWalletRepo(data *Data, logger log.Logger) biz.WalletRepo {
	return &walletRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (w walletRepo) GetBalance(ctx context.Context, wallet *biz.Wallet) (*biz.Wallet, error) {
	return &biz.Wallet{}, nil
}
