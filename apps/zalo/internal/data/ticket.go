package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"vn.vato.zora.be.api/apps/zalo/internal/data/ent"
	"vn.vato.zora.be.api/apps/zalo/internal/data/ent/ticket"
	"vn.vato.zora.be.api/pkg/logs"
)

type ticketRepo struct {
	data *Data
	log  *logs.Helper
}

// func NewTicketRepo(data *Data, logger log.Logger) biz.TicketRepo {
func NewTicketRepo(data *Data, logger log.Logger) any {
	return &ticketRepo{
		data: data,
		log:  logs.NewHelper(logger),
	}
}

// func (r *ticketRepo) Create(ctx context.Context, input biz.CreateTicketDTO) (*ent.Ticket, error) {
func (r *ticketRepo) Create(ctx context.Context, input any) (*ent.Ticket, error) {
	query := r.data.Client.Ticket.Create()
	// query.SetNillableTenantID(input.TenantID)
	// query.SetCode(input.Code)
	// query.SetNillableStatus(input.Status)
	// query.SetNillableWayType(input.WayType)
	// query.SetNillableBookingChannel(input.BookingChannel)
	// query.SetNillablePaymentMethod(input.PaymentMethod)
	// query.SetNillableIsActive(input.IsActive)
	// query.SetNillablePromotionCode(input.PromotionCode)
	// query.SetDepartureTimes(input.DepartureTimes)
	// query.SetNillableOriginAmount(input.OriginAmount)
	// query.SetNillableDiscountAmount(input.DiscountAmount)
	// query.SetNillableRefundAmount(input.RefundAmount)
	// query.SetNillableRemarks(input.Remarks)
	// query.SetReferenceUserID(input.ReferenceUserID)
	// query.SetReferencePhone(input.ReferencePhone)
	// query.SetNillableReferenceMigrationID(input.ReferenceMigrationID)

	return query.Save(ctx)
}

// func (r *ticketRepo) CreateTx(ctx context.Context, tx *ent.Tx, input biz.CreateTicketDTO) (*ent.Ticket, error) {
func (r *ticketRepo) CreateTx(ctx context.Context, tx *ent.Tx, input any) (*ent.Ticket, error) {
	query := r.getClient(tx).Create()
	return query.Save(ctx)
}

func (r *ticketRepo) Update(ctx context.Context, id uuid.UUID, any2 any) (*ent.Ticket, error) {
	// query := r.getClient(nil).Update().Where(ticket.ID(id))
	return nil, nil
}

func (r *ticketRepo) Patch(ctx context.Context, id uuid.UUID, input any) (*ent.Ticket, error) {
	query := r.data.Client.Ticket.UpdateOneID(id)
	return query.Save(ctx)
}

func (r *ticketRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.data.Client.Ticket.DeleteOneID(id).Exec(ctx)
}

func (r *ticketRepo) GetByID(ctx context.Context, id uuid.UUID) (*ent.Ticket, error) {
	return r.data.Client.Ticket.
		Query().
		Where(ticket.ID(id)).
		Only(ctx)
}

func (r *ticketRepo) GetByCode(ctx context.Context, code string) (*ent.Ticket, error) {
	return r.data.Client.Ticket.
		Query().
		Where(ticket.Code(code)).
		Only(ctx)
}
func (r *ticketRepo) GetByTicketWithSeat(ctx context.Context, id *uuid.UUID, code *string) (any, error) {
	if id == nil && (code == nil || *code == "") {
		return nil, fmt.Errorf("id or code must be provided")
	}

	query := r.getClient(nil).Query()
	if id != nil {
		query = query.Where(ticket.ID(*id))
	}
	if code != nil && *code != "" {
		query = query.Where(ticket.Code(*code))
	}

	_, err := query.
		WithTicketSeats().
		First(ctx)

	if err != nil {
		return nil, err
	}

	// seats := make([]biz.TicketSeatDTO, 0, len(o.Edges.TicketSeats))
	// for _, v := range o.Edges.TicketSeats {
	// 	var seat biz.TicketSeatDTO
	// 	if e := copier.Copy(&seat, &v); e != nil {
	// 		return nil, e
	// 	}
	// 	seats = append(seats, seat)
	// }
	//
	// var res biz.TicketDTO
	// if e := copier.Copy(&res, &o); e != nil {
	// 	return nil, e
	// }
	//
	// res.TicketSeats = seats
	//
	// if o.Edges.TicketExtra != nil {
	// 	var extra biz.TicketExtraDTO
	// 	if e := copier.Copy(&extra, o.Edges.TicketExtra); e != nil {
	// 		return nil, e
	// 	}
	// 	res.TicketExtras = append(res.TicketExtras, extra)
	// }
	return nil, nil
}

