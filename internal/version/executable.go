package version

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"

	"github.com/idelchi/godyl/pkg/file"
)

// Executable consists of a full path to a file and its version.
// An attempt to parse the version into a string can be done by using a `Version` type.
type Executable struct {
	// File represents the full path to the executable file.
	File file.File
	// Version holds the parsed version of the executable.
	Version string
}

// NewExecutable creates a new Executable instance from the provided paths.
func NewExecutable(paths ...string) Executable {
	return Executable{File: file.NewFile(paths...)}
}

// Command runs the specified command arguments on the executable using the provided context.
// It returns the output of the command as a trimmed string and any error encountered during execution.
func (e Executable) Command(ctx context.Context, cmdArgs []string) (string, error) {
	var out bytes.Buffer

	cmd := exec.CommandContext(ctx, e.File.Name(), cmdArgs...)
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	return strings.TrimSpace(out.String()), err
}

// ParseVersion attempts to parse the version of the executable using the provided Version object.
// It iterates over predefined command strategies and tries to parse the version from the command output.
// If successful, it sets the Version field of Executable; otherwise, it returns an error.
func (e *Executable) ParseVersion() error {
	timeout := 30 * time.Second

	// Create a context with a timeout to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	version := NewDefaultVersionParser()

	// Iterate through each command strategy
	for _, cmdArgs := range version.Commands {
		// Get the output of the command
		output, err := e.Command(ctx, cmdArgs)
		if err != nil {
			continue
		}

		if version, err := version.ParseString(output); err == nil {
			e.Version = version

			return nil
		} else {
			continue
		}
	}

	e.Version = ""

	return errors.New("unable to parse version from output")
}
