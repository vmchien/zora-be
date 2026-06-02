package def

type QueueTopic string

const (
	SystemJobTopic_InitAccount      QueueTopic = "init_system_account"
	SystemJobTopic_CleanJunkStorage QueueTopic = "clean_junk_storage"
)

const (
	TenantJobTopic_Inventory_CalculateDailyStock QueueTopic = "calculate_daily_stock"
)

var SystemQueueTopics = []QueueTopic{
	SystemJobTopic_InitAccount,
	SystemJobTopic_CleanJunkStorage,
}

var TenantQueueTopics = []QueueTopic{
	TenantJobTopic_Inventory_CalculateDailyStock,
}
