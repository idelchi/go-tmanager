// Package url provides functionality to handle URLs as sources for downloading and managing files.
// It supports initializing URLs, retrieving metadata, and performing installations based on
// download operations. This package integrates with external utilities for file handling and
// matching patterns during installation processes.
package url

import (
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/file"
)

// URL represents a download source with an associated URL and optional token for authorization.
type URL struct {
	URL   string
	Token string

	// Data holds additional metadata related to the URL.
	Data common.Metadata `yaml:"-"`
}

// Get retrieves a specific attribute from the URL's metadata.
func (u *URL) Get(attribute string) string {
	return u.Data.Get(attribute)
}

// Initialize prepares the URL based on the given name. (Ineffective).
func (u *URL) Initialize(name string) error {
	return nil
}

// Exe executes the URL's associated action. (Ineffective).
func (u *URL) Exe() error {
	return nil
}

// Version sets the version for the URL. (Ineffective).
func (u *URL) Version(name string) error {
	return nil
}

// Path sets the path for the URL, storing the provided name in the metadata.
func (u *URL) Path(name string, _ []string, _ string, _ match.Requirements) error {
	u.Data.Set("path", name)
	return nil
}

// Install downloads the file from the URL and processes it based on the provided InstallData.
// It returns the output, the downloaded file, and any error encountered.
func (u *URL) Install(d common.InstallData) (output string, found file.File, err error) {
	return common.Download(d)
}
