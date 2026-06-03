package utils

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

func DetectType(v string) string {
	if regexp.MustCompile(`[A-Za-z]`).MatchString(v) {
		return "code"
	}

	if len(v) >= 10 && regexp.MustCompile(`^[0-9]+$`).MatchString(v) {
		return "phone"
	}

	if regexp.MustCompile(`^[0-9]+$`).MatchString(v) {
		return "userId"
	}

	return "unknown"
}

func RemoveDuplicateStrings(input []string) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, str := range input {
		if _, ok := seen[str]; !ok {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}
	return result
}

func TryParseRouteName(routeName string) (originName, destName string) {
	rn := strings.Split(routeName, "-")
	if len(rn) < 2 {
		rn = strings.Split(routeName, "⇒")
	}

	originName = rn[0]
	destName = rn[len(rn)-1]
	return
}

func TryGetStringPointer(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// FastValueString converts common primitive values to string with low overhead.
//
// Complex objects are collapsed into a marker to avoid reflection-heavy work
// on hot paths.
func FastValueString(v any) string {
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
