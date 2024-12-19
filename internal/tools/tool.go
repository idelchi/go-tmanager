package tools

import (
	"github.com/fatih/structs"

	"github.com/idelchi/godyl/internal/detect"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/command"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/unmarshal"
	"github.com/idelchi/godyl/pkg/utils"

	"gopkg.in/yaml.v3"
)

// Tool represents a single tool configuration.
// It contains various fields that specify details such as the tool's name, version, path, execution settings,
// platform-specific settings, environment variables, and custom strategies for downloading, testing, or deploying.
type Tool struct {
	// Name of the tool, usually a short identifier or title.
	Name string
	// Description of the tool, giving more context about its purpose.
	Description string
	// Version specifies the version of the tool.
	Version Version
	// Path represents the URL or file path where the tool can be fetched or downloaded from.
	Path string
	// Output defines the output path where the tool will be installed or extracted.
	Output string
	// Exe specifies the executable details for the tool, such as patterns or names for locating the binary.
	Exe Exe
	// Platform defines the platform-specific details for the tool, including OS and architecture constraints.
	Platform detect.Platform
	// Aliases represent alternative names or shortcuts for the tool.
	Aliases Aliases
	// Values contains custom values or variables used in the tool's configuration.
	Values map[string]any
	// Fallbacks defines fallback configurations in case the primary configuration fails.
	Fallbacks Fallbacks
	// Hints provide additional matching patterns or heuristics for the tool.
	Hints match.Hints
	// Source defines the source configuration, which determines how the tool is fetched (e.g., GitHub, local files).
	Source sources.Source
	// Tags are labels or markers that can be used to categorize or filter the tool.
	Tags Tags
	// Strategy defines how the tool is deployed, fetched, or managed (e.g., download strategies, handling retries).
	Strategy Strategy
	// Extensions lists additional files or behaviors that are tied to the tool.
	Extensions Extensions
	// Skip defines conditions under which certain steps (e.g., downloading, testing) are skipped.
	Skip Skip
	// Post defines commands that should be run after the main operation, such as post-installation steps.
	Post command.Commands
	// Mode defines the operating mode for the tool, potentially controlling behavior such as silent mode or verbose
	// mode.
	Mode Mode
	// Settings contains custom settings or options that modify the behavior of the tool.
	Settings Settings
	// Env defines the environment variables that are applied when running the tool.
	Env env.Env
	// Check defines a set of instructions for verifying the tool's integrity or functionality.
	Check Checker
	// NoVerifySSL specifies whether SSL verification should be disabled when fetching the tool.
	NoVerifySSL bool `json:"-" mapstructure:"-" yaml:"-"`
}

// UnmarshalYAML implements custom unmarshaling for Tool with KnownFields check.
// This allows the Tool to be unmarshaled from YAML while verifying that only known fields are present,
// ensuring stricter validation and preventing unexpected fields.
func (t *Tool) UnmarshalYAML(value *yaml.Node) error {
	// If it's a scalar (e.g., just the name), handle it directly by assigning it to the Name field.
	if value.Kind == yaml.ScalarNode {
		t.Name = value.Value

		if utils.IsURL(t.Name) {
			t.Path = t.Name
			t.Source.Type = sources.DIRECT
		}

		return nil
	}

	type rawTool Tool

	// Use custom unmarshal logic with KnownFields check to ensure field validation.

	return unmarshal.DecodeWithOptionalKnownFields(value, (*rawTool)(t), true, structs.New(t).Name())
}

// ApplyDefaults applies default values to the Tool configuration.
// If a field is empty or nil, it is replaced with the corresponding default from the Defaults struct.
// TODO(Idelchi): Improve - what if someone wants a value to be ""?
func (t *Tool) ApplyDefaults(d Defaults) {
	utils.SetIfZeroValue(&t.Output, d.Output)
	utils.SetIfZeroValue(&t.Source.Type, d.Source.Type)
	utils.SetIfZeroValue(&t.Source.Github.Token, d.Source.Github.Token)
	utils.SetIfZeroValue(&t.Strategy, d.Strategy)
	utils.SetSliceIfNil(&t.Skip, Condition{Condition: "false"})
	utils.SetIfZeroValue(&t.Mode, d.Mode)
	utils.SetSliceIfNil(&t.Exe.Patterns, d.Exe.Patterns...)
	utils.SetSliceIfNil(&t.Extensions, d.Extensions...)
	utils.SetSliceIfNil(&t.Version.Commands, d.Version.Commands...)
	utils.SetSliceIfNil(&t.Version.Patterns, d.Version.Patterns...)
	utils.SetMapIfNil(&t.Values, d.Values)
	utils.DeepMergeMapsWithoutOverwrite(t.Values, d.Values)
	t.Env.Merge(d.Env)

	// Apply platform-specific defaults and hints.
	t.Platform.Merge(d.Platform)
	t.Hints.Append(d.Hints)
}
