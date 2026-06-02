package enums

type TicketSeatType string

const (
	TicketSeatTypeOutbound TicketSeatType = "outbound"
	TicketSeatTypeReturn   TicketSeatType = "return"
)

func (r TicketSeatType) String() string {
	return string(r)
}

func (r TicketSeatType) IsValid() bool {
	switch r {
	case TicketSeatTypeOutbound,
		TicketSeatTypeReturn:
		return true
	default:
		return false
	}
}

func TicketSeatTypeValues() []string {
	return []string{
		TicketSeatTypeOutbound.String(),
		TicketSeatTypeReturn.String(),
	}
}
