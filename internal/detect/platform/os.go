package platform

import (
	"fmt"
	"strings"
)

// OS represents an operating system.
type OS struct {
	Type string
	Raw  string // Original parsed OS value
}

// OSInfo holds information about an OS type, including aliases.
type OSInfo struct {
	Type    string
	Aliases []string
}

// Supported returns a slice of supported operating system information.
func (OSInfo) Supported() []OSInfo {
	return []OSInfo{
		{
			Type: "linux",
		},
		{
			Type:    "darwin",
			Aliases: []string{"macos", "mac", "osx"},
		},
		{
			Type:    "windows",
			Aliases: []string{"win"},
		},
		{
			Type: "freebsd",
		},
		{
			Type: "android",
		},
		{
			Type: "netbsd",
		},
		{
			Type: "openbsd",
		},
	}
}

// Parse attempts to parse the OS from the given name string.
func (o *OS) Parse(name string) error {
	name = strings.ToLower(name)

	info := OSInfo{}

	for _, info := range info.Supported() {
		for _, alias := range append([]string{info.Type}, info.Aliases...) {
			if strings.Contains(name, alias) {
				o.Type = info.Type
				o.Raw = alias
				return nil
			}
		}
	}

	return fmt.Errorf("unable to parse OS from name: %s", name)
}

// IsUnset returns true if the OS type is not set.
func (o OS) IsUnset() bool {
	return o.Type == ""
}

// Is checks if this OS is exactly the same as another.
func (o OS) Is(other OS) bool {
	return other.Raw == o.Raw && !o.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if this OS is compatible with another.
func (o OS) IsCompatibleWith(other OS) bool {
	if o.IsUnset() || other.IsUnset() {
		return false
	}

	return o.Type == other.Type
}

// String returns a string representation of the OS.
func (o OS) String() string {
	return o.Type
}
