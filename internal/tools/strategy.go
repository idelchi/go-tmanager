package tools

import (
	"fmt"
	"unicode"

	"github.com/Masterminds/semver/v3"
	"github.com/idelchi/godyl/internal/version"
)

// Strategy represents the strategy for handling tool installation or upgrades.
type Strategy string

const (
	// None indicates no strategy, meaning no action will be taken if the tool already exists.
	None Strategy = "none"
	// Upgrade indicates that the tool should only be upgraded if a newer version is available.
	Upgrade Strategy = "upgrade"
	// Force indicates that the tool should be installed or updated regardless of its current state.
	Force Strategy = "force"
)

// Check verifies if the tool needs any action based on the current strategy.
// It returns an error if the tool already exists and the strategy is set to None.
func (s Strategy) Check(t *Tool) error {
	// If the strategy is "None" and the tool already exists, return an error indicating it already exists.
	if t.Strategy == None && t.Exists() {
		return ErrAlreadyExists
	}

	return nil
}

// Upgrade checks if the tool should be upgraded based on the strategy and its current version.
// It compares the existing version with the desired version and returns an error if the tool is already up to date.
func (s Strategy) Upgrade(t *Tool) error {
	// If the tool does not exist, no upgrade is necessary.
	if !t.Exists() {
		return nil
	}

	switch t.Strategy {
	case None:
		// If the strategy is "None" and the tool exists, return an error indicating it already exists.
		return ErrAlreadyExists
	case Upgrade:
		if t.Version.Commands != nil && len(t.Version.Commands) == 0 {
			// No commands to run, so we can't check the version, forcing an upgrade.
			return nil
		}

		// Parse the version of the existing tool.
		exe := version.NewExecutable(t.Output, t.Exe.Name)

		parser := &version.Version{
			Patterns: t.Version.Patterns,
			Commands: t.Version.Commands,
		}

		if err := exe.ParseVersion(parser); err != nil {
			return nil
		}

		source, err := GetFirstSemVerish(exe.Version)
		if err != nil {
			return fmt.Errorf("parsing version %q: %w: %q -> %q", exe.Version, err, exe.Version, t.Version.Version)
		}

		target, err := GetFirstSemVerish(t.Version.Version)
		if err != nil {
			return fmt.Errorf("parsing version %q: %w: %q -> %q", t.Version.Version, err, exe.Version, t.Version.Version)
		}

		// If the versions match, return an error indicating the tool is already up to date.
		if source.Equal(target) {
			return fmt.Errorf("%w: current version %q and target version %q match", ErrUpToDate, exe.Version, t.Version.Version)
		}
		return nil
	case Force:
		// If the strategy is "Force", always proceed with the installation or update.
		return nil
	default:
		return nil
	}
}

// GetFirstSemVerish returns all leading non-digit characters until we hit the actual version part.
func GetFirstSemVerish(version string) (*semver.Version, error) {
	for index := range len(version) {
		candidate := version[index:]
		// First check if it starts with a digit
		if len(candidate) > 0 && !unicode.IsDigit(rune(candidate[0])) {
			continue
		}
		// Then check if it's valid semver
		if v, err := semver.NewVersion(candidate); err == nil {
			return v, nil
		}
	}

	return nil, fmt.Errorf("no valid semver found in %q", version)
}
