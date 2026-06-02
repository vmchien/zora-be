package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/text/unicode/norm"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _newSystemLaunchDate = time.Date(2026, 4, 1, 0, 0, 0, 0, MustVNLocation())

func init() {
	// _newSystemLaunchDate, _ = time.Parse(time.RFC3339, constant.DEPARTURE_TIME_CONDITION)
}

func Marshal(v interface{}) json.RawMessage {
	if v == nil {
		return nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return json.RawMessage(data)
}

func PtrInt(i int) *int { return &i }

func PtrInt64(i int64) *int64 { return &i }

func PtrString(s string) *string { return &s }

func PtrUUID(v uuid.UUID) *uuid.UUID {
	return &v
}

func GetString(val *wrapperspb.StringValue) string {
	if val != nil {
		return val.GetValue()
	}
	return ""
}

// Do not remove these method(s)
func ParseRawToStruct(raw json.RawMessage) (*structpb.Struct, error) {
	var s structpb.Struct
	if err := prototext.Unmarshal([]byte(raw), &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func StructToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

func StructToIntSlice(s *structpb.Struct) []int {
	if s == nil {
		return nil
	}

	result := []int{}
	for _, v := range s.Fields {
		switch kind := v.GetKind().(type) {
		case *structpb.Value_NumberValue:
			fmt.Println("String:", kind.NumberValue)
			if num, ok := v.GetKind().(*structpb.Value_NumberValue); ok {
				result = append(result, int(num.NumberValue))
			}
		case *structpb.Value_ListValue:
			/*fmt.Println("String:", kind.ListValue)*/
			if num, ok := v.GetKind().(*structpb.Value_ListValue); ok {
				// listVal := s.GetFields()["ids"].GetListValue()
				for _, v := range num.ListValue.Values {
					fmt.Println(v.GetNumberValue())
					result = append(result, int(v.GetNumberValue()))
				}
			}
		}

	}

	return result
}

func StructToSliceOfMaps(s *structpb.Struct) []map[string]interface{} {
	if s == nil {
		return nil
	}
	var list []map[string]interface{}
	for _, v := range s.Fields {
		switch kind := v.GetKind().(type) {
		case *structpb.Value_ListValue:
			/*fmt.Println("String:", kind.ListValue)*/
			if num, ok := v.GetKind().(*structpb.Value_ListValue); ok {
				for _, val := range num.ListValue.Values {
					if nested, ok := val.GetKind().(*structpb.Value_StructValue); ok {
						list = append(list, nested.StructValue.AsMap())
					}
				}
			}
		case *structpb.Value_StructValue:
			fmt.Println("String:", kind.StructValue)
			/*for i := 0; i < len(s.Fields); i++ {
				key := strconv.Itoa(i)
				if val, ok := s.Fields[key]; ok {
					if nested, ok := val.GetKind().(*structpb.Value_StructValue); ok {
						list = append(list, nested.StructValue.AsMap())
					}
				}
			}*/
		}
	}

	return list
}

// convert from ent to structpb.struct
func MapToStruct(m map[string]interface{}) *structpb.Struct {
	if m == nil {
		return nil
	}
	s, err := structpb.NewStruct(m)
	if err != nil {
		return nil
	}
	return s
}

func SliceOfMapsToStruct(m []map[string]interface{}) *structpb.Struct {
	if m == nil {
		return nil
	}

	var list []*structpb.Value
	for _, item := range m {
		s, err := structpb.NewStruct(item)
		if err != nil {
			continue
		}
		list = append(list, structpb.NewStructValue(s))
	}

	return &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"items": structpb.NewListValue(&structpb.ListValue{Values: list}),
		},
	}
}

func IntSliceToStruct(arr []int) *structpb.Struct {
	if arr == nil {
		return nil
	}

	var list []*structpb.Value
	for _, v := range arr {
		list = append(list, structpb.NewNumberValue(float64(v)))
	}

	return &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"items": structpb.NewListValue(&structpb.ListValue{Values: list}),
		},
	}
}