func (r *ticketRepo) GetForUpdateByID(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.Ticket, error) {
	return r.getClient(tx).
		Query().
		Where(ticket.ID(id)).
		ForUpdate().
		Only(ctx)
}

func (r *ticketRepo) UpdateStatusAndMethodTx(
	ctx context.Context,
	tx *ent.Tx,
	id uuid.UUID,
	status ticket.Status,
	paymentMethod *int,
) (
	*ent.Ticket, error) {
	query := r.getClient(tx).UpdateOneID(id).
		SetStatus(status)
	if paymentMethod != nil {
		query.SetPaymentMethod(*paymentMethod)
	}
	return query.Save(ctx)
}

func (r *ticketRepo) UpdateStatusAndMethodCancelTx(ctx context.Context, tx *ent.Tx, id uuid.UUID, status ticket.Status, byUserId string) (*ent.Ticket, error) {
	query := r.getClient(tx).UpdateOneID(id).
		SetStatus(status)

	if byUserId != "" {
		query.SetForceUpdatedUserID(byUserId)
	}

	return query.Save(ctx)
}

func (r *ticketRepo) UpdateStatusTx(ctx context.Context, tx *ent.Tx, id uuid.UUID, status ticket.Status) (*ent.Ticket, error) {
	query := r.getClient(tx).UpdateOneID(id).
		SetStatus(status)
	return query.Save(ctx)
}

func (r *ticketRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status ticket.Status) (*ent.Ticket, error) {
	query := r.data.Client.Ticket.UpdateOneID(id).
		SetStatus(status)
	return query.Save(ctx)
}

func (r *ticketRepo) UpdateMethod(ctx context.Context, id uuid.UUID, method int) (*ent.Ticket, error) {
	query := r.data.Client.Ticket.UpdateOneID(id).
		SetPaymentMethod(method)
	return query.Save(ctx)
}

func (r *ticketRepo) getClient(tx *ent.Tx) *ent.TicketClient {
	if tx != nil {
		return tx.Ticket
	}
	return r.data.Client.Ticket
}

func (r *ticketRepo) WithTransaction(ctx context.Context, fn func(tx *ent.Tx) error) error {
	err := r.data.WithTx(ctx, fn)
	if err != nil {
		return err
	}

	return nil
}

// func (r *ticketRepo) BatchUpdateTicketExpired(ctx context.Context) error {
// 	now := time.Now()
// 	fromID := guid.NewMaxTime(now.Add(-5 * time.Minute))
// 	forceExpireID := guid.NewMaxTime(now.Add(-60 * time.Minute))
// 	return r.data.Client.Ticket.Update().
// 		SetStatus(ticket.StatusExpired).
// 		Where(
// 			ticket.IDLT(fromID),
// 			ticket.StatusIn(ticket.StatusBooking, ticket.StatusProcessingBooking),
// 			ticket.Or(
// 				ticket.IDLTE(forceExpireID),
// 				ticket.HasTicketExtraWith(func(selector *sql.Selector) {
// 					selector.Where(
// 						// sql.ExprP("extra_data->>'timeExpiredPayment' IS NOT NULL AND (extra_data->>'timeExpiredPayment')::timestamptz < ?", time.Now()),
// 						sql.And(
// 							sql.ExprP("data->>'timeExpiredPayment' IS NOT NULL"),
// 							sql.ExprP("data->>'timeExpiredPayment' != ''"),
// 							sql.ExprP("data->>'timeExpiredPayment' != '0001-01-01T00:00:00Z'"),
// 							sql.ExprP("(data->>'timeExpiredPayment')::timestamptz < CURRENT_TIMESTAMP"),
// 						),
// 					)
// 				}),
// 			),
// 		).Exec(ctx)
// }

func (r *ticketRepo) UpdateInfoRefundedTx(ctx context.Context, tx *ent.Tx, id uuid.UUID, refundAmount float64) error {
	return tx.Ticket.Update().
		Where(
			ticket.IDEQ(id),
			ticket.StatusEQ(ticket.StatusProcessingCallback),
		).
		SetRefundAmount(refundAmount).
		SetStatus(ticket.StatusCashBack).
		Exec(ctx)
}
