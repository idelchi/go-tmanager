package tools

import (
	"github.com/fatih/structs"
	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

// Version represents the version configuration for a tool.
type Version struct {
	// Version holds the string representation of the parsed version.
	Version string
	// Commands contains the list of command strategies used to extract the version.
	Commands unmarshal.SingleOrSlice[string]
	// Patterns contains the list of regex patterns for parsing the version from output strings.
	Patterns unmarshal.SingleOrSlice[string]
}

// UnmarshalYAML implements custom unmarshaling for Exe,
// allowing the YAML to either provide just the name as a scalar or the full Exe structure.
func (v *Version) UnmarshalYAML(value *yaml.Node) error {
	// If the YAML value is a scalar (e.g., just the version), handle it directly by setting the Version field.
	if value.Kind == yaml.ScalarNode {
		v.Version = value.Value

		return nil
	}

	// Perform custom unmarshaling with field validation, allowing only known fields.
	type raw Version

	return unmarshal.DecodeWithOptionalKnownFields(value, (*raw)(v), true, structs.New(v).Name())
}