func StringSliceContains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// MapToJSONString converts a map[string]interface{} to JSON string
func MapToJSONString(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

// SliceMapToJSONString converts []map[string]interface{} to JSON string
func SliceMapToJSONString(arr []map[string]interface{}) string {
	if arr == nil {
		return ""
	}
	b, err := json.Marshal(arr)
	if err != nil {
		return ""
	}
	return string(b)
}

// JSONStringToMap parses a JSON string into map[string]interface{}
func JSONStringToMap(s string) map[string]interface{} {
	if s == "" {
		return nil
	}
	var m map[string]interface{}
	_ = json.Unmarshal([]byte(s), &m)
	return m
}

// JSONStringToSliceMap parses a JSON string into []map[string]interface{}
func JSONStringToSliceMap(s string) []map[string]interface{} {
	if s == "" {
		return nil
	}
	var arr []map[string]interface{}
	_ = json.Unmarshal([]byte(s), &arr)
	return arr
}

func Slugify(s string) string {
	if s == "" {
		return ""
	}

	// trim and lowercase
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	// normalize: decompose accents (NFD) then drop combining marks (Mn)
	// but first map special Vietnamese characters that decomposition won't handle:
	s = strings.ReplaceAll(s, "đ", "d")
	s = strings.ReplaceAll(s, "Đ", "d") // already lowercase, but safe

	// decompose (NFD)
	t := norm.NFD.String(s)

	// remove all combining marks
	var b strings.Builder
	b.Grow(len(t))
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			// skip combining mark
			continue
		}
		b.WriteRune(r)
	}
	s = b.String()

	// replace any sequence of non a-z0-9 characters with a single hyphen
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")

	// trim leading/trailing hyphens
	s = strings.Trim(s, "-")

	return s
}

func RemoveDuplicateKeys(input ...any) []any {
	seen := make(map[any]struct{})
	var result []any
	for _, item := range input {
		if _, ok := seen[item]; !ok {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func IsUseNewSystemByCode(code string) bool {
	if len(code) == 6 {
		return false
	}
	return true
}

func IsUseNewSystem(t *time.Time) bool {
	if t == nil || _newSystemLaunchDate.After(*t) {
		return false
	}
	return true
}

func ToMap(in interface{}) map[string]interface{} {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return out
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if field.PkgPath != "" {
			continue
		}

		key := field.Tag.Get("json")
		if key == "" || key == "-" {
			key = field.Name
		} else {
			// Tách lấy tên field từ tag (bỏ qua omitempty)
			parts := strings.Split(key, ",")
			key = parts[0]
		}

		// --- ĐOẠN XỬ LÝ QUAN TRỌNG ---
		valInterface := value.Interface()

		// Kiểm tra nếu là time.Time thì convert sang string và bỏ qua zero value
		if timeVal, ok := valInterface.(time.Time); ok {
			if !timeVal.IsZero() {
				out[key] = timeVal.Format(time.RFC3339)
			}
			// Bỏ qua zero time value (không thêm vào map)
		} else {
			out[key] = convertValue(valInterface)
		}
	}

	return out
}

// convertValue xử lý đệ quy cho struct lồng nhau hoặc slice
func convertValue(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		return ToMap(v.Interface())
	case reflect.Slice:
		res := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			res[i] = convertValue(v.Index(i).Interface())
		}
		return res
	default:
		return value
	}
}

func GetValuePointerDefault[T any](input *T, defaultValue T) T {
	if input == nil {
		return defaultValue
	}
	return *input
}

// RoundFloat làm tròn số f tới n chữ số thập phân.
// Với n = 2 → làm tròn 2 chữ số.
func RoundFloat(f float64, n int) float64 {
	pow := math.Pow10(n)
	return math.Round(f*pow) / pow
}
