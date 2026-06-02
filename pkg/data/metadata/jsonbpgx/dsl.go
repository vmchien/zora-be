package jsonbpgx

import (
	"encoding/json"
	"fmt"
	"strings"
)

// --- Interfaces ---

type Condition interface {
	ToSQL(paramIndex int, jsonbCol string) (clause string, args []any, nextParamIndex int)
	Validate() error
	Suggest() string
}

type Update interface {
	ToSQL(paramIndex int, jsonbCol string) (clause string, args []any, nextParamIndex int)
	Validate() error
	Suggest() string
}

// --- DSL Types ---

type FieldPath []string

func (fp FieldPath) String() string {
	parts := make([]string, len(fp))
	for i, p := range fp {
		parts[i] = fmt.Sprintf("'%s'", p)
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ","))
}

// --- Conditions ---

type JSONBCondition struct {
	Field FieldPath
	Op    JSONBOperator
	Value any
}

func (c *JSONBCondition) ToSQL(idx int, jsonbCol string) (string, []any, int) {
	pathExpr := fmt.Sprintf("%s #> $%d", jsonbCol, idx)
	if c.Op == OpArrowText || c.Op == OpHashText || c.Op == OpLike || c.Op == OpILike || c.Op == OpNotLike {
		pathExpr = fmt.Sprintf("%s #>> $%d", jsonbCol, idx)
	}
	switch c.Op {
	case OpLike, OpILike, OpNotLike:
		return fmt.Sprintf("%s %s $%d", pathExpr, c.Op, idx+1),
			[]any{c.Field.String(), c.Value}, idx + 2
	case OpIn, OpNotIn:
		return fmt.Sprintf("%s %s (%s)", pathExpr, c.Op, placeholders(idx+1, c.Value)),
			append([]any{c.Field.String()}, c.Value.([]any)...), idx + 1 + len(c.Value.([]any))
	case OpExists:
		return fmt.Sprintf("jsonb_path_exists(%s, $%d)", jsonbCol, idx),
			[]any{c.Value}, idx + 1
	default:
		clause := fmt.Sprintf("%s %s $%d", pathExpr, c.Op, idx+1)
		return clause, []any{c.Field.String(), c.Value}, idx + 2
	}
}

func (c *JSONBCondition) Validate() error {
	if !IsValidOperator(c.Op) {
		return fmt.Errorf("invalid operator: %s", c.Op)
	}
	return nil
}

func (c *JSONBCondition) Suggest() string {
	return SuggestBehaviorForOperator(c.Op)
}

func Cond(field FieldPath, op JSONBOperator, val any) Condition {
	return &JSONBCondition{Field: field, Op: op, Value: val}
}

// --- Updates ---

type JSONBSetField struct {
	Path  FieldPath
	Value any
}

func (s *JSONBSetField) ToSQL(idx int, jsonbCol string) (string, []any, int) {
	jsonVal, _ := json.Marshal(s.Value)
	clause := fmt.Sprintf("%s = jsonb_set(%s, $%d, $%d::jsonb, true)", jsonbCol, jsonbCol, idx, idx+1)
	return clause, []any{s.Path.String(), jsonVal}, idx + 2
}

func (s *JSONBSetField) Validate() error {
	return nil // Add rules if needed
}

func (s *JSONBSetField) Suggest() string {
	return SuggestBehaviorForFunction(FnJsonbSet)
}

func Set(path FieldPath, value any) Update {
	return &JSONBSetField{Path: path, Value: value}
}

type JSONBMergePatch struct {
	Patch map[string]any
}

func (m *JSONBMergePatch) ToSQL(idx int, jsonbCol string) (string, []any, int) {
	j, _ := json.Marshal(m.Patch)
	clause := fmt.Sprintf("%s = %s || $%d::jsonb", jsonbCol, jsonbCol, idx)
	return clause, []any{j}, idx + 1
}

