package metadata

import (
	"entgo.io/ent/dialect/sql"
	"vn.vato.zora.be.api/pkg/data/metadata/jsonbpgx"
)

type SelectQuery struct {
	table    string
	columns  []string
	jsonbCol string
	preds    []func(*sql.Selector)
	limit    int
	offset   int
	order    string
}

func NewSelectQuery(table string, columns []string, jsonbCol string) *SelectQuery {
	return &SelectQuery{
		table:    table,
		columns:  columns,
		jsonbCol: jsonbCol,
		limit:    100,
	}
}

func (s *SelectQuery) Where(pred func(*sql.Selector)) *SelectQuery {
	s.preds = append(s.preds, pred)
	return s
}

func (s *SelectQuery) WhereCondition(cond jsonbpgx.Condition) *SelectQuery {
	return s.Where(func(sel *sql.Selector) {
		if err := cond.Validate(); err != nil {
			panic(err)
		}
		sqlClause, args, _ := cond.ToSQL(1, s.jsonbCol)
		sel.Where(sql.ExprP(sqlClause, args...))
	})
}

func (s *SelectQuery) Limit(n int) *SelectQuery {
	s.limit = n
	return s
}

func (s *SelectQuery) Offset(n int) *SelectQuery {
	s.offset = n
	return s
}

func (s *SelectQuery) OrderByRaw(sqlClause string) *SelectQuery {
	s.order = sqlClause
	return s
}

func (s *SelectQuery) Build() *sql.Selector {
	builder := sql.Select(s.columns...).From(sql.Table(s.table))
	for _, pred := range s.preds {
		pred(builder)
	}
	if s.order != "" {
		builder.OrderExpr(sql.Expr(s.order))
	}
	if s.limit > 0 {
		builder.Limit(s.limit)
	}
	if s.offset > 0 {
		builder.Offset(s.offset)
	}
	return builder
}
