package enums

type RefundRequestState string

const (
	RefundRequestState_Init      RefundRequestState = "init"
	RefundRequestState_Approved  RefundRequestState = "approved"
	RefundRequestState_Rejected  RefundRequestState = "rejected"
	RefundRequestState_Processed RefundRequestState = "processed"
)

func (r RefundRequestState) String() string {
	return string(r)
}

func (r RefundRequestState) IsValid() bool {
	switch r {
	case RefundRequestState_Init,
		RefundRequestState_Approved,
		RefundRequestState_Rejected,
		RefundRequestState_Processed:
		return true
	default:
		return false
	}
}

func RefundRequestStateValues() []string {
	return []string{
		RefundRequestState_Init.String(),
		RefundRequestState_Approved.String(),
		RefundRequestState_Rejected.String(),
		RefundRequestState_Processed.String(),
	}
}
