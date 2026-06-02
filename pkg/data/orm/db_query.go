package orm

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
)

// ScanRows executes a raw SQL query and scans each row into *T using the provided
// scan function. This avoids reflection entirely and lets the caller decide how
// to map columns to fields.
//
// Usage example:
//
//	type RolePerm struct {
//	    RoleID   uuid.UUID
//	    RoleName string
//	    PermID   uuid.UUID
//	    Resource string
//	    Action   string
//	}
//
//	out, err := db.ScanRows[RolePerm](
//	    ctx, myDB,
//	    `SELECT r.id, r.name, p.id, p.resource, p.action
//	       FROM smp_iam__rbac_roles r
//	       JOIN smp_iam__rbac_role_permissions rp ON rp.role_id = r.id
//	       JOIN smp_iam__rbac_permissions p       ON p.id = rp.permission_id
//	      WHERE r.tenant_id = $1`,
//	    func(rows *sql.Rows, dst *RolePerm) error {
//	        return rows.Scan(&dst.RoleID, &dst.RoleName, &dst.PermID, &dst.Resource, &dst.Action)
//	    },
//	    tenantID,
//	)
//
// Parameters:
//   - d: wraps your *sql.DB (or compatible) connection. It must expose DB.QueryContext.
//   - scan: a callback that reads the current row into *T via rows.Scan(...) and returns an error if any.
//
// Returns:
//   - A []T slice with one element per row, or an error.
//
// Notes:
//   - T can be any type. Typically it's a struct DTO.
//   - This function does not use reflection and therefore allocates minimally.
//   - Caller is responsible for supplying the exact Scan(...) destinations/order to match the query.
func ScanRows[T any](ctx context.Context, conn *Connection, rawQuery string, args []any, scanFn func(*sql.Rows, *T) error) ([]T, error) {
	return scan[T](ctx, conn, rawQuery, args, scanFn)
}

// ScanRow executes a raw SQL query expected to return a single row and scans it into *T
// using the provided scan function. This avoids reflection entirely and lets the caller
// decide how to map columns to fields.
//
// Usage example:
//
//	type UserProfile struct {
//	    ID        uuid.UUID
//	    Username  string
//	    Email     string
//	    CreatedAt time.Time
//	}
//
//	out, err := db.ScanRow[UserProfile](
//	    ctx, myDB,
//	    `SELECT id, username, email, created_at
//	       FROM users
//	      WHERE id = $1`,
//	    func(rows *sql.Rows, dst *UserProfile) error {
//	        return rows.Scan(&dst.ID, &dst.Username, &dst.Email, &dst.CreatedAt)
//	    },
//	    userID,
//	)
//
// Parameters:
//   - d: wraps your *sql.DB (or compatible) connection. It must expose DB.QueryContext.
//   - scan: a callback that reads the current row into *T via rows.Scan(...) and returns an error if any.
//
// Returns:
//
//   - A pointer to T with the scanned data, or sql.ErrNoRows if no rows were returned.
//
// Notes:
//   - T can be any type. Typically it's a struct DTO.
//   - This function does not use reflection and therefore allocates minimally.
//   - Caller is responsible for supplying the exact Scan(...) destinations/order to match the query.
func ScanRow[T any](ctx context.Context, conn *Connection, query string, scanFn func(*sql.Rows, *T) error, args []any) (*T, error) {
	return scanOne[T](ctx, conn, query, args, scanFn)
}

// ScanCount executes a raw SQL COUNT(*) query and returns the total as int64.
// The query must return exactly one row with one numeric column.
func ScanCount(ctx context.Context, conn *Connection, rawQuery string, args ...any) (int32, error) {
	if conn == nil || conn.Driver == nil {
		return 0, fmt.Errorf("db: nil connection or driver")
	}

	var rows sql.Rows
	if err := conn.Driver.Query(ctx, rawQuery, args, &rows); err != nil {
		return 0, fmt.Errorf("db: rawQuery failed: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, fmt.Errorf("sql: no rows in result set")
	}

	var total int32
	if err := rows.Scan(&total); err != nil {
		return 0, fmt.Errorf("db: scan count: %w", err)
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("db: rows iteration: %w", err)
	}

	return total, nil
}

func scan[T any](
	ctx context.Context,
	conn *Connection,
	rawQuery string,
	args []any,
	scanFn func(*sql.Rows, *T) error,
) ([]T, error) {
	if conn == nil || conn.Driver == nil {
		return nil, fmt.Errorf("db: nil connection or driver")
	}

	var rows sql.Rows
	if err := conn.Driver.Query(ctx, rawQuery, args, &rows); err != nil {
		return nil, fmt.Errorf("db: rawQuery failed: %w", err)
	}
	// if rows == nil {
	// 	return nil, fmt.Errorf("db: driver returned nil *sql.Rows")
	// }
	defer rows.Close()

	out := make([]T, 0)
	for rows.Next() {
		var v T
		if err := scanFn(&rows, &v); err != nil {
			return nil, fmt.Errorf("db: scan row: %w", err)
		}
		out = append(out, v)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("db: rows iteration: %w", err)
	}

	return out, nil
}

func scanOne[T any](
	ctx context.Context,
	conn *Connection,
	query string,
	args []any,
	scanFn func(*sql.Rows, *T) error,
) (*T, error) {
	items, err := scan[T](ctx, conn, query, args, scanFn)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("sql: no rows in result set")
	}
	return &items[0], nil
}
