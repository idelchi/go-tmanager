package rusti

import (
	"fmt"
	"path/filepath"
)

type Env map[string]string

func (e Env) ToSlice() []string {
	var env []string
	for k, v := range e {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

func (e *Env) Append(env Env) {
	for k, v := range env {
		(*e)[k] = v
	}
}

func (e *Env) Default(dir string) {
	*e = Env{
		"CARGO_HOME":  filepath.Join(dir, ".cargo"),
		"RUSTUP_HOME": filepath.Join(dir, ".rustup"),
	}
}
