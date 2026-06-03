package metadata

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"vn.vato.zora.be.api/pkg/data/metadata/jsonbpgx"
)

func TestSelectQuery_WithDSLConditions(t *testing.T) {
	tests := []struct {
		name     string
		query    *SelectQuery
		expected []string
		args     []any
	}{
		{
			name: "Basic equality condition",
			query: NewSelectQuery("documents", []string{"id"}, "data").
				WhereCondition(jsonbpgx.Cond(
					jsonbpgx.FieldPath{"type"},
					jsonbpgx.OpArrowText,
					"invoice",
				)),
			expected: []string{"#>>", "->>"},
			args:     []any{"{'type'}", "invoice"},
		},
		{
			name: "IN condition with status",
			query: NewSelectQuery("documents", []string{"id"}, "data").
				WhereCondition(&jsonbpgx.JSONBCondition{
					Field: jsonbpgx.FieldPath{"status"},
					Op:    jsonbpgx.OpContainedIn,
					Value: []any{"draft", "done"},
				}),
			expected: []string{"<@"},
			args:     []any{"{'status'}", []any{"draft", "done"}},
		},
		{
			name: "LIKE title",
			query: NewSelectQuery("documents", []string{"id"}, "data").
				WhereCondition(&jsonbpgx.JSONBCondition{
					Field: jsonbpgx.FieldPath{"title"},
					Op:    jsonbpgx.OpILike,
					Value: "%report%",
				}),
			expected: []string{"ILIKE"},
			args:     []any{"{'title'}", "%report%"},
		},
		{
			name: "Ordering, Limit and Offset",
			query: NewSelectQuery("documents", []string{"id"}, "data").
				OrderByRaw("id DESC").
				Limit(10).
				Offset(20),
			expected: []string{"ORDER BY id DESC", "LIMIT", "OFFSET"},
			args:     nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := tt.query.Build()
			sqlStr, args := sel.Query()

			fmt.Println("Query: ", sqlStr)
			fmt.Println("Args: ", args)

			for _, part := range tt.expected {
				require.Contains(t, sqlStr, part)
			}
			if tt.args != nil {
				require.Equal(t, tt.args, args)
			}
		})
	}
}
