package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"vn.vato.zora.be.api/apps/wallet/internal/biz"
)

type topupRepo struct {
	data *Data
	log  *log.Helper
}

func NewTopUpRepo(data *Data, logger log.Logger) biz.TopUpRepo {
	return &topupRepo{data: data, log: log.NewHelper(logger)}
}

func (r *topupRepo) CreateOrder(ctx context.Context, order *biz.TopupOrder) (*biz.TopupOrder, error) {
	// TODO: INSERT INTO topup_orders (user_id, amount, status, expired_at)
	order.ID = fmt.Sprintf("ORD-%d-%d", order.UserID, time.Now().UnixMilli())
	order.Status = "pending"
	order.ExpiredAt = time.Now().Add(15 * time.Minute)
	r.log.Infof("CreateOrder: id=%s userID=%d amount=%d", order.ID, order.UserID, order.Amount)
	return order, nil
}

func (r *topupRepo) UpdateOrder(ctx context.Context, order *biz.TopupOrder) (*biz.TopupOrder, error) {
	// TODO: UPDATE topup_orders SET status = ?, expired_at = ? WHERE id = ?
	r.log.Infof("UpdateOrder: id=%s status=%s", order.ID, order.Status)
	return order, nil
}

func (r *topupRepo) FindOrderByID(ctx context.Context, orderID string) (*biz.TopupOrder, error) {
	// TODO: SELECT * FROM topup_orders WHERE id = ?
	return &biz.TopupOrder{ID: orderID, Status: "pending"}, nil
}

func (r *topupRepo) ListOrderByStatus(ctx context.Context, status string) ([]*biz.TopupOrder, error) {
	// TODO: SELECT * FROM topup_orders WHERE status = ?
	return []*biz.TopupOrder{}, nil
}

func (r *topupRepo) ListAllOrders(ctx context.Context) ([]*biz.TopupOrder, error) {
	// TODO: SELECT * FROM topup_orders
	return []*biz.TopupOrder{}, nil
}
