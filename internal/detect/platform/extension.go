package platform

import (
	"path/filepath"
	"strings"
)

// Extension represents a file extension, typically used for executable files on different operating systems.
type Extension string

// Default returns the default file extension based on the operating system.
// For Windows, it returns ".exe", and for other operating systems, it returns an empty string.
func (e Extension) Default(os OS) Extension {
	switch os.Type {
	case "windows":
		return Extension(".exe")
	default:
		return Extension("")
	}
}

// String returns the Extension as a string.
func (e Extension) String() string {
	return string(e)
}

func (e *Extension) Parse(name string) error {
	switch ext := filepath.Ext(name); ext {
	case ".gz":
		if strings.HasSuffix(name, ".tar.gz") {
			*e = Extension(".tar.gz")
		} else {
			*e = Extension(ext)
		}
	default:
		*e = Extension(ext)
	}

	return nil
}
