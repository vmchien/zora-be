package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"vn.vato.zora.be.api/apps/wallet/internal/biz"
)

type billingRepo struct {
	data *Data
	log  *log.Helper
}

func NewBillingRepo(data *Data, logger log.Logger) biz.BillingRepo {
	return &billingRepo{data: data, log: log.NewHelper(logger)}
}

func (r *billingRepo) GetBalance(ctx context.Context, userID int64) (*biz.WalletBalance, error) {
	// TODO: SELECT balance_zatify, balance_zns FROM wallets WHERE user_id = ?
	return &biz.WalletBalance{BalanceZatify: 0, BalanceZNS: 0, Currency: "VND"}, nil
}

func (r *billingRepo) DepositZatify(ctx context.Context, userID int64, amount int64, refCode string) error {
	// TODO: UPDATE wallets SET balance_zatify = balance_zatify + ? WHERE user_id = ?
	//       INSERT INTO transactions (user_id, type, amount, status, ref_code) VALUES (...)
	r.log.Infof("DepositZatify: userID=%d amount=%d refCode=%s", userID, amount, refCode)
	return nil
}

func (r *billingRepo) Withdraw(ctx context.Context, userID int64, amount int64, txType string) (*biz.WalletBalance, error) {
	// TODO: kiểm tra balance đủ, UPDATE wallets SET balance_zatify = balance_zatify - ?
	//       INSERT INTO transactions (...)
	r.log.Infof("Withdraw: userID=%d amount=%d type=%s", userID, amount, txType)
	return &biz.WalletBalance{BalanceZatify: 0, BalanceZNS: 0, Currency: "VND"}, nil
}

func (r *billingRepo) TransferZatifyToZNS(ctx context.Context, userID int64, amount int64) (*biz.WalletBalance, error) {
	// TODO: trong 1 transaction DB:
	//   UPDATE wallets SET balance_zatify = balance_zatify - ?, balance_zns = balance_zns + ? WHERE user_id = ?
	//   INSERT INTO transactions (type='transfer_zns', ...)
	r.log.Infof("TransferZatifyToZNS: userID=%d amount=%d", userID, amount)
	return &biz.WalletBalance{BalanceZatify: 0, BalanceZNS: float64(amount), Currency: "VND"}, nil
}

func (r *billingRepo) GetTransactions(ctx context.Context, userID int64, filter *biz.TransactionFilter) ([]*biz.Transaction, int, error) {
	// TODO: SELECT ... FROM transactions WHERE user_id = ? AND type = ? AND created_at BETWEEN ? AND ?
	//       ORDER BY created_at DESC LIMIT 20 OFFSET (page-1)*20
	return []*biz.Transaction{}, 0, nil
}
