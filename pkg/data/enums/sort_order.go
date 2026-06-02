package enums

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func (r SortOrder) String() string {
	return string(r)
}

func (r SortOrder) IsValid() bool {
	switch r {
	case SortOrderAsc,
		SortOrderDesc:
		return true
	default:
		return false
	}
}

func SortOrderValues() []string {
	return []string{
		SortOrderAsc.String(),
		SortOrderDesc.String(),
	}
}
