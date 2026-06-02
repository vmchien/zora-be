package enums

type NotificationChannel string

const (
	NotificationChannel_Unknown NotificationChannel = "unknown"
	NotificationChannel_Push    NotificationChannel = "push"
	NotificationChannel_Sms     NotificationChannel = "sms"
	NotificationChannel_Email   NotificationChannel = "email"
)

func (r NotificationChannel) String() string {
	return string(r)
}

func (r NotificationChannel) IsValid() bool {
	switch r {
	case NotificationChannel_Unknown,
		NotificationChannel_Push,
		NotificationChannel_Sms,
		NotificationChannel_Email:
		return true
	default:
		return false
	}
}

func NotificationChannelValues() []string {
	return []string{
		NotificationChannel_Unknown.String(),
		NotificationChannel_Push.String(),
		NotificationChannel_Sms.String(),
		NotificationChannel_Email.String(),
	}
}
