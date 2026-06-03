package enums

type TicketSeatStatus string

const (
	TicketSeatStatusInitial        TicketSeatStatus = "initial"
	TicketSeatStatusPaid           TicketSeatStatus = "paid"
	TicketSeatStatusExported       TicketSeatStatus = "exported"
	TicketSeatStatusCancelled      TicketSeatStatus = "cancelled"
	TicketSeatStatusRefunded       TicketSeatStatus = "refunded"
	TicketSeatStatusCreateByFuta   TicketSeatStatus = "created_by_futa"
	TicketSeatStatusCanceledByFuta TicketSeatStatus = "cancelled_by_futa"
)

func (r TicketSeatStatus) String() string {
	return string(r)
}

func (r TicketSeatStatus) IsValid() bool {
	switch r {
	case TicketSeatStatusInitial,
		TicketSeatStatusPaid,
		TicketSeatStatusExported,
		TicketSeatStatusCancelled,
		TicketSeatStatusRefunded,
		TicketSeatStatusCreateByFuta,
		TicketSeatStatusCanceledByFuta:
		return true
	default:
		return false
	}
}

func TicketSeatStatusValues() []string {
	return []string{
		TicketSeatStatusInitial.String(),
		TicketSeatStatusPaid.String(),
		TicketSeatStatusExported.String(),
		TicketSeatStatusCancelled.String(),
		TicketSeatStatusRefunded.String(),
		TicketSeatStatusCreateByFuta.String(),
		TicketSeatStatusCanceledByFuta.String(),
	}
}
