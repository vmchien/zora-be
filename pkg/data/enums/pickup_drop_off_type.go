package enums

type PickupDropOffType string

const (
	PickupDropOffShuttle     PickupDropOffType = "shuttle"
	PickupDropOffSelfTransit PickupDropOffType = "self"
	PickupDropOffOffice      PickupDropOffType = "office"
	PickupDropOffTransit     PickupDropOffType = "transit"
)

func (r PickupDropOffType) String() string {
	return string(r)
}

func (r PickupDropOffType) IsValid() bool {
	switch r {
	case PickupDropOffSelfTransit,
		PickupDropOffOffice,
		PickupDropOffTransit:
		return true
	default:
		return false
	}
}

func PickupDropOffValues() []string {
	return []string{
		PickupDropOffShuttle.String(),
		PickupDropOffSelfTransit.String(),
		PickupDropOffOffice.String(),
		PickupDropOffTransit.String(),
	}
}

// TODO: check constant type -> default is office?
func GetPickupDropOffType(value float64) PickupDropOffType {
	switch value {
	case 1:
		return PickupDropOffShuttle
	case 2:
		return PickupDropOffSelfTransit
	case 3:
		return PickupDropOffOffice
	default:
		return PickupDropOffOffice
	}
}

func GetPickupDropOffValue(t PickupDropOffType) float64 {
	switch t {
	case PickupDropOffShuttle:
		return 1
	case PickupDropOffSelfTransit:
		return 2
	case PickupDropOffOffice:
		return 3
	default:
		return 3
	}
}
