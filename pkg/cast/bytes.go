package cast

import (
	"encoding/base64"
	"encoding/json"
)

func BytesToBase64URLBytes(buf []byte) []byte {
	encoder := base64.URLEncoding
	allocLen := encoder.EncodedLen(len(buf))

	result := make([]byte, allocLen)
	encoder.Encode(result, buf)

	return result
}

func ToBytes(val any) []byte {
	marshal, err := json.Marshal(val)
	if err != nil {
		return nil
	}
	return marshal
}
