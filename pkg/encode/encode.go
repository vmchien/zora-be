package encode

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
)

func Encode(data any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TryEncode[T any](data T) []byte {
	result, err := Encode(data)
	if err != nil {
		panic(err)
	}
	return result
}

func Decode(buf []byte, data any) error {
	bb := bytes.NewBuffer(buf)
	encoder := gob.NewDecoder(bb)
	err := encoder.Decode(data)

	return err
}

func TryDecode[T any](buf []byte) *T {
	var result T

	err := Decode(buf, &result)
	if err != nil {
		panic(err)
	}

	return &result
}

func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func DecodeBase64(input string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
