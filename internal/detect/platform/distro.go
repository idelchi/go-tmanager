package platform

import (
	"fmt"
	"strings"
)

// Distribution represents a Linux distribution.
type Distribution struct {
	Type string
	Raw  string // Original parsed distribution value
}

// DistroInfo holds information about a distribution type, including aliases.
type DistroInfo struct {
	Type    string
	Aliases []string
}

// Supported returns a slice of supported distribution information.
func (DistroInfo) Supported() []DistroInfo {
	return []DistroInfo{
		{
			Type: "debian",
			// Aliases: []string{"debian"},
		},
		{
			Type: "ubuntu",
			// Aliases: []string{"ubuntu"},
		},
		{
			Type: "centos",
			// Aliases: []string{"centos"},
		},
		{
			Type:    "redhat",
			Aliases: []string{"rhel"},
		},
		{
			Type: "arch",
			// Aliases: []string{"arch"},
		},
		{
			Type: "alpine",
			// Aliases: []string{"alpine"},
		},
		{
			Type:    "raspbian",
			Aliases: []string{"raspberry"},
		},
	}
}

// Parse attempts to parse the distribution from the given name string.
func (d *Distribution) Parse(name string) error {
	name = strings.ToLower(name)

	info := DistroInfo{}

	for _, info := range info.Supported() {
		// for _, alias := range info.Aliases {
		for _, alias := range append([]string{info.Type}, info.Aliases...) {
			if strings.Contains(name, alias) {
				d.Type = info.Type
				d.Raw = alias
				return nil
			}
		}
	}

	return fmt.Errorf("unable to parse distribution from name: %s", name)
}

// IsUnset returns true if the distribution type is not set.
func (d Distribution) IsUnset() bool {
	return d.Type == ""
}

// Is checks if this distribution is exactly the same as another.
func (d Distribution) Is(other Distribution) bool {
	return other.Raw == d.Raw && !d.IsUnset() && !other.IsUnset()
}

// IsCompatibleWith checks if this distribution is compatible with another.
func (d Distribution) IsCompatibleWith(other Distribution) bool {
	if d.IsUnset() || other.IsUnset() {
		return false
	}

	return d.Type == other.Type
}

// String returns a string representation of the distribution.
func (d Distribution) String() string {
	return d.Type
}

// Default returns the default Distribution, which is an empty Distribution.
func Default() Distribution {
	return Distribution{}
}
