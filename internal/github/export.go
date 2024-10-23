package github

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	exportMutex sync.Mutex
	exportFile  = "tests/assets2.json"
	importFile  = "tests/assets.json"
)

// Export retrieves the latest release for the repository and stores its assets in a JSON file.
func (g *Repository) Export(release *Release) error {
	exportMutex.Lock()
	defer exportMutex.Unlock()

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(exportFile), 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Load existing data
	var data map[string][]Asset
	if fileData, err := os.ReadFile(exportFile); err == nil {
		if err := json.Unmarshal(fileData, &data); err != nil {
			return fmt.Errorf("failed to unmarshal existing data: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to read existing file: %w", err)
	}

	if data == nil {
		data = make(map[string][]Asset)
	}

	// Add or update the entry
	key := fmt.Sprintf("%s/%s", g.Owner, g.Repo)
	data[key] = release.Assets

	// Save the updated data
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := os.WriteFile(exportFile, jsonData, 0o644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// LatestReleaseFromExport retrieves the latest release information from the exported JSON file.
func (g *Repository) LatestReleaseFromExport() (*Release, error) {
	// Read the exported file
	fileData, err := os.ReadFile(importFile)
	if err != nil {
		return g.LatestRelease()
	}

	var data map[string][]Asset
	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal export data: %w", err)
	}

	key := fmt.Sprintf("%s/%s", g.Owner, g.Repo)
	assets, ok := data[key]
	if !ok {
		return nil, fmt.Errorf("no data found for repository %s", key)
	}

	release := &Release{
		Assets: assets,
	}

	return release, nil
}
