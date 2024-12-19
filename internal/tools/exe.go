package tools

import (
	"github.com/fatih/structs"

	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Exe represents the configuration for the executable, including
// the name under which the binary will be stored and the patterns
// to search for the binary within the downloaded files.
type Exe struct {
	// Name is the name under which the binary will be stored in the output folder.
	Name string
	// Patterns specifies the patterns used to locate the binary in the downloaded folder.
	// This can either be a single string or a slice of strings.
	Patterns unmarshal.SingleOrSliceType[string]
}

// UnmarshalYAML implements custom unmarshaling for Exe,
// allowing the YAML to either provide just the name as a scalar or the full Exe structure.
func (e *Exe) UnmarshalYAML(value *yaml.Node) error {
	// If the YAML value is a scalar (e.g., just the name), handle it directly by setting the Name field.
	if value.Kind == yaml.ScalarNode {
		e.Name = value.Value

		return nil
	}

	// Perform custom unmarshaling with field validation, allowing only known fields.
	type raw Exe

	return unmarshal.DecodeWithOptionalKnownFields(value, (*raw)(e), true, structs.New(e).Name())
}
