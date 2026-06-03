package format

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"
)

// defaultFuncMap defines helper functions usable inside templates.
// Extend this map as needed for your project.
func defaultFuncMap() map[string]any {
	return map[string]any{
		// {{currency .amount "VND"}} → "5,500,000 VND"
		"currency": func(v any, code string) string {
			return fmt.Sprintf("%s %s", formatNumber(v), code)
		},

		// {{ago .created_at}} → "2h ago", "3d ago"
		"ago": func(t any) string {
			ts, ok := toTime(t)
			if !ok {
				return ""
			}
			return timeAgo(ts)
		},

		"upper":    strings.ToUpper,
		"lower":    strings.ToLower,
		"trim":     strings.TrimSpace,
		"truncate": truncate,
	}
}

/* ---------------- Helper functions ---------------- */

func toTime(v any) (time.Time, bool) {
	switch x := v.(type) {
	case time.Time:
		return x, true
	case *time.Time:
		if x == nil {
			return time.Time{}, false
		}
		return *x, true
	case string:
		layouts := []string{
			time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05", "2006-01-02",
		}
		for _, l := range layouts {
			if t, err := time.Parse(l, x); err == nil {
				return t, true
			}
		}
	}
	return time.Time{}, false
}

func formatNumber(v any) string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return withThousands(rv.Int())
	case reflect.Float32, reflect.Float64:
		f := rv.Float()
		intPart, frac := math.Modf(f)
		if math.Abs(frac) < 1e-9 {
			return withThousands(int64(intPart))
		}
		return fmt.Sprintf("%s.%02d", withThousands(int64(intPart)), int(math.Round(frac*100)))
	default:
		return fmt.Sprintf("%v", v)
	}
}

func withThousands(n int64) string {
	s := fmt.Sprintf("%d", n)
	if n < 0 {
		s = s[1:]
	}
	var out []byte
	for i, c := range []byte(reverse(s)) {
		if i != 0 && i%3 == 0 {
			out = append(out, ',')
		}
		out = append(out, c)
	}
	res := reverse(string(out))
	if n < 0 {
		return "-" + res
	}
	return res
}

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func timeAgo(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	default:
		return t.Format("2006-01-02")
	}
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 || len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}
