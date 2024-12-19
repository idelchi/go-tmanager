package tools

import (
	"fmt"
	"slices"
	"strings"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/pkg/unmarshal"
)

// Extensions represents a collection of file extensions.
type Extensions = unmarshal.SingleOrSliceType[string]

// ExtensionsToHint converts a list of extensions into a match.Hint.
func ExtensionsToHint(exts Extensions) match.Hint {
	var noExtensionPart string
	var extensionParts []string

	for _, ext := range slices.Compact(exts) {
		if ext == "" {
			noExtensionPart = "^[^.]+$"
		} else {
			escapedExt := strings.ReplaceAll(ext, ".", `\.`) // Escape dots in extensions
			extensionParts = append(extensionParts, fmt.Sprintf(".*%s$", escapedExt))
		}
	}

	var pattern string
	// Combine both parts
	if noExtensionPart != "" && len(extensionParts) > 0 {
		pattern = fmt.Sprintf("(%s|%s)", noExtensionPart, strings.Join(extensionParts, "|"))
	} else if noExtensionPart != "" {
		pattern = noExtensionPart
	} else {
		pattern = strings.Join(extensionParts, "|")
	}

	return match.Hint{
		Pattern: pattern,
		Must:    true,
	}
}
