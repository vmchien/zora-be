package constant

const (
	// Ticket lifecycle
	NotificationEventGroup_Exported  = "exported"
	NotificationEventGroup_Cancelled = "cancelled"
	NotificationEventGroup_Refunded  = "refunded"
	NotificationEventGroup_Changed   = "changed"

	// Cancel flow (separated to avoid sending OTP + success templates together)
	NotificationEventGroup_CancelOtp     = "cancel_otp"
	NotificationEventGroup_CancelVatoPay = "cancel_vatopay"
	NotificationEventGroup_CancelEWallet = "cancel_e_wallet"
	NotificationEventGroup_CancelNapas   = "cancel_napas"
	NotificationEventGroup_CancelAdmin   = "cancel_admin"
	NotificationEventGroup_CancelGeneric = "cancel_generic"

	// Payment flow (separated to avoid conflict between success and fail)
	NotificationEventGroup_PaymentProcessing = "payment_processing"
	NotificationEventGroup_PaymentSuccess    = "payment_success"
	NotificationEventGroup_PaymentFail       = "payment_fail"
	NotificationEventGroup_RefundFail        = "refund_fail"

	// Shuttle / pickup-dropoff
	NotificationEventGroup_Shuttle = "pickup_drop_off"
)

const (
	NotificationTemplateCode_Shuttle = "PICKUP_DROP_OFF"
	NotificationTemplateCode_EndTrip = "END_TRIP"
)

const (
	NotificationJob_MaxRetries = 3
)
