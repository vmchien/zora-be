package enums

type WayType string

const (
	OneWay    WayType = "one_way"
	RoundTrip WayType = "round_trip"
)

func (r WayType) String() string {
	return string(r)
}

func (r WayType) IsValid() bool {
	switch r {
	case OneWay,
		RoundTrip:
		return true
	default:
		return false
	}
}

func WayTypeValues() []string {
	return []string{
		OneWay.String(),
		RoundTrip.String(),
	}
}
