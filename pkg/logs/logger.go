package logs

import (
	"context"
	"fmt"
	stdlog "log"
	"os"
	"path"
	"runtime"
	"time"

	kratos_zap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	kratos_zerolog "github.com/go-kratos/kratos/contrib/log/zerolog/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"vn.vato.zora.be.api/pkg/constant"
)

var Logger log.Logger

type LOG_DRIVER string

const (
	LOG_DRIVER_ZAP      LOG_DRIVER = "zap"
	LOG_DRIVER_ZERO_LOG LOG_DRIVER = "zero-log"
	LOG_DRIVER_STD_LOG  LOG_DRIVER = "std-log"
	LOG_DRIVER_GCP      LOG_DRIVER = "gcp-log"
)

func Init(driver LOG_DRIVER) log.Logger {
	var base log.Logger
	switch driver {
	case LOG_DRIVER_ZAP:
		base = newZapLogger()
	case LOG_DRIVER_ZERO_LOG:
		base = newZerologLogger()
	case LOG_DRIVER_STD_LOG:
		base = newStdOutLogger()
	case LOG_DRIVER_GCP:
		base = NewGCPLogDriver()
	default:
		stdlog.Fatalf("unsupported log driver: %s", driver)
	}

	// TODO: disable caller in production

	l := log.With(
		base,
		"domain", domainValuer(),
		"ts", log.DefaultTimestamp,
		"caller", callerValuer(6),
		"trace.id", traceValuer(),
		"span.id", spanValuer(),
	)

	log.SetLogger(l)
	return l
}

// ---------------------
// 👇 Zerolog setup
// ---------------------
func newZerologLogger() log.Logger {
	// Optional: set time format (default is RFC3339Nano)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.LevelFieldName = "severity"
	return kratos_zerolog.NewLogger(new(zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Logger().
		Hook(gcpSeverityHook{})))
}

type gcpSeverityHook struct{}

func (g gcpSeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	switch level {
	case zerolog.DebugLevel:
		e.Str("severity", "DEBUG")
	case zerolog.InfoLevel:
		e.Str("severity", "INFO")
	case zerolog.WarnLevel:
		e.Str("severity", "WARNING")
	case zerolog.ErrorLevel:
		e.Str("severity", "ERROR")
	case zerolog.FatalLevel, zerolog.PanicLevel:
		e.Str("severity", "CRITICAL")
	default:
		e.Str("severity", "DEFAULT")
	}
}

// ---------------------
// 👇 Zap setup
// ---------------------
func newZapLogger() log.Logger {
	cfg := zap.NewDevelopmentConfig() // 👈 dùng NewProductionConfig nếu là môi trường prod
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.CallerKey = "caller"
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	cfg.EncoderConfig.LevelKey = "severity"
	cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch level {
		case zapcore.DebugLevel:
			enc.AppendString("DEBUG")
		case zapcore.InfoLevel:
			enc.AppendString("INFO")
		case zapcore.WarnLevel:
			enc.AppendString("WARNING")
		case zapcore.ErrorLevel:
			enc.AppendString("ERROR")
		case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
			enc.AppendString("CRITICAL")
		default:
			enc.AppendString("DEFAULT")
		}
	}

	zlogger, err := cfg.Build()
	if err != nil {
		stdlog.Fatalf("failed to build zap logger: %v", err)
	}

	return kratos_zap.NewLogger(zlogger)
}

// ---------------------
// 👇 STD-OUT setup
// ---------------------
func newStdOutLogger() log.Logger {
	// id, _ := os.Hostname()
	return log.With(log.NewStdLogger(os.Stdout))
}

func GetLoggerFromContext(ctx context.Context, logger log.Logger) *log.Helper {
	val := ctx.Value(constant.CTX_KEY_LOGGER)
	logger, ok := val.(log.Logger)
	if !ok {
		spanCtx := trace.SpanContextFromContext(ctx)
		traceID := spanCtx.TraceID().String()
		spanID := spanCtx.SpanID().String()

		if traceID == "" {
			traceID = "unknown"
		}
		if spanID == "" {
			spanID = "unknown"
		}
		reqLogger := log.With(logger,
			constant.CTX_KEY_LOG_TRACE_ID, traceID,
			constant.CTX_KEY_LOG_SPAN_ID, spanID,
		)
		return newHelper(reqLogger)
	}

	return newHelper(logger)
}

func traceValuer() log.Valuer {
	return func(ctx context.Context) any {
		if ctx == nil {
			return ""
		}
		if traceID, ok := ctx.Value(constant.CTX_KEY_LOG_TRACE_ID).(string); ok && traceID != "" {
			return traceID
		}
		spanCtx := trace.SpanContextFromContext(ctx)
		if spanCtx.IsValid() {
			return spanCtx.TraceID().String()
		}
		return ""
	}
}

func spanValuer() log.Valuer {
	return func(ctx context.Context) any {
		if ctx == nil {
			return ""
		}
		if spanID, ok := ctx.Value(constant.CTX_KEY_LOG_SPAN_ID).(string); ok && spanID != "" {
			return spanID
		}
		spanCtx := trace.SpanContextFromContext(ctx)
		if spanCtx.IsValid() {
			return spanCtx.SpanID().String()
		}
		return ""
	}
}

func domainValuer() log.Valuer {
	return func(ctx context.Context) any {
		if ctx == nil {
			return ""
		}
		if domain, ok := ctx.Value(constant.CTX_KEY_DOMAIN).(string); ok && domain != "" {
			return domain
		}
		return ""
	}
}

func callerValuer(skip int) log.Valuer {
	return func(ctx context.Context) any {
		// for i := 0; i < 10; i++ {
		// 	_, file, line, ok := runtime.Caller(i)
		// 	if !ok {
		// 		break
		// 	}
		// 	fmt.Printf("skip=%d -> %s:%d\n", i, path.Base(file), line)
		// }
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			return ""
		}
		return fmt.Sprintf("%s:%d", path.Base(file), line)
	}
}
