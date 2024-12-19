// Package unmarshal provides utilities for unmarshalling YAML data that can
// represent either a single item or a slice of items. It includes a generic
// type `SingleOrSlice` to handle this pattern, allowing flexible unmarshalling
// from YAML input.
//
// The package also provides functions to decode YAML nodes while optionally
// enforcing that only known fields are present, improving error handling
// in case of unexpected fields in the YAML input.
package unmarshal

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// SingleOrSliceType represents a custom type that can unmarshal from YAML as either
// a single element or a slice of elements. It is a generic type that works for
// any type T.
type SingleOrSliceType[T any] []T

// UnmarshalYAML implements the yaml.Unmarshaler interface for SingleOrSlice.
// It allows the YAML value to be unmarshaled either as a single element or a slice.
// The unmarshaled result is assigned to the receiver.
func (ss *SingleOrSliceType[T]) UnmarshalYAML(value *yaml.Node) error {
	result, err := SingleOrSlice[T](value, true)
	if err != nil {
		return err
	}

	*ss = result

	return nil
}

// SingleOrSlice is a helper function that unmarshals a YAML node into
// a slice of type T. It handles cases where the node could represent a single
// item or a list. If the node is a scalar or a mapping, it wraps it in a sequence node.
// The useKnownFields parameter controls whether unknown fields trigger an error.
func SingleOrSlice[T any](node *yaml.Node, useKnownFields bool) ([]T, error) {
	// If it's a scalar or mapping node, wrap it in a sequence node
	if node.Kind == yaml.ScalarNode || node.Kind == yaml.MappingNode {
		node = &yaml.Node{
			Kind:    yaml.SequenceNode,
			Content: []*yaml.Node{node},
		}
	}

	var result []T

	// Use DecodeWithOptionalKnownFields to decode the node into the result slice.
	if err := DecodeWithOptionalKnownFields(node, &result, useKnownFields, "any"); err != nil {
		return nil, err
	}

	return result, nil
}

// DecodeWithOptionalKnownFields is a helper function that decodes a YAML node into
// the provided output (out) interface. It optionally enforces that all fields in the YAML
// node are known to the target type if useKnownFields is set to true.
// The input parameter is used for error message formatting.
func DecodeWithOptionalKnownFields(value *yaml.Node, out any, useKnownFields bool, input string) error {
	// Re-encode the yaml.Node to bytes
	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	if err := enc.Encode(value); err != nil {
		return err
	}

	if err := enc.Close(); err != nil {
		return err
	}

	// Decode from the buffer
	decoder := yaml.NewDecoder(&buf)
	if useKnownFields {
		decoder.KnownFields(true)
	}

	// Decode into the provided interface
	err := decoder.Decode(out)

	return yamlTypeErrorConversion(err, input)
}

// yamlTypeErrorConversion converts yaml.TypeError errors into more informative messages
// by including the actual type of the input. It modifies the error message when it detects
// a "not found in type" message, appending the input type to the message.
func yamlTypeErrorConversion(err error, input string) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*yaml.TypeError); !ok {
		return err
	}

	if !strings.Contains(err.Error(), " not found in type") {
		return err
	}

	typeErr := err.(*yaml.TypeError)

	for i, err := range typeErr.Errors {
		parts := strings.SplitN(err, " not found in type", 2)
		if len(parts) != 2 {
			continue
		}

		typeErr.Errors[i] = fmt.Sprintf("%s not found in type %q", parts[0], input)
	}

	return err
}
