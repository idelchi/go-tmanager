package tools

import (
	"strconv"

	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Skip represents a list of conditions under which certain operations (e.g., downloading, testing) should be skipped.
type Skip []Condition

// Condition defines a condition and an optional reason for skipping an operation.
type Condition struct {
	// Condition is a string that represents a boolean expression (e.g., "true" or "false") that determines whether the
	// operation should be skipped.
	Condition string
	// Reason provides an optional explanation for why the operation is being skipped.
	Reason string
}

// True checks if any condition in the Skip list evaluates to true.
// It returns a boolean indicating if the skip should occur, the associated reason, and any error encountered while
// evaluating the condition.
func (s Skip) True() (bool, string, error) {
	for _, condition := range s {
		// Parse the condition string into a boolean value.
		if val, err := strconv.ParseBool(condition.Condition); err != nil {
			return false, condition.Reason, err
		} else {
			if val {
				return true, condition.Reason, nil
			}
		}
	}
	return false, "", nil
}

// UnmarshalYAML implements custom unmarshaling for Skip,
// allowing the Skip field to be either a single condition or a list of conditions.
func (s *Skip) UnmarshalYAML(value *yaml.Node) error {
	// If the YAML value is a scalar (e.g., just a single condition), handle it directly.
	if value.Kind == yaml.ScalarNode {
		*s = []Condition{{Condition: value.Value}}
		return nil
	}

	// Otherwise, treat it as a list of conditions and unmarshal accordingly.
	result, err := unmarshal.UnmarshalSingleOrSlice[Condition](value, true)
	if err != nil {
		return err
	}
	*s = result
	return nil
}
