package pretty

import (
	"bytes"
	"encoding/json"

	"github.com/showa-93/go-mask"

	"gopkg.in/yaml.v3"
)

// YAML returns a prettified YAML representation of the provided object.
func YAML(obj any) string {
	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)

	if err := enc.Encode(&obj); err != nil {
		return err.Error()
	}

	return buf.String()
}

// YAMLMasked returns a prettified YAML representation of the provided object
// with masked sensitive fields. It uses JSONMasked internally to mask the fields.
func YAMLMasked(obj any) string {
	return YAML(MaskJSON(obj))
}

// JSON returns a prettified JSON representation of the provided object.
func JSON(obj any) string {
	bytes, err := json.MarshalIndent(obj, "  ", "    ")
	if err != nil {
		return err.Error()
	}

	return string(bytes)
}

// JSONMasked returns a prettified JSON representation of the provided object
// with masked sensitive fields.
func JSONMasked(obj any) string {
	return JSON(MaskJSON(obj))
}

// MaskJSON masks sensitive fields in the provided object according to predefined masking rules
// and returns the masked object. The masker replaces sensitive data with a mask character.
func MaskJSON(obj any) any {
	masker := mask.NewMasker()

	// Set the masking character to "-"
	masker.SetMaskChar("-")

	// Register a function to mask strings by filling them with the masking character.
	masker.RegisterMaskStringFunc(mask.MaskTypeFilled, masker.MaskFilledString)

	t, err := mask.Mask(obj)
	if err != nil {
		return err.Error()
	}

	return t
}
