package enums

type RefundJobStatus string

const (
	RefundJobStatus_Pending      RefundJobStatus = "pending"
	RefundJobStatus_Processing   RefundJobStatus = "processing"
	RefundJobStatus_Succeeded    RefundJobStatus = "succeeded"
	RefundJobStatus_Failed       RefundJobStatus = "failed"
	RefundJobStatus_ManualReview RefundJobStatus = "manual_review"
)

func (r RefundJobStatus) String() string {
	return string(r)
}

func (r RefundJobStatus) IsValid() bool {
	switch r {
	case RefundJobStatus_Pending,
		RefundJobStatus_Processing,
		RefundJobStatus_Succeeded,
		RefundJobStatus_Failed,
		RefundJobStatus_ManualReview:
		return true
	default:
		return false
	}
}

func RefundJobStatusValues() []string {
	return []string{
		RefundJobStatus_Pending.String(),
		RefundJobStatus_Processing.String(),
		RefundJobStatus_Succeeded.String(),
		RefundJobStatus_Failed.String(),
		RefundJobStatus_ManualReview.String(),
	}
}
