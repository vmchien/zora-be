package jsonbpgx

// OperatorMeta contains metadata for JSONB operators.
type OperatorMeta struct {
	Label    string
	Symbol   JSONBOperator
	Category Category
}

// FunctionMeta contains metadata for JSONB functions.
type FunctionMeta struct {
	Label    string
	Name     JSONBFunction
	Category Category
}

// AllOperators is a list of all known JSONB operators with their metadata.
var AllOperators = []OperatorMeta{
	{"Get JSON field", OpArrow, CategoryAccess},
	{"Get JSON field as text", OpArrowText, CategoryAccess},
	{"Get JSON path", OpHash, CategoryAccess},
	{"Get JSON path as text", OpHashText, CategoryAccess},
	{"Contains (left contains right)", OpContains, CategoryComparison},
	{"Contained in (left is in right)", OpContainedIn, CategoryComparison},
	{"Key exists", OpKeyExists, CategoryKey},
	{"Any key exists", OpAnyKeyExists, CategoryKey},
	{"All keys exist", OpAllKeysExist, CategoryKey},
	{"JSONPath match", OpJSONPathMatch, CategoryPath},
	{"JSONPath predicate", OpJSONPathPredicate, CategoryPath},
	{"Delete key/index", OpDeleteKey, CategoryMutation},
	{"Delete path", OpDeletePath, CategoryMutation},
	{"Concatenate JSON", OpConcat, CategoryMutation},
}

// AllFunctions is a list of all known JSONB functions with their metadata.
var AllFunctions = []FunctionMeta{
	{"Type of JSON value", FnJsonbTypeof, CategoryAccess},
	{"Array length", FnJsonbArrayLength, CategoryAccess},
	{"Each object entry", FnJsonbEach, CategoryAccess},
	{"Each entry as text", FnJsonbEachText, CategoryAccess},
	{"Object keys", FnJsonbObjectKeys, CategoryAccess},
	{"Array elements", FnJsonbArrayElements, CategoryAccess},
	{"Array elements as text", FnJsonbArrayElementsText, CategoryAccess},
	{"Extract path", FnJsonbExtractPath, CategoryAccess},
	{"Extract path text", FnJsonbExtractPathText, CategoryAccess},
	{"Build object", FnJsonbBuildObject, CategoryConstruction},
	{"Build array", FnJsonbBuildArray, CategoryConstruction},
	{"Set field", FnJsonbSet, CategoryMutation},
	{"Insert field", FnJsonbInsert, CategoryMutation},
	{"Strip nulls", FnJsonbStripNulls, CategoryMutation},
	{"Delete key", FnJsonbDelete, CategoryMutation},
	{"Delete path", FnJsonbDeletePath, CategoryMutation},
	{"Exists", FnJsonbExists, CategoryKey},
	{"Exists any", FnJsonbExistsAny, CategoryKey},
	{"Exists all", FnJsonbExistsAll, CategoryKey},
	{"Path query", FnJsonbPathQuery, CategoryPath},
	{"Path query array", FnJsonbPathQueryArray, CategoryPath},
	{"Path query first", FnJsonbPathQueryFirst, CategoryPath},
	{"Path exists", FnJsonbPathExists, CategoryPath},
	{"Path match", FnJsonbPathMatch, CategoryPath},
	{"Path predicate", FnJsonbPathPredicate, CategoryPath},
	{"To record", FnJsonbToRecord, CategoryRecord},
	{"To recordset", FnJsonbToRecordset, CategoryRecord},
	{"Populate record", FnJsonbPopulateRecord, CategoryRecord},
	{"Populate recordset", FnJsonbPopulateRecordset, CategoryRecord},
}

// ValidateOperator checks if a given string matches a known JSONBOperator.
func ValidateOperator(op string) bool {
	for _, meta := range AllOperators {
		if string(meta.Symbol) == op {
			return true
		}
	}
	return false
}

// ValidateFunction checks if a given string matches a known JSONBFunction.
func ValidateFunction(fn string) bool {
	for _, meta := range AllFunctions {
		if string(meta.Name) == fn {
			return true
		}
	}
	return false
}

