package env

import (
	"os"
	"runtime"
	"strings"
)

// Normalized returns a copy of the Env with all keys normalized to uppercase on Windows.
// On other operating systems, it returns the Env unchanged.
func (e Env) Normalized() Env {
	if runtime.GOOS == "windows" {
		return e.Normalize()
	}

	return e
}

// Normalize returns a copy of the Env with all keys converted to uppercase.
// This is primarily useful for ensuring consistent key handling on case-insensitive systems like Windows.
func (e Env) Normalize() Env {
	normalized := make(Env, len(e))

	for k, v := range e {
		normalized[strings.ToUpper(k)] = v
	}

	return normalized
}

// Expand expands environment variables in the Env values using os.ExpandEnv.
// Each value in the Env is processed, replacing any occurrences of ${var} or $var with the corresponding value from the
// environment.
func (e *Env) Expand() {
	for k, v := range *e {
		(*e)[k] = os.ExpandEnv(v)
	}
}
