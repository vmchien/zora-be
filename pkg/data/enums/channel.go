package enums

type Channel string

const (
	ChannelUnknown Channel = "unknown"
	ChannelWeb     Channel = "web_client"
	ChannelMobile  Channel = "mobile_app"
	ChannelKiosk   Channel = "kiosk_app"
)

func (r Channel) String() string {
	return string(r)
}

func (r Channel) IsValid() bool {
	switch r {
	case ChannelUnknown,
		ChannelWeb,
		ChannelMobile,
		ChannelKiosk:
		return true
	default:
		return false
	}
}

func ChannelValues() []string {
	return []string{
		ChannelUnknown.String(),
		ChannelWeb.String(),
		ChannelMobile.String(),
		ChannelKiosk.String(),
	}
}

func IsValidChannel(channel string) bool {
	return Channel(channel).IsValid()
}
