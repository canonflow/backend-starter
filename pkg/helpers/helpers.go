package helpers

import "strings"

func Parser[E, T any](fn func(E) (T, error), arg E, def T) T {
	result, err := fn(arg)
	if err != nil {
		return def
	}

	return result
}

func SplitStringTrim(value string, sep string) []string {
	if value == "" {
		return []string{}
	}

	if sep == "" {
		sep = ","
	}

	v := strings.Split(value, sep)

	result := make([]string, len(v))

	for _, val := range v {
		result = append(result, strings.TrimSpace(val))
	}

	return result
}
