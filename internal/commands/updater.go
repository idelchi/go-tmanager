package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/inconshreveable/go-update"

	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/file"
)

// GodylUpdater is responsible for updating the godyl tool using the specified update strategy and defaults.
type GodylUpdater struct {
	Strategy    tools.Strategy // Strategy defines how updates are applied (e.g., Upgrade, Downgrade, None).
	Defaults    tools.Defaults // Defaults holds tool-specific default values for the update process.
	NoVerifySSL bool           // NoVerifySSL disables SSL verification for the update process.
}

// Update performs the update process for the godyl tool, applying the specified strategy.
func (gu GodylUpdater) Update(version string) error {
	// Set default strategy if none is provided.
	if gu.Strategy == tools.None {
		gu.Strategy = tools.Upgrade
	}

	// Determine the tool path from build info, defaulting to "idelchi/godyl" if not available.
	path := "idelchi/godyl"
	info, ok := debug.ReadBuildInfo()
	if ok {
		path = strings.TrimPrefix(info.Main.Path, "github.com/")
	}

	// Create a new Tool object with the appropriate strategy and source.
	tool := tools.Tool{
		Name: path,
		Source: sources.Source{
			Type: sources.GITHUB,
		},
		Strategy:    gu.Strategy,
		NoVerifySSL: gu.NoVerifySSL,
	}

	// Apply any default values to the tool.
	tool.ApplyDefaults(gu.Defaults)
	if err := tool.Resolve(nil, nil); err != nil {
		return fmt.Errorf("resolving tool: %w", err)
	}

	if tool.Version.Version == version {
		fmt.Printf("godyl (%v) is already up-to-date\n", version)

		if gu.Strategy == tools.Force {
			fmt.Println("Forcing updating...")
		} else {
			return nil
		}
	}

	fmt.Printf("Update requested from %q -> %q\n", version, tool.Version.Version)

	// Download the tool.
	output, err := gu.Get(tool)

	defer func() {
		folder := file.Folder(output)
		folder.Remove()
	}()

	if err != nil {
		return fmt.Errorf("getting godyl: %w", err)
	}

	// Replace the existing godyl binary with the newly downloaded version.
	if err := gu.Replace(filepath.Join(output, tool.Exe.Name)); err != nil {
		return fmt.Errorf("replacing godyl: %w", err)
	}

	if runtime.GOOS == "windows" {
		if err := winCleanup(); err != nil {
			return fmt.Errorf("issuing delete command: %w", err)
		}
	}

	fmt.Println("Godyl updated successfully")

	return nil
}

// Replace applies the new godyl binary by replacing the current executable with the downloaded one.
func (gu GodylUpdater) Replace(path string) error {
	body, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file %q: %w", path, err)
	}
	defer body.Close()

	options := update.Options{}
	if runtime.GOOS == "windows" {
		// options.OldSavePath = filepath.Join(filepath.Dir(path), ".godyl.exe.old")
	}

	if err := update.Apply(body, options); err != nil {
		return err
	}

	return err
}

// Get downloads the tool based on its source, placing it in a temporary directory, and returns the output path.
func (gu GodylUpdater) Get(tool tools.Tool) (string, error) {
	// Create a temporary directory to store the downloaded tool.
	var dir file.Folder
	// For Windows, get the directory of the current executable.
	if runtime.GOOS == "windows" {
		current, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("getting current executable: %w", err)
		}

		folder := filepath.Dir(current)
		if err := dir.CreateRandomInDir(folder); err != nil {
			return "", fmt.Errorf("creating temporary directory: %w", err)
		}
	} else {
		if err := dir.CreateRandomInTempDir(); err != nil {
			return "", fmt.Errorf("creating temporary directory: %w", err)
		}
	}

	tool.Output = dir.Path()

	// Resolve any dependencies or settings for the tool.
	if err := tool.Resolve(nil, nil); err != nil {
		return "", fmt.Errorf("resolving tool: %w", err)
	}

	// Download the tool and capture any messages or errors.
	if output, msg, err := tool.Download(); err != nil {
		return "", fmt.Errorf("downloading tool: %w: %s: %s", err, output, msg)
	}

	fmt.Printf("Downloading %q from %q\n", tool.Name, tool.Path)

	return tool.Output, nil
}
