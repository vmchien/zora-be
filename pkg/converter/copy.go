package converter

import (
	"encoding/json"

	"github.com/jinzhu/copier"
)

func Copy(dest any, src any) error {
	return copier.CopyWithOption(dest, src, copier.Option{
		IgnoreEmpty:   false,
		CaseSensitive: false,
		DeepCopy:      false,
		Converters: []copier.TypeConverter{
			// DTO to Proto
			mapToProtoStructConverter,
			pointerMapToProtoStructConverter,
			arrayMapToProtoListValueConverter,
			pointerArrayMapToProtoListValueConverter,
			timeToProtoTimeConverter,
			pointerTimeToProtoTimeConverter,

			// Proto to DTO
			protoListValueToArrayMapConverter,
			protoListValueToPointerArrayMapConverter,
			protoStructToMapConverter,
			protoStructToPointerMapConverter,
			protoTimeToTimeConverter,
			protoTimeToPointerTimeConverter,
			protoArrayStringToPointerArrayStringConverter,
		},
		FieldNameMapping: nil,
	})
}

func Unmarshal(dest any, src any) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &dest)
}

func UnmarshalString(dest any, jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), &dest)
}