func (m *JSONBMergePatch) Validate() error {
	return nil
}

func (m *JSONBMergePatch) Suggest() string {
	return SuggestBehaviorForOperator(OpConcat)
}

func Merge(patch map[string]any) Update {
	return &JSONBMergePatch{Patch: patch}
}

// --- QueryBuilder ---

type QueryBuilder struct {
	table    string
	jsonbCol string
	conds    []Condition
	order    string
	limit    int
	offset   int
}

func NewQueryBuilder(table string, jsonbCol string) *QueryBuilder {
	return &QueryBuilder{table: table, jsonbCol: jsonbCol, limit: 1, offset: 0}
}

func (qb *QueryBuilder) Where(c Condition) *QueryBuilder {
	qb.conds = append(qb.conds, c)
	return qb
}

func (qb *QueryBuilder) OrderBy(field string, asc bool) *QueryBuilder {
	dir := "ASC"
	if !asc {
		dir = "DESC"
	}
	qb.order = fmt.Sprintf("ORDER BY %s %s", field, dir)
	return qb
}

func (qb *QueryBuilder) Limit(n int) *QueryBuilder {
	qb.limit = n
	return qb
}

func (qb *QueryBuilder) Offset(n int) *QueryBuilder {
	qb.offset = n
	return qb
}

func (qb *QueryBuilder) Build() (string, []any, error) {
	args := []any{}
	idx := 1
	where := []string{}
	for _, cond := range qb.conds {
		if err := cond.Validate(); err != nil {
			return "", nil, err
		}
		s, a, i := cond.ToSQL(idx, qb.jsonbCol)
		where = append(where, s)
		args = append(args, a...)
		idx = i
	}

	sql := fmt.Sprintf("SELECT * FROM %s", qb.table)
	if len(where) > 0 {
		sql += " WHERE " + strings.Join(where, " AND ")
	}
	if qb.order != "" {
		sql += " " + qb.order
	}
	if qb.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", qb.limit)
	}
	if qb.offset > 0 {
		sql += fmt.Sprintf(" OFFSET %d", qb.offset)
	}
	return sql, args, nil
}

// --- UpdateBuilder ---

type UpdateBuilder struct {
	table    string
	jsonbCol string
	update   Update
	conds    []Condition
}

func NewUpdateBuilder(table string, jsonbCol string) *UpdateBuilder {
	return &UpdateBuilder{table: table, jsonbCol: jsonbCol}
}

func (ub *UpdateBuilder) Set(up Update) *UpdateBuilder {
	ub.update = up
	return ub
}

func (ub *UpdateBuilder) Where(c Condition) *UpdateBuilder {
	ub.conds = append(ub.conds, c)
	return ub
}

func (ub *UpdateBuilder) Build() (string, []any, error) {
	args := []any{}
	idx := 1

	if err := ub.update.Validate(); err != nil {
		return "", nil, err
	}
	setSQL, setArgs, newIdx := ub.update.ToSQL(idx, ub.jsonbCol)
	args = append(args, setArgs...)
	idx = newIdx

	where := []string{}
	for _, cond := range ub.conds {
		if err := cond.Validate(); err != nil {
			return "", nil, err
		}
		s, a, i := cond.ToSQL(idx, ub.jsonbCol)
		where = append(where, s)
		args = append(args, a...)
		idx = i
	}

	sql := fmt.Sprintf("UPDATE %s SET %s", ub.table, setSQL)
	if len(where) > 0 {
		sql += " WHERE " + strings.Join(where, " AND ")
	}
	return sql, args, nil
}

func placeholders(start int, values any) string {
	vals, ok := values.([]any)
	if !ok {
		return "$" + fmt.Sprint(start)
	}
	ph := make([]string, len(vals))
	for i := range vals {
		ph[i] = fmt.Sprintf("$%d", start+i)
	}
	return strings.Join(ph, ", ")
}
