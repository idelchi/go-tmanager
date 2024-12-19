package tools

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/utils"
)

func ErrCausesEarlyReturn(err error) bool {
	return errors.Is(ErrAlreadyExists, err) ||
		errors.Is(ErrUpToDate, err) ||
		errors.Is(ErrDoesHaveTags, err) ||
		errors.Is(ErrDoesNotHaveTags, err) ||
		errors.Is(ErrSkipped, err) ||
		errors.Is(ErrFailed, err)
}

var (
	// ErrAlreadyExists indicates that the tool already exists in the system.
	ErrAlreadyExists = fmt.Errorf("tool already exists")
	// ErrUpToDate indicates that the tool is already up to date.
	ErrUpToDate = fmt.Errorf("tool is up to date")
	// ErrDoesHaveTags indicates that the tool has tags that are in the excluded tags list.
	ErrDoesHaveTags = fmt.Errorf("tool contains excluded tags")
	// ErrDoesNotHaveTags indicates that the tool does not contain required tags.
	ErrDoesNotHaveTags = fmt.Errorf("tool does not contain included tags")
	// ErrSkipped indicates that the tool has been skipped due to conditions.
	ErrSkipped = fmt.Errorf("tool skipped")
	// ErrFailed indicates that the tool has failed to install or resolve.
	ErrFailed = fmt.Errorf("tool failed")
)

// Resolve attempts to resolve the tool's source and strategy based on the provided tags.
// It handles fallbacks and applies templating to the tool's fields as needed.
func (t *Tool) Resolve(withTags, withoutTags []string) error {
	if len(t.Name) == 0 {
		return fmt.Errorf("%w: tool name is empty", ErrFailed)
	}

	// Normalize values to ensure consistency in the .Values map.
	t.Values = utils.NormalizeMap(t.Values)

	// Load environment variables from the system.
	t.Env.Merge(env.FromEnv())

	// Expand and set the output folder path.
	output := file.Folder(t.Output)
	if err := output.Expand(); err != nil {
		return err
	}
	t.Output = output.Path()

	// Set the strategy to Force if the mode is "extract".
	if t.Mode == "extract" {
		t.Strategy = Force
	}

	// Save the path for templating later.
	path := t.Path

	t.Extensions = slices.Compact(t.Extensions)
	t.Aliases = slices.Compact(t.Aliases)

	// Expand environment variables.
	t.Env.Expand()

	if err := t.TemplateFirst(); err != nil {
		return err
	}

	t.Fallbacks = slices.Compact(t.Fallbacks)

	// Build the fallback sources from the primary source type and additional fallbacks.
	fallbacks := append([]sources.Type{t.Source.Type}, t.Fallbacks...)

	var lastErr error
	// Try resolving with each fallback in order.
	for _, fallback := range slices.Compact(fallbacks) {
		if fallback == sources.RUST {
			// Skip Rust fallback for now.
			continue
		}

		if err := t.tryResolveFallback(fallback, path, withTags, withoutTags); ErrCausesEarlyReturn(err) {
			return err
		} else if err != nil {
			lastErr = err

			continue // Move on to the next fallback.
		}

		return nil // Success, no need to try further fallbacks.
	}

	// If all fallbacks fail, return the last encountered error.
	return lastErr
}

// CheckSkipConditions verifies whether the tool should be skipped based on its tags or strategy.
func (t *Tool) CheckSkipConditions(withTags, withoutTags []string) error {
	if !t.Tags.Has(withTags) {
		return fmt.Errorf("%w: %v: tool tags: %v", ErrDoesNotHaveTags, withTags, t.Tags)
	}

	if !t.Tags.HasNot(withoutTags) {
		return fmt.Errorf("%w: %v: tool tags: %v", ErrDoesHaveTags, withoutTags, t.Tags)
	}

	if err := t.Strategy.Check(t); err != nil {
		return err
	}

	if skip, msg, _ := t.Skip.True(); skip {
		return fmt.Errorf("%w: %q", ErrSkipped, msg)
	}

	return nil
}

// tryResolveFallback attempts to resolve a tool using a specific fallback source type.
func (t *Tool) tryResolveFallback(fallback sources.Type, path string, withTags, withoutTags []string) error {
	// Append the tool's name as a tag.
	t.Tags.Append(t.Name)

	// Check if the tool should be skipped based on its conditions.
	if err := t.CheckSkipConditions(withTags, withoutTags); err != nil {
		return err
	}

	// Set the source type to the current fallback.
	t.Source.Type = fallback

	// Get the installer for the current source type.
	populator, err := t.Source.Installer()
	if err != nil {
		return err
	}

	// Initialize the installer.
	if err := populator.Initialize(t.Name); err != nil {
		return err
	}

	// Retrieve executable details from the installer.
	if err := populator.Exe(); err != nil {
		return err
	}

	// Apply templating to the tool's fields.
	utils.SetIfEmpty(&t.Exe.Name, populator.Get("exe"))
	utils.SetIfEmpty(&t.Exe.Name, t.Name)

	// Re-check skip conditions after applying templates.
	if err := t.CheckSkipConditions(withTags, withoutTags); err != nil {
		return err
	}

	// Retrieve the tool's version from the installer if it is not already set.
	if utils.IsEmpty(t.Version.Version) {
		if err := populator.Version(t.Name); err != nil {
			return err
		}
	}

	utils.SetIfEmpty(&t.Version.Version, populator.Get("version"))

	if err := t.TemplateLast(); err != nil {
		return err
	}

	// Determine the tool's path if not already set.
	if utils.IsEmpty(t.Path) {
		hints := t.Hints
		hints.Add(ExtensionsToHint(t.Extensions))

		if err := populator.Path(t.Name, nil, t.Version.Version, match.Requirements{
			Platform: t.Platform,
			Hints:    hints,
		}); err != nil {
			return err
		}
	}
	utils.SetIfEmpty(&t.Path, populator.Get("path"))
	utils.SetIfEmpty(&t.Path, path)

	// Append platform-specific file extension to the executable name.
	if !strings.HasSuffix(t.Exe.Name, t.Platform.Extension.String()) {
		t.Exe.Name += t.Platform.Extension.String()
	}

	// Set patterns for finding the executable.
	utils.SetSliceIfNil(&t.Exe.Patterns, fmt.Sprintf("^%s$", t.Exe.Name))

	// Append platform-specific extensions to aliases.
	for i, alias := range t.Aliases {
		t.Aliases[i] = alias + t.Platform.Extension.String()
	}

	// Attempt to upgrade the tool using the current strategy.
	if err := t.Strategy.Upgrade(t); err != nil {
		return err
	}

	// Validate the tool's configuration.
	return t.Validate()
}

// Validate validates the Tool's configuration using the validator package.
func (t *Tool) Validate() error {
	validate := validator.New()
	if err := validate.Struct(t); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}
	return nil
}

// Exists checks if the tool's executable already exists in the output path.
func (t *Tool) Exists() bool {
	f := file.NewFile(t.Output, t.Exe.Name)
	return f.Exists() && f.IsFile()
}

// Download downloads the tool using its configured source and installer.
func (t *Tool) Download() (string, file.File, error) {
	installer, err := t.Source.Installer()
	if err != nil {
		return "", "", err
	}

	data := common.InstallData{
		Path:        t.Path,
		Name:        t.Name,
		Exe:         t.Exe.Name,
		Patterns:    t.Exe.Patterns,
		Output:      t.Output,
		Aliases:     t.Aliases,
		Mode:        t.Mode.String(),
		Env:         t.Env,
		NoVerifySSL: t.NoVerifySSL,
	}

	return installer.Install(data)
}
