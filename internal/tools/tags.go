package tools

import (
	"slices"

	"github.com/idelchi/godyl/pkg/unmarshal"

	"gopkg.in/yaml.v3"
)

// Tags represents a list of tags associated with a tool.
// Tags can be used to categorize or filter tools based on specific labels or keywords.
type Tags []string

// UnmarshalYAML implements custom unmarshaling for Tags,
// allowing the field to be either a single string or a list of strings.
func (t *Tags) UnmarshalYAML(value *yaml.Node) error {
	result, err := unmarshal.UnmarshalSingleOrSlice[string](value, false)
	if err != nil {
		return err
	}
	*t = result
	return nil
}

// Append adds new tags to the Tags list, avoiding duplicates.
func (t *Tags) Append(tags ...string) {
	for _, tag := range tags {
		if !slices.Contains(*t, tag) {
			*t = append(*t, tag)
		}
	}
}

// Has checks if any of the provided tags are present in the Tags list.
// Returns true if at least one tag matches.
func (t Tags) Has(tags Tags) bool {
	if len(tags) == 0 {
		return true
	}

	for _, tag := range tags {
		if slices.Contains(t, tag) {
			return true
		}
	}

	return false
}

// HasNot checks if none of the provided tags are present in the Tags list.
// Returns true if none of the tags match.
func (t Tags) HasNot(tags Tags) bool {
	if len(tags) == 0 {
		return true
	}
	return !t.Has(tags)
}
