package def

import (
	"entgo.io/ent/dialect/sql"
	"fmt"
	"strings"
)

type Operator string

type Predicate struct {
	Operator        Operator
	Name            string
	SqlOperatorName string
}
type FieldPredicate struct {
	Name              string
	AllowedTypes      []string
	AllowedPredicates []Predicate
}

type FieldPredicateFilter struct {
	Field    string
	Operator string
	Value    any
}

const DefaultFieldPredicate = EQ

const (
	Unknown      Operator = ""
	EQ           Operator = "="
	NEQ          Operator = "!="
	GT           Operator = ">"
	GTE          Operator = ">="
	LT           Operator = "<"
	LTE          Operator = "=<"
	In           Operator = "IN"
	NotIn        Operator = "NOT IN"
	Contains     Operator = "Contains"
	HasPrefix    Operator = "HasPrefix"
	HasSuffix    Operator = "HasSuffix"
	ContainsFold Operator = "ContainsFold"
	EqualFold    Operator = "EqualFold"
	IsNull       Operator = "IsNil"
	NotNull      Operator = "NotNil"
	// TODO: define for JSON
)

var operatorNameMap = map[Operator]string{
	EQ:           "EQ",
	NEQ:          "NEQ",
	GT:           "GT",
	GTE:          "GTE",
	LT:           "LT",
	LTE:          "LTE",
	In:           "In",
	NotIn:        "NotIn",
	Contains:     "Contains",
	HasPrefix:    "HasPrefix",
	HasSuffix:    "HasSuffix",
	ContainsFold: "ContainsFold",
	EqualFold:    "EqualFold",
	IsNull:       "IsNull",
	NotNull:      "NotNull",
}

var FieldPredicates = []FieldPredicate{
	{
		Name:         "Bool",
		AllowedTypes: []string{"bool"},
		AllowedPredicates: []Predicate{
			{
				Operator:        EQ,
				Name:            operatorNameMap[EQ],
				SqlOperatorName: operatorNameMap[EQ],
			},
			{
				Operator:        NEQ,
				Name:            operatorNameMap[NEQ],
				SqlOperatorName: operatorNameMap[NEQ],
			},
		},
	},
	{
		Name: "Numeric",
		AllowedTypes: []string{
			"int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"float32", "float64",
		},
		AllowedPredicates: []Predicate{
			{
				Operator:        EQ,
				Name:            operatorNameMap[EQ],
				SqlOperatorName: operatorNameMap[EQ],
			},
			{
				Operator:        NEQ,
				Name:            operatorNameMap[NEQ],
				SqlOperatorName: operatorNameMap[NEQ],
			},
			{
				Operator:        GT,
				Name:            operatorNameMap[GT],
				SqlOperatorName: operatorNameMap[GT],
			},
			{
				Operator:        GTE,
				Name:            operatorNameMap[GTE],
				SqlOperatorName: operatorNameMap[GTE],
			},
			{
				Operator:        LT,
				Name:            operatorNameMap[LT],
				SqlOperatorName: operatorNameMap[LT],
			},
			{
				Operator:        LTE,
				Name:            operatorNameMap[LTE],
				SqlOperatorName: operatorNameMap[LTE],
			},
			{
				Operator:        In,
				Name:            operatorNameMap[In],
				SqlOperatorName: operatorNameMap[In],
			},
			{
				Operator:        NotIn,
				Name:            operatorNameMap[NotIn],
				SqlOperatorName: operatorNameMap[NotIn],
			},
		},
	},
	{
		Name:         "Time",
		AllowedTypes: []string{"time", "time.Time"},
		AllowedPredicates: []Predicate{
			{
				Operator:        EQ,
				Name:            operatorNameMap[EQ],
				SqlOperatorName: operatorNameMap[EQ],
			},
			{
				Operator:        NEQ,
				Name:            operatorNameMap[NEQ],
				SqlOperatorName: operatorNameMap[NEQ],
			},
			{
				Operator:        GT,
				Name:            operatorNameMap[GT],
				SqlOperatorName: operatorNameMap[GT],
			},
			{
				Operator:        GTE,
				Name:            operatorNameMap[GTE],
				SqlOperatorName: operatorNameMap[GTE],
			},
			{
				Operator:        LT,
				Name:            operatorNameMap[LT],
				SqlOperatorName: operatorNameMap[LT],
			},
			{
				Operator:        LTE,
				Name:            operatorNameMap[LTE],
				SqlOperatorName: operatorNameMap[LTE],
			},
			{
				Operator:        In,
				Name:            operatorNameMap[In],
				SqlOperatorName: operatorNameMap[In],
			},
			{
				Operator:        NotIn,
				Name:            operatorNameMap[NotIn],
				SqlOperatorName: operatorNameMap[NotIn],
			},
			{
				Operator:        IsNull,
				Name:            operatorNameMap[IsNull],
				SqlOperatorName: operatorNameMap[IsNull],
			},
			{
				Operator:        NotNull,
				Name:            operatorNameMap[NotNull],
				SqlOperatorName: operatorNameMap[NotNull],
			},
		},
	},
	{
		Name:         "UUID",
		AllowedTypes: []string{"uuid.UUID"},
		AllowedPredicates: []Predicate{
			{
				Operator:        EQ,
				Name:            operatorNameMap[EQ],
				SqlOperatorName: operatorNameMap[EQ],
			},
			{
				Operator:        NEQ,
				Name:            operatorNameMap[NEQ],
				SqlOperatorName: operatorNameMap[NEQ],
			},
			{
				Operator:        In,
				Name:            operatorNameMap[In],
				SqlOperatorName: operatorNameMap[In],
			},
			{
				Operator:        NotIn,
				Name:            operatorNameMap[NotIn],
				SqlOperatorName: operatorNameMap[NotIn],
			},
			{
				Operator:        IsNull,
				Name:            operatorNameMap[IsNull],
				SqlOperatorName: operatorNameMap[IsNull],
			},
			{
				Operator:        NotNull,
				Name:            operatorNameMap[NotNull],
				SqlOperatorName: operatorNameMap[NotNull],
			},
		},
	},
	{
		Name:         "String",
		AllowedTypes: []string{"string"},
		AllowedPredicates: []Predicate{
			{
				Operator:        EQ,
				Name:            operatorNameMap[EQ],
				SqlOperatorName: operatorNameMap[EQ],
			},
			{
				Operator:        NEQ,
				Name:            operatorNameMap[NEQ],
				SqlOperatorName: operatorNameMap[NEQ],
			},
			{
				Operator:        GT,
				Name:            operatorNameMap[GT],
				SqlOperatorName: operatorNameMap[GT],
			},
			{
				Operator:        GTE,
				Name:            operatorNameMap[GTE],
				SqlOperatorName: operatorNameMap[GTE],
			},
			{
				Operator:        LT,
				Name:            operatorNameMap[LT],
				SqlOperatorName: operatorNameMap[LT],
			},
			{
				Operator:        LTE,
				Name:            operatorNameMap[LTE],
				SqlOperatorName: operatorNameMap[LTE],
			},
			{
				Operator:        In,
				Name:            operatorNameMap[In],
				SqlOperatorName: operatorNameMap[In],
			},
			{
				Operator:        NotIn,
				Name:            operatorNameMap[NotIn],
				SqlOperatorName: operatorNameMap[NotIn],
			},
			{
				Operator:        IsNull,
				Name:            operatorNameMap[IsNull],
				SqlOperatorName: operatorNameMap[IsNull],
			},
			{
				Operator:        NotNull,
				Name:            operatorNameMap[NotNull],
				SqlOperatorName: operatorNameMap[NotNull],
			},
			{
				Operator:        Contains,
				Name:            operatorNameMap[Contains],
				SqlOperatorName: operatorNameMap[Contains],
			},
			{
				Operator:        HasPrefix,
				Name:            operatorNameMap[HasPrefix],
				SqlOperatorName: operatorNameMap[HasPrefix],
			},
			{
				Operator:        HasSuffix,
				Name:            operatorNameMap[HasSuffix],
				SqlOperatorName: operatorNameMap[HasSuffix],
			},
			{
				Operator:        ContainsFold,
				Name:            operatorNameMap[ContainsFold],
				SqlOperatorName: operatorNameMap[ContainsFold],
			},
			{
				Operator:        EqualFold,
				Name:            operatorNameMap[EqualFold],
				SqlOperatorName: operatorNameMap[EqualFold],
			},
		},
	},
}

