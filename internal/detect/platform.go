package detect

import (
	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/idelchi/godyl/pkg/utils"
)

// Platform defines the characteristics of the platform, including OS, architecture, library, extension, and
// distribution.
type Platform struct {
	OS           platform.OS           // OS represents the operating system (e.g., Linux, Windows, macOS).
	Architecture platform.Architecture // Architecture defines the platform's CPU architecture and version.
	Library      platform.Library      // Library specifies the system library (e.g., GNU, Musl, MSVC).
	Extension    platform.Extension    // Extension represents the default file extension for executables.
	Distribution platform.Distribution // Distribution refers to the Linux distribution, if applicable.
}

// Parse attempts to parse the platform's OS, architecture, and library from a given string.
func (p *Platform) Parse(name string) {
	p.OS.Parse(name)
	p.Architecture.Parse(name)
	p.Library.Parse(name)
	p.Extension.Default(p.OS)
}

// Merge combines another Platform's fields into the current Platform, setting fields that are empty.
func (p *Platform) Merge(other Platform) {
	utils.SetIfZeroValue(&p.OS, other.OS)
	utils.SetIfZeroValue(&p.Architecture.Type, other.Architecture.Type)
	utils.SetIfZeroValue(&p.Architecture.Version, other.Architecture.Version)
	utils.SetIfZeroValue(&p.Architecture.Raw, other.Architecture.Raw)
	utils.SetIfZeroValue(&p.Library, other.Library)
	utils.SetIfZeroValue(&p.Extension, other.Extension)
	utils.SetIfZeroValue(&p.Distribution, other.Distribution)
}

// ToMap converts the Platform struct to a map for use in templates.
func (p Platform) ToMap() map[string]any {
	platformMap := make(map[string]any)
	platformMap["OS"] = p.OS.String()
	platformMap["ARCH"] = p.Architecture.Type
	platformMap["ARCH_VERSION"] = p.Architecture.Version
	platformMap["LIBRARY"] = p.Library.String()
	platformMap["EXTENSION"] = p.Extension.String()
	platformMap["DISTRIBUTION"] = p.Distribution.String()

	return platformMap
}
