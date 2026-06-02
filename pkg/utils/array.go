package utils

import (
	"reflect"
	"strconv"
	"strings"
)

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func NotContainsArrayInArray[T comparable](smaller []T, bigger []T) bool {
	for _, v := range smaller {
		ok := Contains(bigger, v)
		if !ok {
			return true
		}
	}
	return false
}

func Compare2Array[T comparable](A, B []T) (inAnotB, inBnotA, inAB []T) {
	m := make(map[T]uint8)
	for _, val := range A {
		m[val] |= (1 << 0)
	}
	for _, val := range B {
		m[val] |= (1 << 1)
	}
	inAnotB = []T{}
	inBnotA = []T{}
	inAB = []T{}
	for i, v := range m {
		a := v&(1<<0) != 0
		b := v&(1<<1) != 0
		switch {
		case a && b:
			inAB = append(inAB, i)
		case a && !b:
			inAnotB = append(inAnotB, i)
		case !a && b:
			inBnotA = append(inBnotA, i)
		}
	}
	return
}

// input string ex: "1,2,3,4"
func ConvertStringToUint32(s string) []uint32 {
	var listBranch []uint32
	str := strings.ReplaceAll(s, " ", "")
	arr := strings.Split(str, ",")
	for _, v := range arr {
		val, err := strconv.Atoi(v)
		if err == nil {
			listBranch = append(listBranch, uint32(val))
		}
	}
	return listBranch
}

func ArrayToMapId[T any, K comparable](data []T, keyField string) map[K][]T {
	res := make(map[K][]T)
	for i := 0; i < len(data); i++ {
		indexVal := data[i]
		var idValue reflect.Value
		value := reflect.ValueOf(indexVal)
		if value.Kind() == reflect.Ptr {
			value = reflect.ValueOf(value.Interface()).Elem()
		}
		idValue = value.FieldByName(keyField)
		if idValue.IsValid() {
			id := idValue.Interface().(K)
			res[id] = append(res[id], indexVal)
		}
	}
	return res
}

func DistinctArray[T comparable](data []T) []T {
	lenght := len(data)
	check := make(map[T]bool, lenght)
	res := []T{}

	for i := 0; i < lenght; i++ {
		if val, ok := check[data[i]]; !(val && ok) {
			res = append(res, data[i])
			check[data[i]] = true
		}
	}

	return res
}

type UniqueArray[T comparable] struct {
	Value   []T
	Checker map[T]bool
}

func NewUniqueArray[T comparable](value []T) (sua UniqueArray[T]) {
	sua = UniqueArray[T]{}
	sua.Value = value
	sua.Checker = make(map[T]bool)
	return
}

func (sua *UniqueArray[T]) Append(val ...T) (err error) {
	for _, v := range val {
		if oldval, already := sua.Checker[v]; !(oldval && already) {
			sua.Value = append(sua.Value, v)
			sua.Checker[v] = true
		}
	}
	return
}

func GetAllFieldString[T comparable](obj T) []string {
	v := reflect.ValueOf(obj).Elem()
	fields := []string{}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String && field.String() != "" {
			fields = append(fields, field.String())
		}
	}
	return fields
}