func isPredicateValid(fieldType, fieldOperator string) Operator {
	for _, fp := range FieldPredicates {
		if hasFieldType(fp.AllowedTypes, fieldType) {
			return parseOperator(fp.AllowedPredicates, fieldOperator)
		}
	}
	return Unknown
}

func parseOperator(predicates []Predicate, dest string) Operator {
	for _, p := range predicates {
		if strings.EqualFold(p.Name, dest) {
			return p.Operator
		}
	}
	return Unknown
}

func hasFieldType(types []string, dest string) bool {
	for _, t := range types {
		if t == dest {
			return true
		}
	}
	return false
}

func BuildFieldPredicate(fieldName string, fieldType string, fieldOperator string, value interface{}) func(*sql.Selector) {
	return func(s *sql.Selector) {
		if op := isPredicateValid(fieldType, fieldOperator); op != Unknown {
			switch op {
			case EQ:
				s.Where(sql.EQ(fieldName, value))
			case NEQ:
				s.Where(sql.NEQ(fieldName, value))
			case GT:
				s.Where(sql.GT(fieldName, value))
			case GTE:
				s.Where(sql.GTE(fieldName, value))
			case LT:
				s.Where(sql.LT(fieldName, value))
			case LTE:
				s.Where(sql.LTE(fieldName, value))
			case In:
				s.Where(sql.In(fieldName, value))
			case NotIn:
				s.Where(sql.NotIn(fieldName, value))
			case IsNull:
				s.Where(sql.IsNull(fieldName))
			case NotNull:
				s.Where(sql.NotNull(fieldName))
			case Contains:
				if strings.EqualFold("string", fieldType) {
					s.Where(sql.Contains(fieldName, fmt.Sprintf("%s", value)))
				}
			case HasPrefix:
				if strings.EqualFold("string", fieldType) {
					s.Where(sql.HasPrefix(fieldName, fmt.Sprintf("%s", value)))
				}
			case HasSuffix:
				if strings.EqualFold("string", fieldType) {
					s.Where(sql.HasSuffix(fieldName, fmt.Sprintf("%s", value)))
				}
			case ContainsFold:
				if strings.EqualFold("string", fieldType) {
					s.Where(sql.ContainsFold(fieldName, fmt.Sprintf("%s", value)))
				}
			case EqualFold:
				if strings.EqualFold("string", fieldType) {
					s.Where(sql.ExprP(fieldName, fmt.Sprintf("%s", value)))
				}
			default:
				s.Where(sql.EQ(fieldName, value))
			}
		}
	}
}
