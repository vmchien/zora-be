package jsonbpgx

// JSONBFunction is a type-safe enum for PostgreSQL 17 JSONB functions.
// Reference: https://www.postgresql.org/docs/17/functions-json.html
type JSONBFunction string

// String returns the string representation of the function.
func (fn JSONBFunction) String() string {
	return string(fn)
}

// IsValidFunction checks if the function is a known PostgreSQL JSONB function.
func IsValidFunction(fn JSONBFunction) bool {
	switch fn {
	case
		FnJsonbTypeof, FnJsonbArrayLength, FnJsonbEach, FnJsonbEachText,
		FnJsonbObjectKeys, FnJsonbArrayElements, FnJsonbArrayElementsText,
		FnJsonbExtractPath, FnJsonbExtractPathText, FnJsonbBuildObject,
		FnJsonbBuildArray, FnJsonbSet, FnJsonbInsert, FnJsonbStripNulls,
		FnJsonbDelete, FnJsonbDeletePath, FnJsonbExists, FnJsonbExistsAny,
		FnJsonbExistsAll, FnJsonbPathQuery, FnJsonbPathQueryArray,
		FnJsonbPathQueryFirst, FnJsonbPathExists, FnJsonbPathMatch,
		FnJsonbPathPredicate, FnJsonbToRecord, FnJsonbToRecordset,
		FnJsonbPopulateRecord, FnJsonbPopulateRecordset:
		return true
	default:
		return false
	}
}

const (
	// Accessors and navigation
	FnJsonbTypeof            JSONBFunction = "jsonb_typeof"
	FnJsonbArrayLength       JSONBFunction = "jsonb_array_length"
	FnJsonbEach              JSONBFunction = "jsonb_each"
	FnJsonbEachText          JSONBFunction = "jsonb_each_text"
	FnJsonbObjectKeys        JSONBFunction = "jsonb_object_keys"
	FnJsonbArrayElements     JSONBFunction = "jsonb_array_elements"
	FnJsonbArrayElementsText JSONBFunction = "jsonb_array_elements_text"
	FnJsonbExtractPath       JSONBFunction = "jsonb_extract_path"
	FnJsonbExtractPathText   JSONBFunction = "jsonb_extract_path_text"

	// Construction
	FnJsonbBuildObject JSONBFunction = "jsonb_build_object"
	FnJsonbBuildArray  JSONBFunction = "jsonb_build_array"

	// Manipulation
	FnJsonbSet        JSONBFunction = "jsonb_set"
	FnJsonbInsert     JSONBFunction = "jsonb_insert"
	FnJsonbStripNulls JSONBFunction = "jsonb_strip_nulls"
	FnJsonbDelete     JSONBFunction = "jsonb_delete"
	FnJsonbDeletePath JSONBFunction = "jsonb_delete_path"

	// Comparison & containment
	FnJsonbExists    JSONBFunction = "jsonb_exists"
	FnJsonbExistsAny JSONBFunction = "jsonb_exists_any"
	FnJsonbExistsAll JSONBFunction = "jsonb_exists_all"

	// Path queries
	FnJsonbPathQuery      JSONBFunction = "jsonb_path_query"
	FnJsonbPathQueryArray JSONBFunction = "jsonb_path_query_array"
	FnJsonbPathQueryFirst JSONBFunction = "jsonb_path_query_first"
	FnJsonbPathExists     JSONBFunction = "jsonb_path_exists"
	FnJsonbPathMatch      JSONBFunction = "jsonb_path_match"
	FnJsonbPathPredicate  JSONBFunction = "jsonb_path_predicate"

	// Misc
	FnJsonbToRecord          JSONBFunction = "jsonb_to_record"
	FnJsonbToRecordset       JSONBFunction = "jsonb_to_recordset"
	FnJsonbPopulateRecord    JSONBFunction = "jsonb_populate_record"
	FnJsonbPopulateRecordset JSONBFunction = "jsonb_populate_recordset"
)
