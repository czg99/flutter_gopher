package ffi

import (
	"errors"
	"fmt"
)

type NullString struct {
	Value string
	Valid bool // Valid is true if String is not NULL
}

func Swap(a, b int) (int, int) {
	return b, a
}

func GetNullString(value *string) NullString {
	if value == nil {
		return NullString{
			Valid: false,
		}
	}
	return NullString{
		Value: *value,
		Valid: true,
	}
}

func GetStringList(values []int) ([]string, error) {
	if len(values) == 0 {
		return nil, errors.New("values is empty")
	}
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = fmt.Sprint(v)
	}
	return result, nil
}

func StringToBytes(value string) []byte {
	return []byte(value)
}

func TestPanic() {
	panic("I can't crash")
}
