package commands

import (
	"fmt"
	"os"
	"os/exec"
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

	if err := update.Apply(body, update.Options{}); err != nil {
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

func winCleanupWorking() error {
	fmt.Println("Issuing a delete command for the old godyl binary")

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)
	var folder file.Folder
	if err := folder.CreateRandomInTempDir(); err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}

	logFile := filepath.Join(exeDir, "godyl_cleanup.log")
	oldBinary := filepath.Join(exeDir, ".godyl.exe.old")
	batchFile := filepath.Join(exeDir, "cleanup.bat")

	batchContent := fmt.Sprintf(`@echo off
echo Started cleanup at %%TIME%% >> "%s"
attrib -h -s "%s" >nul 2>&1

for /L %%%%G in (1,1,10) do (
    del /F "%s" >nul 2>&1
    if not exist "%s" (
        echo Deleted on attempt %%%%G at %%TIME%% >> "%s"
        goto :cleanup_success
    )
    echo Failed attempt %%%%G at %%TIME%% >> "%s"
    timeout /t 1 /nobreak >nul
)

:cleanup_failed
echo Final cleanup failed at %%TIME%% >> "%s"
goto :end

:cleanup_success
echo Cleanup completed successfully at %%TIME%% >> "%s"

:end
del "%s"
`, logFile, oldBinary, oldBinary, oldBinary, logFile, logFile, logFile, logFile, batchFile)

	// Write the batch file
	if err := os.WriteFile(batchFile, []byte(batchContent), 0o644); err != nil {
		return fmt.Errorf("creating batch file: %w", err)
	}

	// Use absolute path for the batch file and run it from its directory
	cmd := exec.Command("cmd", "/C", "cd", "/d", exeDir, "&", "start", "/MIN", filepath.Base(batchFile))
	cmd.Dir = exeDir

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting cleanup script: %w", err)
	}

	return nil
}

func winCleanup() error {
	fmt.Println("Issuing a delete command for the old godyl binary")

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}
	exeDir := file.NewFile(exePath).Dir()

	var folder file.Folder
	if err := folder.CreateRandomInTempDir(); err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}

	oldBinary := file.NewFile(exeDir.Path(), ".godyl.exe.old")
	batchFile := file.NewFile(folder.Path(), "cleanup.bat")
	logFile := file.NewFile(folder.Path(), "cleanup_debug.log")

	fmt.Printf("Batch file stored in: %s\n", batchFile.Path())

	batchContent := fmt.Sprintf(`@echo off
set "OLD_BINARY=%s"
set "BATCH_FILE=%s"
set "FOLDER=%s"
set "LOG_FILE=%s"

echo Starting cleanup script at %%TIME%% >> "%%LOG_FILE%%"

echo Removing attributes >> "%%LOG_FILE%%"
attrib -h -s "%%OLD_BINARY%%" >nul 2>&1

echo Beginning deletion attempts >> "%%LOG_FILE%%"
for /L %%%%G in (1,1,10) do (
    echo Attempt %%%%G >> "%%LOG_FILE%%"
    del /F "%%OLD_BINARY%%" >nul 2>&1
    if not exist "%%OLD_BINARY%%" (
        echo Deleted successfully on attempt %%%%G >> "%%LOG_FILE%%"
        goto :cleanup_success
    )
    echo Still exists after attempt %%%%G >> "%%LOG_FILE%%"
    timeout /t 1 /nobreak >nul
)

:cleanup_failed
echo Cleanup failed after 10 attempts >> "%%LOG_FILE%%"
goto :after_cleanup

:cleanup_success
echo Cleanup succeeded >> "%%LOG_FILE%%"

:after_cleanup
echo Changing directory to TEMP before deleting files and folder >> "%%LOG_FILE%%"
cd /d "%%TEMP%%"

REM echo Deleting batch file >> "%%LOG_FILE%%"
REM del "%%BATCH_FILE%%" >nul 2>&1

echo Removing temporary folder >> "%%LOG_FILE%%"
rmdir /s /q "%%FOLDER%%" >nul 2>&1

echo Cleanup script finished at %%TIME%% >> "%%LOG_FILE%%"
`, oldBinary.Path(), batchFile.Path(), folder.Path(), logFile.Path())

	if err := os.WriteFile(batchFile.Path(), []byte(batchContent), 0o644); err != nil {
		return fmt.Errorf("creating batch file: %w", err)
	}

	// Fire and forget, run minimized
	cmd := exec.Command("cmd", "/C", "start", "/MIN", batchFile.Path())
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting cleanup script: %w", err)
	}

	return nil
}

func winCleanup2() error {
	fmt.Println("Issuing a delete command for the old godyl binary")

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}
	exeDir := file.NewFile(exePath).Dir()

	var folder file.Folder
	if err := folder.CreateRandomInTempDir(); err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}

	oldBinary := file.NewFile(exeDir.Path(), ".godyl.exe.old")
	batchFile := file.NewFile(folder.Path(), "cleanup.bat")

	fmt.Printf("Batch file stored in: %s\n", batchFile.Path())

	// Using absolute paths, no logging, and no printing to a file.
	batchContent := fmt.Sprintf(`@echo off
attrib -h -s "%s" >nul 2>&1

for /L %%%%G in (1,1,10) do (
    del /F "%s" >nul 2>&1
    if not exist "%s" (
        goto :cleanup_success
    )
    timeout /t 1 /nobreak >nul
)

:cleanup_failed
goto :end

:cleanup_success
:end
del "%s"
`, oldBinary.Path(), oldBinary.Path(), oldBinary.Path(), batchFile.Path())

	if err := os.WriteFile(batchFile.Path(), []byte(batchContent), 0o644); err != nil {
		return fmt.Errorf("creating batch file: %w", err)
	}

	// Use 'start /MIN' with absolute paths; no need for cd since all paths are absolute.
	cmd := exec.Command("cmd", "/C", "start", "/MIN", batchFile.Path())
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting cleanup script: %w", err)
	}

	return nil
}
