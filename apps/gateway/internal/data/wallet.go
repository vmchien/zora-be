package data

import (
	"context"

	walletv1 "vn.vato.zora.be.api/api/wallet/v1"
	"vn.vato.zora.be.api/apps/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type WalletClient struct {
	client walletv1.WalletServiceClient
	log    *log.Helper
}

func NewWalletRepo(client walletv1.WalletServiceClient, logger log.Logger) biz.WalletRepo {
	return &WalletClient{
		client: client,
		log:    log.NewHelper(logger),
	}
}

func (w WalletClient) GetBalance(ctx context.Context, wallet *biz.Wallet) (*biz.Wallet, error) {
	result, err := w.client.GetBalance(ctx, &walletv1.BalanceRequest{UserId: wallet.UserId})
	if err != nil {
		return nil, err
	}
	return &biz.Wallet{
		UserId: result.GetBalanceZns(),
	}, nil
}
