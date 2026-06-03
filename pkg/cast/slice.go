package cast

import (
	"fmt"
)

func CastAnyToStringSlice[T fmt.Stringer](input []T) []string {
	var result []string
	for _, v := range input {
		result = append(result, v.String())
	}
	return result
}

func CastEnumToInt[T ~int](input []T) []int {
	result := make([]int, len(input))
	for i, v := range input {
		result[i] = int(v)
	}
	return result
}

func CastInt32ToEnum[T ~int](input []int32) []T {
	result := make([]T, len(input))
	for i, v := range input {
		result[i] = T(v)
	}
	return result
}

func CastEnumToInt32[T ~int](input []T) []int32 {
	result := make([]int32, len(input))
	for i, v := range input {
		result[i] = int32(v)
	}
	return result
}

func CopyToSliceAny[T any](slice []T) []any {
	result := make([]any, 0, len(slice))
	for _, v := range slice {
		result = append(result, v)
	}
	return result
}
