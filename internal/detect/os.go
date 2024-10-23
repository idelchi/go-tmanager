package detect

import (
	"fmt"

	"github.com/idelchi/godyl/internal/detect/platform"
	"github.com/shirou/gopsutil/host"
)

// Detect gathers information about the current platform, such as the operating system, architecture,
// distribution, library, and file extension, and populates the Platform struct accordingly.
func (p *Platform) Detect() error {
	var os platform.OS
	var arch platform.Architecture
	var library platform.Library
	var distro platform.Distribution
	var extension platform.Extension

	info, err := host.Info()
	if err != nil {
		return fmt.Errorf("getting host information: %w", err)
	}

	// Determine the OS from runtime information
	if err := os.Parse(info.OS); err != nil {
		return err
	}

	// Determine the Linux distribution from system information
	distro.Parse(info.Platform)

	// Set the default library based on the OS and distribution
	library = library.Default(os, distro)

	// Determine the architecture from the system's kernel architecture
	if err := arch.Parse(info.KernelArch); err != nil {
		return err
	}

	if arch.Is64Bit() && os.Type == "linux" {
		is32Bit, err := platform.Is32Bit()
		if err == nil && is32Bit {
			arch.To32BitUserLand()
		}
	}

	// Populate the Platform struct with the detected values
	*p = Platform{
		OS:           os,
		Distribution: distro,
		Architecture: arch,
		Library:      library,
		Extension:    extension.Default(os),
	}

	return nil
}
