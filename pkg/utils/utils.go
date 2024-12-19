package utils

import (
	"net/url"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// IsURL checks if the input string is a valid URL.
func IsURL(str string) bool {
	u, err := url.Parse(str)

	return err == nil && u.Scheme != "" && u.Host != ""
}

// IsZeroValue checks if the input value is its zero value.
func IsZeroValue[S comparable](input S) bool {
	var zero S

	return input == zero
}

// SetIfZeroValue sets the value of input to the specified value if it is its zero value.
func SetIfZeroValue[S comparable](input *S, value S) {
	if IsZeroValue(*input) {
		*input = value
	}
}

// SetSliceIfNil sets the value of input to the provided values slice if input is nil.
func SetSliceIfNil[S ~[]T, T any](input *S, values ...T) {
	if *input == nil {
		*input = append([]T(nil), values...)
	}
}

// SetMapIfNil sets the value of input to the provided defaultMap if input is nil.
func SetMapIfNil[M ~map[K]V, K comparable, V any](input *M, values M) {
	if *input == nil {
		*input = make(M, len(values))
		for k, v := range values {
			(*input)[k] = v
		}
	}
}

// NormalizeMap normalizes the keys of a map to title case recursively.
// If the value is another map, it will recursively normalize its keys as well.
func NormalizeMap(m map[string]any) map[string]any {
	normalizedMap := make(map[string]any)
	c := cases.Title(language.English)

	for key, value := range m {
		upperKey := c.String(key)

		switch v := value.(type) {
		case map[string]any:
			normalizedMap[upperKey] = NormalizeMap(v)
		default:
			normalizedMap[upperKey] = v
		}
	}

	return normalizedMap
}

// DeepMergeMapsWithoutOverwrite merges two maps of map[string]any, adding values
// from second to first without overwriting the existing values in first.
// It performs a deep merge, handling nested maps recursively.
func DeepMergeMapsWithoutOverwrite(first, second map[string]any) {
	for key, secondVal := range second {
		if firstVal, exists := first[key]; exists {
			// If both values are maps, recursively merge them
			if firstMap, ok1 := firstVal.(map[string]any); ok1 {
				if secondMap, ok2 := secondVal.(map[string]any); ok2 {
					DeepMergeMapsWithoutOverwrite(firstMap, secondMap)
				}
			} // If the key exists but isn't a map, do nothing (keep the original value)
		} else {
			// If the key doesn't exist in first, add it from second
			first[key] = secondVal
		}
	}
}
