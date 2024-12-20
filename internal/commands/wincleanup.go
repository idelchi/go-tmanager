package commands

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/idelchi/godyl/pkg/file"
)

//go:embed scripts/*
var cleanupFiles embed.FS

type cleanupData struct {
	OldBinary string
	BatchFile string
	Folder    string
	LogFile   string
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

	// Read and parse the template
	tmpl, err := template.ParseFS(cleanupFiles, "scripts/cleanup.bat.template")
	if err != nil {
		return fmt.Errorf("parsing cleanup template: %w", err)
	}

	// Create the batch file
	f, err := os.Create(batchFile.Path())
	if err != nil {
		return fmt.Errorf("creating batch file: %w", err)
	}
	defer f.Close()

	// Execute the template with the data
	data := cleanupData{
		OldBinary: oldBinary.Path(),
		BatchFile: batchFile.Path(),
		Folder:    folder.Path(),
		LogFile:   logFile.Path(),
	}

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("executing cleanup template: %w", err)
	}

	// Fire and forget, run minimized
	cmd := exec.Command("cmd", "/C", "start", "/MIN", batchFile.Path())
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting cleanup script: %w", err)
	}

	return nil
}
