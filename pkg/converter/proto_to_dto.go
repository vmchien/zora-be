package converter

import (
	"fmt"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var protoListValueToArrayMapConverter = copier.TypeConverter{
	SrcType: &structpb.ListValue{},
	DstType: []map[string]interface{}{},
	Fn: func(src interface{}) (interface{}, error) {
		if lv, ok := src.(*structpb.ListValue); ok {
			rs, err := _listToArrayMap(lv)
			if rs == nil {
				return nil, nil
			}
			return rs, err
		}
		return nil, fmt.Errorf("invalid type, expected *structpb.ListValue{}")
	},
}
var protoListValueToPointerArrayMapConverter = copier.TypeConverter{
	SrcType: &structpb.ListValue{},
	DstType: &[]map[string]interface{}{},
	Fn: func(src interface{}) (interface{}, error) {
		if lv, ok := src.(*structpb.ListValue); ok {
			rs, err := _listToArrayMap(lv)
			if rs == nil {
				return nil, nil
			}
			return &rs, err
		}
		return nil, fmt.Errorf("invalid type, expected *structpb.ListValue{}")
	},
}
var protoStructToMapConverter = copier.TypeConverter{
	SrcType: &structpb.Struct{},
	DstType: map[string]interface{}{},
	Fn: func(src interface{}) (interface{}, error) {
		if s, ok := src.(*structpb.Struct); ok {
			if s == nil {
				return nil, nil
			}
			return s.AsMap(), nil
		}
		return nil, fmt.Errorf("invalid type, expected *structpb.Struct{}")
	},
}
var protoStructToPointerMapConverter = copier.TypeConverter{
	SrcType: &structpb.Struct{},
	DstType: &map[string]interface{}{},
	Fn: func(src interface{}) (interface{}, error) {
		if s, ok := src.(*structpb.Struct); ok {
			if s == nil {
				return nil, nil
			}
			m := s.AsMap()
			return &m, nil
		}
		return nil, fmt.Errorf("invalid type, expected *structpb.Struct{}")
	},
}
var protoTimeToTimeConverter = copier.TypeConverter{
	SrcType: &timestamppb.Timestamp{},
	DstType: time.Time{},
	Fn: func(src interface{}) (interface{}, error) {
		if m, ok := src.(*timestamppb.Timestamp); ok {
			return m.AsTime(), nil
		}
		return nil, fmt.Errorf("invalid type, expected timestamppb.Timestamp{}")
	},
}
var protoTimeToPointerTimeConverter = copier.TypeConverter{
	SrcType: &timestamppb.Timestamp{},
	DstType: &time.Time{},
	Fn: func(src interface{}) (interface{}, error) {
		if m, ok := src.(*timestamppb.Timestamp); ok {
			t := m.AsTime()
			return &t, nil
		}
		return nil, fmt.Errorf("invalid type, expected timestamppb.Timestamp{}")
	},
}
var protoArrayStringToPointerArrayStringConverter = copier.TypeConverter{
	SrcType: []string{},
	DstType: &[]string{},
	Fn: func(src interface{}) (interface{}, error) {
		if m, ok := src.([]string); ok {
			if len(m) == 0 {
				return nil, nil
			}
			return &m, nil
		}
		return nil, fmt.Errorf("invalid type, expected []string{}")
	},
}

func _listToArrayMap(src *structpb.ListValue) ([]map[string]interface{}, error) {
	if src == nil {
		return nil, nil
	}
	var rs []map[string]interface{}
	for _, v := range src.Values {
		if v == nil {
			continue
		}
		if sv := v.GetStructValue(); sv != nil {
			rs = append(rs, sv.AsMap())
			continue
		}
		// Fallback: attempt to coerce via AsInterface()
		if mv, ok := v.AsInterface().(map[string]interface{}); ok {
			rs = append(rs, mv)
		}
		// Non-object items are ignored by design; adjust if needed
	}
	return rs, nil
}
