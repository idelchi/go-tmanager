package commands

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func PrintYAMLBytes(yamlBytes []byte) any {
	var data any
	if err := yaml.Unmarshal(yamlBytes, &data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	return data
}
