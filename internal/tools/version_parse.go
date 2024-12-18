package tools

import (
	"github.com/fatih/structs"
	"github.com/idelchi/godyl/pkg/unmarshal"
	"gopkg.in/yaml.v3"
)

type Version struct {
	Version  string
	Commands unmarshal.SingleOrSlice[string]
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
