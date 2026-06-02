package enums

type NotificationJobStatus string

const (
	NotificationJobStatus_Failed     NotificationJobStatus = "failed"
	NotificationJobStatus_Pending    NotificationJobStatus = "pending"
	NotificationJobStatus_Processing NotificationJobStatus = "processing"
	NotificationJobStatus_Sent       NotificationJobStatus = "sent"
)

func (r NotificationJobStatus) String() string {
	return string(r)
}

func (r NotificationJobStatus) IsValid() bool {
	switch r {
	case NotificationJobStatus_Failed,
		NotificationJobStatus_Pending,
		NotificationJobStatus_Processing,
		NotificationJobStatus_Sent:
		return true
	default:
		return false
	}
}

func NotificationJobStatusValues() []string {
	return []string{
		NotificationJobStatus_Failed.String(),
		NotificationJobStatus_Pending.String(),
		NotificationJobStatus_Processing.String(),
		NotificationJobStatus_Sent.String(),
	}
}