// SuggestBehaviorForOperator provides suggested behavior or usage notes for a given JSONB operator.
func SuggestBehaviorForOperator(op JSONBOperator) string {
	switch op {
	case OpArrow:
		return "Access a JSON object field by key, returns JSON value."
	case OpArrowText:
		return "Access a JSON object field by key, returns text value."
	case OpHash:
		return "Access a nested JSON value using path array, returns JSON."
	case OpHashText:
		return "Access a nested JSON value using path array, returns text."
	case OpContains:
		return "Checks if left JSON contains right JSON (e.g., for partial match)."
	case OpContainedIn:
		return "Checks if left JSON is contained within right JSON."
	case OpKeyExists:
		return "Checks if a given key exists in the JSON object."
	case OpAnyKeyExists:
		return "Checks if any of the listed keys exist in the JSON object."
	case OpAllKeysExist:
		return "Checks if all of the listed keys exist in the JSON object."
	case OpJSONPathMatch:
		return "Evaluates a JSONPath expression and returns true if matched."
	case OpJSONPathPredicate:
		return "Tests a JSONPath predicate expression for truth value."
	case OpDeleteKey:
		return "Removes a key from a JSON object or index from an array."
	case OpDeletePath:
		return "Removes a value from JSON by path."
	case OpConcat:
		return "Concatenates two JSONB values (merges them)."
	default:
		return "Unknown operator or unsupported usage."
	}
}

// SuggestBehaviorForFunction provides suggested behavior or usage notes for a given JSONB function.
func SuggestBehaviorForFunction(fn JSONBFunction) string {
	switch fn {
	case FnJsonbTypeof:
		return "Returns the type of the top-level JSON value as text."
	case FnJsonbArrayLength:
		return "Returns the number of elements in a JSON array."
	case FnJsonbEach:
		return "Expands a JSON object into a set of key/value pairs."
	case FnJsonbEachText:
		return "Expands a JSON object into a set of key/text-value pairs."
	case FnJsonbObjectKeys:
		return "Returns the set of keys in the JSON object."
	case FnJsonbArrayElements:
		return "Expands a JSON array into a set of JSON values."
	case FnJsonbArrayElementsText:
		return "Expands a JSON array into a set of text values."
	case FnJsonbExtractPath:
		return "Extracts a JSON sub-object using a path of keys."
	case FnJsonbExtractPathText:
		return "Extracts a text value from a JSON sub-object using a path of keys."
	case FnJsonbBuildObject:
		return "Builds a JSON object from key/value pairs."
	case FnJsonbBuildArray:
		return "Builds a JSON array from a variadic list of values."
	case FnJsonbSet:
		return "Updates an existing JSON value by setting a specified path to a new value."
	case FnJsonbInsert:
		return "Inserts a value at the specified location in a JSON object or array."
	case FnJsonbStripNulls:
		return "Removes all object fields with null values from the input JSON."
	case FnJsonbDelete:
		return "Deletes a key or index from a JSON object or array."
	case FnJsonbDeletePath:
		return "Deletes a value at a specified path from a JSON object."
	case FnJsonbExists:
		return "Checks whether a key exists in a JSON object."
	case FnJsonbExistsAny:
		return "Checks whether any keys in an array exist in the JSON object."
	case FnJsonbExistsAll:
		return "Checks whether all keys in an array exist in the JSON object."
	case FnJsonbPathQuery:
		return "Queries JSON data using a JSONPath expression."
	case FnJsonbPathQueryArray:
		return "Returns the result of a JSONPath query as an array."
	case FnJsonbPathQueryFirst:
		return "Returns the first result of a JSONPath query."
	case FnJsonbPathExists:
		return "Checks whether a JSONPath matches any part of the JSON data."
	case FnJsonbPathMatch:
		return "Tests whether a JSONPath expression matches the JSON structure."
	case FnJsonbPathPredicate:
		return "Evaluates a boolean JSONPath predicate expression."
	case FnJsonbToRecord:
		return "Expands a JSON object to a SQL composite type."
	case FnJsonbToRecordset:
		return "Expands a JSON array of objects to a set of SQL composite types."
	case FnJsonbPopulateRecord:
		return "Populates a SQL composite type from a JSON object, with existing defaults."
	case FnJsonbPopulateRecordset:
		return "Populates a set of SQL composite types from a JSON array of objects."
	default:
		return "Unknown function or unsupported usage."
	}
}
