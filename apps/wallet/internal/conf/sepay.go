package conf

// Data_Sepay holds Sepay payment gateway credentials.
// Mirrors the Sepay message in conf.proto.
// This file can be removed after running `make proto` to regenerate conf.pb.go.
type Data_Sepay struct {
	ApiKey        string `json:"api_key"`
	WebhookSecret string `json:"webhook_secret"`
}
