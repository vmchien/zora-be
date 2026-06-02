package constant

const (

	// Kafka Consumer/Producer IDs (vato/tanquangdung/etc..)
	UserNotificationConsumerKafkaId        = "futa.buslines.user.notification.consumer"
	UserTicketUpcomingConsumerKafkaId      = "futa.buslines.user.ticket.upcomming.consumer"
	UserTicketPurchaseDailyConsumerKafkaId = "futa.buslines.user.ticket.purchase.daily.consumer"
	BookingExportConsumerKafkaId           = "futa.vato.buslines.booking.export.consumer"
	BookingRefundConsumerKafkaId           = "futa.vato.buslines.booking.refund.consumer"
	SyncTicketConsumerKafkaId              = "futa.vato.buslines.booking.sync.ticket.consumer"
	AlertConsumerKafkaId                   = "futa.buslines.alert.notification.consumer"

	BookingNotificationProducerKafkaId = "futa.vato.buslines.booking.notification.producer"
	BookingRefundProducerKafkaId       = "futa.vato.buslines.booking.refund.producer"
	AlertNotificationProducerKafkaId   = "futa.vato.buslines.booking.alert.producer"
)

const (
	BookingTicketKey    = "booking"
	PaymentTicketKey    = "payment"
	ExportTicketKey     = "export"
	CancelTicketKey     = "cancel"
	RefundTicketKey     = "refund"
	ExpiredTicketKey    = "expired"
	RefundFailTicketKey = "refund_fail"
)

const (
	UserNotificationConsumerKafkaGroupId_VATO         = "futa.buslines.user.notification.vato.consumer"
	UserNotificationConsumerKafkaGroupId_TANQUANGDUNG = "futa.buslines.user.notification.tanquangdung.consumer"
)

const (
	UserNotificationProducerKafkaId = "futa.buslines.user.notification.fcm.producer"
)
