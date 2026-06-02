package orm

/*
PostgreSQL Error Code Reference (SQLSTATE)

Danh sách các lỗi phổ biến thường gặp khi làm việc với PostgreSQL
đặc biệt hữu ích khi sử dụng driver pgx (qua pgx/stdlib) và Ent:

| Tên lỗi                         | SQLSTATE | Ý nghĩa                                                                |
|--------------------------------|----------|------------------------------------------------------------------------|
| Unique violation               | 23505    | Vi phạm unique constraint (e.g. email, username)                      |
| Foreign key violation          | 23503    | Vi phạm khóa ngoại (khóa không tồn tại)                              |
| Not null violation             | 23502    | Cột bị rỗng khi đã khai báo NOT NULL                                 |
| Check constraint fail          | 23514    | Vi phạm điều kiện CHECK constraint                                   |
| Exclusion violation            | 23P01    | Vi phạm exclusion constraint                                          |
| Invalid text representation    | 22P02    | Lỗi parse kiểu dữ liệu (ví dụ nhập chữ vào cột kiểu số)              |
| String data too long           | 22001    | Chuỗi vượt quá độ dài tối đa của cột                                  |
| Numeric value out of range     | 22003    | Giá trị số vượt quá giới hạn (int16/int32/int64)                      |

Nguồn: PostgreSQL official docs + pgx/pgconn.PgError.Code

Có thể dùng để viết các hàm như:
- IsUniqueViolation(err)
- IsForeignKeyViolation(err)
- IsNotNullViolation(err)
...
*/

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func extractPgError(err error) (*pgconn.PgError, bool) {
	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	return pgErr, ok
}

func IsUniqueViolation(err error) bool {
	if e, ok := extractPgError(err); ok {
		return e.Code == "23505"
	}
	return false
}

func IsForeignKeyViolation(err error) bool {
	if e, ok := extractPgError(err); ok {
		return e.Code == "23503"
	}
	return false
}

func IsNotNullViolation(err error) bool {
	if e, ok := extractPgError(err); ok {
		return e.Code == "23502"
	}
	return false
}

func IsCheckViolation(err error) bool {
	if e, ok := extractPgError(err); ok {
		return e.Code == "23514"
	}
	return false
}

func IsStringTooLong(err error) bool {
	if e, ok := extractPgError(err); ok {
		return e.Code == "22001"
	}
	return false
}
