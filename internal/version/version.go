package version

import (
	"errors"
	"regexp"
	"strings"
)

// Version provides functionality for parsing version strings using regex patterns
// and defines command strategies for extracting the version from an executable.
type Version struct {
	// Patterns contains the list of regex patterns for parsing the version from output strings.
	Patterns []*regexp.Regexp
	// Commands contains the list of command strategies used to extract the version.
	Commands []string
	// String holds the string representation of the parsed version.
	String string
}

// NewDefaultVersionParser creates a new Version parser with default regex patterns
// for common version formats and predefined command strategies for extracting the version.
func NewDefaultVersionParser() *Version {
	return &Version{
		Patterns: []*regexp.Regexp{
			// Pattern for versions formatted as X.X.X, surrounded by any characters.
			regexp.MustCompile(`.*?(\d+\.\d+\.\d+).*`),
			// Pattern for versions formatted as X.X, surrounded by any characters.
			regexp.MustCompile(`.*?(\d+\.\d+).*`),
		},
		Commands: []string{
			"--version", // Attempt to get version with --version flag.
			"-v",        // Attempt to get version with -v flag.
			"-version",  // Attempt to get version with -version flag.
			"version",   // Attempt to get version with version command.
		},
	}
}

// ParseString attempts to parse the provided output string using the defined regex patterns.
// It normalizes multi-line output into a single line and tries to match the patterns.
// Returns the first matched version string or an error if no match is found.
func (v *Version) ParseString(output string) (string, error) {
	// Normalize the output into a single line by replacing newlines with spaces.
	normalizedOutput := strings.Join(strings.Split(output, "\n"), " ")

	// Try to match each regex pattern on the normalized output.
	for _, pattern := range v.Patterns {
		if matches := pattern.FindStringSubmatch(normalizedOutput); len(matches) > 1 {
			// Return the first matched version group.
			return matches[1], nil
		}
	}

	return "", errors.New("unable to parse version from output")
}
