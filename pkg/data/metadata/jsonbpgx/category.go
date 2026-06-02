package jsonbpgx

// Category represents a category of JSONB operations in PostgreSQL.
type Category string

// String returns the string representation of the category.
func (c Category) String() string {
	return string(c)
}

// AllCategories returns a slice of all defined JSONB operation categories.
func AllCategories() []Category {
	return []Category{
		CategoryAccess,
		CategoryComparison,
		CategoryKey,
		CategoryPath,
		CategoryMutation,
		CategoryConstruction,
		CategoryRecord,
	}
}

const (
	CategoryAccess       Category = "Access"
	CategoryComparison   Category = "Comparison"
	CategoryKey          Category = "Key"
	CategoryPath         Category = "Path"
	CategoryMutation     Category = "Mutation"
	CategoryConstruction Category = "Construction"
	CategoryRecord       Category = "Record"
)
