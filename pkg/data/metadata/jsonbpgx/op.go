package jsonbpgx

// JSONBOperator is a type-safe enum for PostgreSQL 17 JSON/JSONB operators.
// Reference: https://www.postgresql.org/docs/17/functions-json.html
type JSONBOperator string

// String returns the string representation of the operator.
func (op JSONBOperator) String() string {
	return string(op)
}

// IsValidOperator checks if the operator is a known PostgreSQL JSONB operator.
func IsValidOperator(op JSONBOperator) bool {
	switch op {
	case OpArrow, OpArrowText, OpHash, OpHashText,
		OpContains, OpContainedIn, OpKeyExists,
		OpAnyKeyExists, OpAllKeysExist, OpJSONPathMatch,
		OpJSONPathPredicate, OpDeleteKey, OpDeletePath, OpConcat,
		OpLike, OpILike, OpNotLike, OpIn, OpNotIn, OpExists:
		return true
	default:
		return false
	}
}

const (
	OpArrow             JSONBOperator = "->"
	OpArrowText         JSONBOperator = "->>"
	OpHash              JSONBOperator = "#>"
	OpHashText          JSONBOperator = "#>>"
	OpContains          JSONBOperator = "@>"
	OpContainedIn       JSONBOperator = "<@"
	OpKeyExists         JSONBOperator = "?"
	OpAnyKeyExists      JSONBOperator = "?|"
	OpAllKeysExist      JSONBOperator = "?&"
	OpJSONPathMatch     JSONBOperator = "@?"
	OpJSONPathPredicate JSONBOperator = "@@"
	OpDeleteKey         JSONBOperator = "-"
	OpDeletePath        JSONBOperator = "#-"
	OpConcat            JSONBOperator = "||"
	OpLike              JSONBOperator = "LIKE"
	OpILike             JSONBOperator = "ILIKE"
	OpNotLike           JSONBOperator = "NOT LIKE"
	OpIn                JSONBOperator = "IN"
	OpNotIn             JSONBOperator = "NOT IN"
	OpExists            JSONBOperator = "EXISTS"
)
