package logs

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
)

// gcpLogPayload is the minimal JSON payload written to stdout.
type gcpLogPayload struct {
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// GCPLogDriver is a minimal, allocation-conscious Kratos logger for GCP.
//
// Output shape:
//   - severity
//   - message
//
// All extra key/value pairs are flattened into message.
type GCPLogDriver struct {
	out *os.File
}

var byteBufferPool = sync.Pool{
	New: func() any {
		b := make([]byte, 0, 256)
		return &b
	},
}

// NewGCPLogDriver creates a new ultra-lean stdout logger for Kratos.
func NewGCPLogDriver() log.Logger {
	return &GCPLogDriver{
		out: os.Stdout,
	}
}

// Log implements the Kratos log.Logger interface.
func (l *GCPLogDriver) Log(level log.Level, keyvals ...any) error {
	bufPtr := byteBufferPool.Get().(*[]byte)
	buf := (*bufPtr)[:0]

	buf = buildLogMessageZeroPrepend(buf, keyvals)

	payload := gcpLogPayload{
		Severity: toGCPSeverity(level),
		Message:  string(buf),
	}

	data, err := json.Marshal(payload)

	*bufPtr = buf[:0]
	byteBufferPool.Put(bufPtr)

	if err != nil {
		return err
	}

	_, err = l.out.Write(append(data, '\n'))
	return err
}

// buildLogMessageZeroPrepend builds the final message in a single pass.
//
// Rules:
//   - "msg" is written as plain text
//   - "level" and "severity" are ignored
//   - all other fields are flattened as "key=value"
//   - output order follows the original keyvals order
//
// Example:
//
//	"msg","request failed","trace.id","abc","user_id",123
//
// Result:
//
//	"request failed trace.id=abc user_id=123"
func buildLogMessageZeroPrepend(dst []byte, keyvals []any) []byte {
	wroteAny := false

	for i := 0; i < len(keyvals); i += 2 {
		key := fastString(keyvals[i])
		if key == "" {
			continue
		}

		var value any
		if i+1 < len(keyvals) {
			value = keyvals[i+1]
		}

		switch key {
		case "level", "severity":
			continue
		case "msg":
			msg := fastValueString(value)
			if msg == "" {
				continue
			}
			if wroteAny {
				dst = append(dst, ' ')
			}
			dst = append(dst, msg...)
			wroteAny = true
		default:
			if wroteAny {
				dst = append(dst, ' ')
			}
			dst = append(dst, key...)
			dst = append(dst, '=')
			dst = append(dst, fastValueString(value)...)
			wroteAny = true
		}
	}

	return dst
}

// fastValueString converts common primitive values to string with low overhead.
//
// Complex objects are collapsed into a marker to avoid reflection-heavy work
// on hot paths.
func fastValueString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case []byte:
		return string(x)
	case error:
		return x.Error()
	case bool:
		if x {
			return "true"
		}
		return "false"
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(x, 'f', -1, 64)
	case json.Number:
		return x.String()
	default:
		return "[complex]"
	}
}

// fastString converts a key or simple value to string cheaply.
func fastString(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case []byte:
		return string(x)
	case error:
		return x.Error()
	case int:
		return strconv.Itoa(x)
	case int8:
		return strconv.FormatInt(int64(x), 10)
	case int16:
		return strconv.FormatInt(int64(x), 10)
	case int32:
		return strconv.FormatInt(int64(x), 10)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint8:
		return strconv.FormatUint(uint64(x), 10)
	case uint16:
		return strconv.FormatUint(uint64(x), 10)
	case uint32:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	default:
		return ""
	}
}

// toGCPSeverity maps Kratos levels to Google Cloud Logging severity values.
func toGCPSeverity(level log.Level) string {
	switch level {
	case log.LevelDebug:
		return "DEBUG"
	case log.LevelInfo:
		return "INFO"
	case log.LevelWarn:
		return "WARNING"
	case log.LevelError:
		return "ERROR"
	case log.LevelFatal:
		return "CRITICAL"
	default:
		return "DEFAULT"
	}
}
