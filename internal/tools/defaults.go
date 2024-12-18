package tools

import (
	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/env"
)

// Defaults holds default configuration values for a tool.
// These values can be applied to a tool configuration to provide sensible defaults
// for executable paths, output locations, platform-specific settings, and more.
type Defaults struct {
	// Exe specifies default executable details such as patterns for identifying the binary.
	Exe Exe
	// Output specifies the default output path for the tool.
	Output string
	// Platform defines default platform-specific settings (e.g., OS and architecture).
	Platform detect.Platform
	// Values contains default custom values for the tool configuration.
	Values map[string]any
	// Fallbacks defines default fallback configurations or sources in case the primary configuration fails.
	Fallbacks []string
	// Hints provide default matching patterns or heuristics for the tool.
	Hints match.Hints
	// Source defines the default source configuration for fetching the tool (e.g., GitHub, local files).
	Source sources.Source
	// Tags are default labels or markers for categorizing the tool.
	Tags Tags
	// Strategy defines the default deployment or fetching strategy for the tool.
	Strategy Strategy
	// Extensions lists default additional file extensions related to the tool.
	Extensions []string
	// Env defines default environment variables applied when running the tool.
	Env env.Env
	// Mode specifies the default operating mode for the tool (e.g., silent mode, verbose mode).
	Mode Mode
	// Version specifies the default version details for the tool.
	Version Version
}

// Initialize detects the current platform and applies platform-specific defaults to the Defaults struct.
// It also sets up default extensions based on the detected platform.
func (d *Defaults) Initialize() error {
	// Detect the current platform (e.g., OS, architecture).
	platform := detect.Platform{}
	if err := platform.Detect(); err != nil {
		return err
	}

	// Merge the detected platform details with the default platform settings.
	d.Platform.Merge(platform)

	return nil
}
