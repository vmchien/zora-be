package orm

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"vn.vato.zora.be.api/pkg/logs"
)

type ctxKey struct{}

type tracer struct {
	l *logs.Helper
}

type queryInfo struct {
	sql  string
	args []any
	t0   time.Time
}

func (t tracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	qi := queryInfo{
		sql:  data.SQL,
		args: data.Args,
		t0:   time.Now(),
	}
	return context.WithValue(ctx, ctxKey{}, qi)
}

func (t tracer) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	if v := ctx.Value(ctxKey{}); v != nil {
		qi := v.(queryInfo)
		t.writeLog(ctx, qi, data)
		// t.l.Debugf("[pgx] took=%s | sql=%s | args=%v | tag=%s | err=%v", time.Since(qi.t0), qi.sql, qi.args, data.CommandTag, data.Err)
		// fmt.Printf("[pgx] took=%s | sql=%s | args=%v | tag=%s | err=%v\n",
		// 	time.Since(qi.t0), qi.sql, qi.args, data.CommandTag, data.Err)
	} else {
		t.l.Debugf(ctx, "[pgx] took=? | tag=%s | err=%v\n", data.CommandTag, data.Err)
		// fmt.Printf("[pgx] took=? | tag=%s | err=%v\n", data.CommandTag, data.Err)
	}
}

func (t tracer) writeLog(ctx context.Context, qi queryInfo, data pgx.TraceQueryEndData) {
	// color := "\033[32m"
	// elapsed := time.Since(qi.t0)
	// if elapsed > time.Second {
	// 	color = "\033[31m"
	// } else if elapsed > 200*time.Millisecond {
	// 	color = "\033[33m"
	// }
	// coloredTime := fmt.Sprintf("%s[pgx] took=%s %s | tag=%s | err=%v", color, time.Since(qi.t0), "\033[0m", data.CommandTag, data.Err)
	coloredTime := fmt.Sprintf("[pgx] took=%s | tag=%s | err=%v", time.Since(qi.t0), data.CommandTag, data.Err)
	t.l.Debugf(ctx, coloredTime)
}
