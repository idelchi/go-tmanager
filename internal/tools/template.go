package tools

import (
	"github.com/idelchi/godyl/internal/templates"
	"github.com/idelchi/godyl/pkg/utils"
)

func (t *Tool) ToTemplateMap(flatten ...map[string]any) map[string]any {
	templateMap := map[string]any{
		"Name":    t.Name,
		"Env":     t.Env,
		"Values":  t.Values,
		"Version": t.Version,
		"Exe":     t.Exe.Name,
		"Output":  t.Output,
	}

	for _, o := range flatten {
		for k, v := range o {
			templateMap[k] = v
		}
	}

	return templateMap
}

// Template applies templating to various fields of the Tool struct, such as version, path, and checksum.
// It processes these fields using Go templates and updates them with the templated values.
func (t *Tool) TemplateFirst() error {
	values := t.ToTemplateMap(t.Platform.ToMap())

	// Apply templating to Source.Type
	output, err := templates.Apply(t.Source.Type.String(), values)
	if err != nil {
		return err
	}
	t.Source.Type.From(output)

	// Apply templating to the Skip conditions
	for i := range t.Skip {
		err = templates.ApplyAndSet(&t.Skip[i].Condition, values)
		if err != nil {
			return err
		}
	}

	// Validate the Skip conditions
	if _, _, err := t.Skip.True(); err != nil {
		return err
	}

	if err := templates.ApplyAndSet(&t.Output, values); err != nil {
		return err
	}

	if err := templates.ApplyAndSet(&t.Version, values); err != nil {
		return err
	}

	return nil
}

func (t *Tool) TemplateLast() error {
	values := t.ToTemplateMap(t.Platform.ToMap())

	if err := templates.ApplyAndSet(&t.Exe.Name, values); err != nil {
		return err
	}

	// Apply templating to Exe.Patterns
	for i := range t.Exe.Patterns {
		if err := templates.ApplyAndSet(&t.Exe.Patterns[i], values); err != nil {
			return err
		}
	}

	// Apply templating to Extensions
	for i := range t.Extensions {
		if err := templates.ApplyAndSet(&t.Extensions[i], values); err != nil {
			return err
		}
	}

	// Apply templating to Source.Commands
	for i, cmd := range t.Source.Commands {
		output, err := templates.Apply(cmd.String(), values)
		if err != nil {
			return err
		}
		t.Source.Commands[i].From(output)
	}

	// Apply templating to Post commands
	for i, cmd := range t.Post {
		output, err := templates.Apply(cmd.String(), values)
		if err != nil {
			return err
		}
		t.Post[i].From(output)
	}

	// Apply templating to Hints patterns and weights
	for i := range t.Hints {
		if err := templates.ApplyAndSet(&t.Hints[i].Pattern, values); err != nil {
			return err
		}
		if err := templates.ApplyAndSet(&t.Hints[i].Weight, values); err != nil {
			return err
		}
		// Set a default weight of "1" if not specified and convert it to an integer
		utils.SetIfEmpty(&t.Hints[i].Weight, "1")

		if err := t.Hints[i].SetWeight(); err != nil {
			return err
		}
	}

	if err := templates.ApplyAndSet(&t.Path, values); err != nil {
		return err
	}

	return nil
}
