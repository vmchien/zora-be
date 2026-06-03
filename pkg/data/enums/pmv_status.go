package enums

type PMVStatus string

const (
	PMVStatusUnknown   PMVStatus = "Unknown"
	PMVStatusWait      PMVStatus = "Wait"
	PMVStatusSuccess   PMVStatus = "Success"
	PMVStatusRefund    PMVStatus = "Refund"
	PMVStatusError     PMVStatus = "Error"
	PMVStatusCancelled PMVStatus = "Cancelled"
)

func (s PMVStatus) String() string {
	return string(s)
}

func (s PMVStatus) IsValid() bool {
	switch s {
	case PMVStatusWait,
		PMVStatusSuccess,
		PMVStatusRefund,
		PMVStatusError,
		PMVStatusCancelled:
		return true
	default:
		return false
	}
}

func PMVStatusValues() []string {
	return []string{
		PMVStatusWait.String(),
		PMVStatusSuccess.String(),
		PMVStatusRefund.String(),
		PMVStatusError.String(),
		PMVStatusCancelled.String(),
	}
}

var TicketStatusMap = map[PMVStatus]TicketStatus{
	PMVStatusWait:      TicketStatusBooking,
	PMVStatusSuccess:   TicketStatusCompleted,
	PMVStatusRefund:    TicketStatusCashBack,
	PMVStatusError:     TicketStatusBookingExpired,
	PMVStatusCancelled: TicketStatusBookingExpired,
	PMVStatusUnknown:   TicketStatusUnknown,
}

func GetPMVStatus(s string) PMVStatus {
	status := PMVStatus(s)
	if status.IsValid() {
		return status
	}
	return PMVStatusUnknown
}
