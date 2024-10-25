// Package goc provides functionality for handling Go-based installations and
// managing Go commands using GitHub repositories. It integrates with GitHub
// to fetch and install Go projects, set paths, and manage metadata related
// to the installation process.
package goc

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/idelchi/godyl/internal/goi"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/pkg/file"
)

// Go represents a Go project sourced from a GitHub repository.
type Go struct {
	github  *github.GitHub
	Command string          `yaml:"command"` // Optional custom command for the Go project
	Data    common.Metadata `yaml:"-"`       // Metadata about the Go project
}

// SetGitHub sets the GitHub repository used for the Go project.
func (g *Go) SetGitHub(gh *github.GitHub) {
	g.github = gh
}

// Get retrieves a specific attribute from the GitHub repository's metadata.
func (g *Go) Get(attribute string) string {
	return g.github.Data.Get(attribute)
}

// Initialize sets up the Go project based on the given name, using the associated GitHub repository.
func (g *Go) Initialize(name string) error {
	return g.github.Initialize(name)
}

// Exe sets the executable for the Go project.
func (g *Go) Exe() error {
	return g.github.Exe()
}

// Version fetches and sets the version for the Go project.
func (g *Go) Version(name string) error {
	return g.github.Version(name)
}

// Path sets the path for the Go project based on its version, using the format github.com/{owner}/{repo}@{version}.
func (g *Go) Path(_ string, _ []string, version string, _ match.Requirements) error {
	g.github.Data.Set("path", fmt.Sprintf("github.com/%s/%s@%s", g.github.Owner, g.github.Repo, version))
	return nil
}

var mu sync.Mutex

// Install installs the Go project by downloading and setting up the required files,
// and returns the output, the found file, and any error encountered during installation.
func (g *Go) Install(d common.InstallData) (output string, found file.File, err error) {
	mu.Lock()
	binary, err := goi.New(d.NoVerifySSL)
	if err != nil {
		return "", "", err
	}
	mu.Unlock()

	installer := goi.Installer{
		Binary: binary,
	}

	var folder file.Folder
	folder.CreateRandomInTempDir()

	installer.Binary.Env.Append(
		goi.Env{
			"GOBIN": folder.Path(),
		},
	)

	name := strings.TrimSuffix(d.Exe, filepath.Ext(d.Exe))

	paths := []string{
		d.Path,
		strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/cmd/%s@", name, name), 1),
		strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/cmd@", name), 1),
	}

	if g.Command != "" {
		paths = []string{
			strings.Replace(d.Path, fmt.Sprintf("/%s@", name), fmt.Sprintf("/%s/%s@", name, g.Command), 1),
		}
	}

	for i, path := range paths {
		paths[i] = strings.ToLower(path)
	}

	for _, path := range paths {
		output, err = installer.Install(path)

		if err == nil {
			d.Path = path
			found, err := common.FindAndSymlink(file.NewFile(folder.Path()), d)

			return output, found, err
		}
	}

	return output, "", err
}
