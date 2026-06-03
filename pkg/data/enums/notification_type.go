package enums

type NotificationType string

const (
	NotificationType_Unknown   NotificationType = "unknown"
	NotificationType_Booking   NotificationType = "booking"
	NotificationType_System    NotificationType = "system"
	NotificationType_Promotion NotificationType = "promotion"
)

func (r NotificationType) String() string {
	return string(r)
}

func (r NotificationType) IsValid() bool {
	switch r {
	case NotificationType_Unknown,
		NotificationType_Booking,
		NotificationType_System,
		NotificationType_Promotion:
		return true
	default:
		return false
	}
}

func NotificationTypeValues() []string {
	return []string{
		NotificationType_Unknown.String(),
		NotificationType_Booking.String(),
		NotificationType_System.String(),
		NotificationType_Promotion.String(),
	}
}
