package command

import (
	"fmt"
	"strings"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/file"
)

// Commands represents a slice of shell commands.
type Commands []Command

// Get retrieves a specific attribute of the commands. (Ineffective).
func (c *Commands) Get(attribute string) string {
	return ""
}

// Initialize prepares the Commands based on the given string. (Ineffective).
func (*Commands) Initialize(_ string) error {
	return nil
}

// Exe executes the commands. (Ineffective).
func (*Commands) Exe() error {
	return nil
}

// Version sets the version for the commands. (Ineffective).
func (*Commands) Version(_ string) error {
	return nil
}

// Path sets up the path for the commands, using the provided parameters. (Ineffective).
func (*Commands) Path(_ string, _ []string, _ string, _ match.Requirements) error {
	return nil
}

// Combined returns all commands in the Commands slice as a single Command,
// with each command joined by semicolons.
func (c Commands) Combined() Command {
	stringCommands := make([]string, len(c))
	for i, cmd := range c {
		stringCommands[i] = string(cmd)
	}
	return Command(strings.Join(stringCommands, "; "))
}

// Install runs the combined commands for installation using the provided InstallData,
// captures the output, and returns it alongside any errors or found file information.
func (c Commands) Install(d common.InstallData) (output string, found file.File, err error) {
	cmd := c.Combined()

	// Execute the combined command
	output, err = cmd.Shell(d.Env.ToSlice()...)
	if err != nil {
		return output, "", fmt.Errorf("running combined commands: %w", err)
	}

	return strings.TrimRight(output, "\n"), "", nil
}
