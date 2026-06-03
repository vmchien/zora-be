package enums

type RefundRequestStatus string

const (
	RefundRequestStatus_Draft      RefundRequestStatus = "draft"
	RefundRequestStatus_Pending    RefundRequestStatus = "pending"
	RefundRequestStatus_Processing RefundRequestStatus = "processing"
	RefundRequestStatus_Succeeded  RefundRequestStatus = "succeeded"
	RefundRequestStatus_Failed     RefundRequestStatus = "failed"
)

func (r RefundRequestStatus) String() string {
	return string(r)
}

func (r RefundRequestStatus) IsValid() bool {
	switch r {
	case RefundRequestStatus_Draft,
		RefundRequestStatus_Pending,
		RefundRequestStatus_Processing,
		RefundRequestStatus_Succeeded,
		RefundRequestStatus_Failed:
		return true
	default:
		return false
	}
}

func RefundRequestStatusValues() []string {
	return []string{
		RefundRequestStatus_Draft.String(),
		RefundRequestStatus_Pending.String(),
		RefundRequestStatus_Processing.String(),
		RefundRequestStatus_Succeeded.String(),
		RefundRequestStatus_Failed.String(),
	}
}
