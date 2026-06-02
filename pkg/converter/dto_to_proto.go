package converter

import (
	"fmt"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var mapToProtoStructConverter = copier.TypeConverter{
	SrcType: map[string]interface{}{},
	DstType: &structpb.Struct{},
	Fn: func(src interface{}) (interface{}, error) {
		if m, ok := src.(map[string]interface{}); ok {
			return structpb.NewStruct(m)
		}
		return nil, fmt.Errorf("invalid type, expected map[string]interface{}")
	},
}
var pointerMapToProtoStructConverter = copier.TypeConverter{
	SrcType: &map[string]interface{}{},
	DstType: &structpb.Struct{},
	Fn: func(src interface{}) (interface{}, error) {
		if m, ok := src.(*map[string]interface{}); ok {
			if m == nil {
				return structpb.NewStruct(nil)
			}
			return structpb.NewStruct(*m)
		}
		return nil, fmt.Errorf("invalid type, expected *map[string]interface{}")
	},
}
var arrayMapToProtoListValueConverter = copier.TypeConverter{
	SrcType: []map[string]interface{}{},
	DstType: &structpb.ListValue{},
	Fn: func(src interface{}) (interface{}, error) {
		if m, ok := src.([]map[string]interface{}); ok {
			return _arrayMapToList(m)
		}
		return nil, fmt.Errorf("invalid type, expected map[string]interface{}")
	},
}

var pointerArrayMapToProtoListValueConverter = copier.TypeConverter{
	SrcType: &[]map[string]interface{}{},
	DstType: &structpb.ListValue{},
	Fn: func(src interface{}) (interface{}, error) {
		if m, ok := src.(*[]map[string]interface{}); ok {
			if m == nil {
				return structpb.NewListValue(nil), nil
			}
			return _arrayMapToList(*m)
		}
		return nil, fmt.Errorf("invalid type, expected map[string]interface{}")
	},
}
var timeToProtoTimeConverter = copier.TypeConverter{
	SrcType: time.Time{},
	DstType: &timestamppb.Timestamp{},
	Fn: func(src interface{}) (interface{}, error) {
		if m, ok := src.(time.Time); ok {
			return timestamppb.New(m), nil
		}
		return nil, fmt.Errorf("invalid type, expected time.Time{}")
	},
}
var pointerTimeToProtoTimeConverter = copier.TypeConverter{
	SrcType: &time.Time{},
	DstType: &timestamppb.Timestamp{},
	Fn: func(src interface{}) (interface{}, error) {
		if m, ok := src.(*time.Time); ok {
			if m == nil {
				return timestamppb.New(time.Time{}), nil
			}
			return timestamppb.New(*m), nil
		}
		return nil, fmt.Errorf("invalid type, expected time.Time{}")
	},
}

func _arrayMapToList(src []map[string]interface{}) (*structpb.ListValue, error) {
	var rs structpb.ListValue
	for _, v := range src {
		val, err := structpb.NewValue(v)
		if err != nil {
			return nil, err
		}
		rs.Values = append(rs.Values, val)
	}
	return &rs, nil
}
