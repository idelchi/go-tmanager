package env

import (
	"fmt"
	"maps"
	"strings"
)

// Add splits a `key=value` string and adds it to the Env map.
// It returns an error if the input is not properly formatted, expecting exactly one '=' separator.
func (e *Env) Add(kv string) error {
	const expectedParts = 2

	parts := strings.SplitN(kv, "=", expectedParts)
	if len(parts) != expectedParts {
		return fmt.Errorf("%w: %q", ErrEnvMalformed, kv)
	}

	(*e)[parts[0]] = parts[1]

	return nil
}

// Merge merges another Env into the current Env, without overwriting existing keys in the current Env.
// If a key in the other Env already exists in the current Env, it is not updated.
func (e *Env) Merge(envs ...Env) {
	for _, env := range envs {
		maps.Copy(env, *e)
		*e = env
	}
}

// Merged returns a new Env by merging the given Env into the current Env,
// without overwriting existing keys in the original Env.
// This method does not mutate the original Env.
func (e Env) Merged(envs ...Env) Env {
	merged := maps.Clone(e)

	for _, env := range envs {
		merged.Merge(env)
	}

	return merged
}
