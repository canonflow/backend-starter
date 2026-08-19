package helpers

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func Parser[E, T any](fn func(E) (T, error), arg E, def T) T {
	result, err := fn(arg)
	if err != nil {
		return def
	}

	return result
}

func ParseUint(s string) (uint, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(v), nil
}

func SplitStringTrim(value string, sep string) []string {
	if value == "" {
		return []string{}
	}

	if sep == "" {
		sep = ","
	}

	v := strings.Split(value, sep)

	result := make([]string, 0, len(v))

	for _, val := range v {
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}

	return result
}

func GenerateUUID() string {
	return uuid.NewString()
}

func SliceContains[T comparable](slice []T, val T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}

	return false
}

func ReverseSlice[T any](s *[]T) {
	arr := *s
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
