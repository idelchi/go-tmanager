package env

import (
	"fmt"
	"regexp"
	"strings"
)

// Get retrieves the value associated with the given key or returns an error if the key is not found.
func (e Env) Get(key string) (string, error) {
	if v, ok := e[key]; ok {
		return v, nil
	}

	return "", fmt.Errorf("%w: %q", ErrEnvVarNotFound, key)
}

// GetOrDefault retrieves the value for the given key, or returns the provided defaultValue if the key is not found.
func (e Env) GetOrDefault(key, defaultValue string) string {
	if value, err := e.Get(key); err != nil {
		return defaultValue
	} else {
		return value
	}
}

// GetAll returns a new Env containing all key-value pairs that satisfy the given predicate function.
func (e Env) GetAll(predicate func(key, value string) bool) Env {
	result := make(Env)

	for k, v := range e {
		if predicate(k, v) {
			result[k] = v
		}
	}

	return result
}

// GetAllWithPrefix returns all environment variables with keys starting with the given prefix.
func (e Env) GetAllWithPrefix(prefix string) Env {
	return e.GetAll(func(key, _ string) bool {
		return strings.HasPrefix(key, prefix)
	})
}

// GetAllWithSuffix returns all environment variables with keys ending with the given suffix.
func (e Env) GetAllWithSuffix(suffix string) Env {
	return e.GetAll(func(key, _ string) bool {
		return strings.HasSuffix(key, suffix)
	})
}

// GetAllMatching returns all environment variables with keys matching the given regex pattern.
// It returns an error if the provided regex pattern is invalid.
func (e Env) GetAllMatching(pattern string) (Env, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	return e.GetAll(func(key, _ string) bool {
		return re.MatchString(key)
	}), nil
}

// GetAllWithValue returns all environment variables with the exact given value.
func (e Env) GetAllWithValue(value string) Env {
	return e.GetAll(func(_, v string) bool {
		return v == value
	})
}
