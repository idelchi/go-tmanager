package tools

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/idelchi/go-next-tag/pkg/stdin"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/utils"
	"gopkg.in/yaml.v3"
)

// Tools represents a collection of Tool configurations.
type Tools []Tool

// Load reads a tool configuration file and loads it into the Tools collection.
// If the configuration is not a YAML file, it assumes a tool is being referenced by name or URL and creates a simple tool entry.
func (t *Tools) Load(cfg string) (err error) {
	// Check if the configuration is not a YAML file.
	if !strings.HasSuffix(cfg, ".yml") && !strings.HasSuffix(cfg, ".yaml") && cfg != "-" {
		// If the configuration starts with "http", assume it's a URL.
		if utils.IsURL(cfg) {

			tool := Tool{
				Name: filepath.Base(cfg),
				Path: cfg,
				Mode: Extract,
			}

			tool.Source.Type = sources.DIRECT

			// Create a new Tool with the URL as the Path and Name.
			*t = Tools{
				tool,
			}
		} else {
			// If it's not a URL, treat it as a simple tool name.
			*t = Tools{
				Tool{
					Name: cfg,
					Mode: Extract,
				},
			}
		}
		return nil
	}

	var data []byte

	if cfg == "-" {
		if !stdin.IsPiped() {
			return errors.New("no data piped to stdin")
		}
		input, err := stdin.Read()
		if err != nil {
			return err
		}

		data = []byte(input)
	} else {
		// Read the YAML configuration file from disk.
		input, err := os.ReadFile(cfg)
		if err != nil {
			return err
		}

		data = input
	}

	// Unmarshal the YAML content into the Tools collection.
	err = yaml.Unmarshal(data, t)
	if err != nil {
		return err
	}

	return nil
}
