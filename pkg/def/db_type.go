package def

// const (
// 	JSONB = "jsonb"
// 	JSON  = "json"
// )

const (
	PostgresJsonIndexType = "GIN"
	MySqlJsonIndexType    = "FULLTEXT"
)

var (
	JsonStruct      struct{}
	JsonMap         map[string]interface{}
	JsonArrayMap    []map[string]interface{}
	JsonArray       []interface{}
	JsonArrayString []string
)
