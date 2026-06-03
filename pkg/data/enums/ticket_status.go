package enums

type TicketStatus string

const (
	TicketStatusUnknown              TicketStatus = "unknown"
	TicketStatusBookingCancelVato    TicketStatus = "cancel_by_app_vato"   // Khách hàng huỷ vé từ app VATO
	TicketStatusBookingCancelFuta    TicketStatus = "cancel_by_futa"       // Huỷ vé từ FuTa
	TicketStatusBooking              TicketStatus = "booking"              // Đặt giữ vé, chưa thanh toán
	TicketStatusProcessingBooking    TicketStatus = "processing_booking"   // Giao dịch đang được xử lý
	TicketStatusPaymentSuccess       TicketStatus = "payment_success"      // Thanh toán thành công
	TicketStatusPaymentFail          TicketStatus = "payment_fail"         // Thanh toán thất bại
	TicketStatusCompleted            TicketStatus = "completed"            // Xuất vé
	TicketStatusTicketUpdateFromFuta TicketStatus = "update_by_futa"       // Cập nhật trạng thái từ FUTA
	TicketStatusTicketWasChanged     TicketStatus = "ticket_was_changed"   // Vé đã được thay đổi
	TicketStatusBookingCancelAdmin   TicketStatus = "booking_cancel_admin" // Huỷ vé từ Admin
	TicketStatusCashBack             TicketStatus = "cash_back"            // Hoàn tiền
	TicketStatusProcessingCallback   TicketStatus = "processing_callback"  // Giao dịch đang được xử lý (callback)
	TicketStatusBookingExpired       TicketStatus = "expired"              // Hết hạn xử lý: thanh toán
)

func (r TicketStatus) String() string {
	return string(r)
}

func (r TicketStatus) IsValid() bool {
	switch r {
	case TicketStatusUnknown,
		TicketStatusBookingCancelVato,
		TicketStatusBookingCancelFuta,
		TicketStatusBooking,
		TicketStatusProcessingBooking,
		TicketStatusPaymentSuccess,
		TicketStatusPaymentFail,
		TicketStatusCompleted,
		TicketStatusTicketUpdateFromFuta,
		TicketStatusTicketWasChanged,
		TicketStatusBookingCancelAdmin,
		TicketStatusCashBack,
		TicketStatusProcessingCallback,
		TicketStatusBookingExpired:
		return true
	default:
		return false
	}
}

func TicketStatusValues() []string {
	return []string{
		TicketStatusUnknown.String(),
		TicketStatusBookingCancelVato.String(),
		TicketStatusBookingCancelFuta.String(),
		TicketStatusBooking.String(),
		TicketStatusProcessingBooking.String(),
		TicketStatusPaymentSuccess.String(),
		TicketStatusPaymentFail.String(),
		TicketStatusCompleted.String(),
		TicketStatusTicketUpdateFromFuta.String(),
		TicketStatusTicketWasChanged.String(),
		TicketStatusBookingCancelAdmin.String(),
		TicketStatusCashBack.String(),
		TicketStatusProcessingCallback.String(),
		TicketStatusBookingExpired.String(),
	}
}

func ParseTicketStatusFromInt(id int) TicketStatus {
	mapping := map[int]TicketStatus{
		1:  TicketStatusBookingCancelVato,
		2:  TicketStatusBookingCancelFuta,
		3:  TicketStatusBooking,
		4:  TicketStatusProcessingBooking,
		5:  TicketStatusPaymentSuccess,
		6:  TicketStatusPaymentFail,
		7:  TicketStatusCompleted,
		8:  TicketStatusTicketUpdateFromFuta,
		9:  TicketStatusTicketWasChanged,
		10: TicketStatusBookingCancelAdmin,
		11: TicketStatusCashBack,
		12: TicketStatusProcessingCallback,
	}

	status, exists := mapping[id]
	if !exists {
		return TicketStatusUnknown
	}
	return status
}
