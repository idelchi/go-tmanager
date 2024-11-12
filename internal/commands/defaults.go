package commands

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"

	"github.com/idelchi/godyl/internal/tools"

	"gopkg.in/yaml.v3"
)

//go:embed defaults.yml
var defaultsFile []byte

// Defaults holds all the configuration options for godyl, including tool-specific defaults.
type Defaults struct {
	// Inline tool-specific defaults.
	tools.Defaults `yaml:",inline"`
}

// Unmarshal parses the provided YAML data into the Defaults struct.
func (d *Defaults) Unmarshal(data []byte) error {
	return yaml.Unmarshal(data, d)
}

// FromFile reads and parses a YAML file from the given path into the Defaults struct.
func (d *Defaults) FromFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return d.Unmarshal(data)
}

// Default loads the embedded default YAML configuration.
func (d *Defaults) Default() error {
	return d.Unmarshal(defaultsFile)
}

// Validate checks the Defaults struct to ensure all required fields are properly set.
func (d *Defaults) Validate() error {
	validate := validator.New()
	if err := validate.Struct(d); err != nil {
		return fmt.Errorf("validating Defaults: %w", err)
	}
	return nil
}

// Merge applies values from a Config object into the Defaults struct, only if corresponding values are set.
func (d *Defaults) Merge(cfg Config) (err error) {
	if IsSet("output") {
		d.Output = cfg.Output
	}

	if IsSet("source") {
		d.Source.Type = cfg.Source
	}

	if IsSet("strategy") {
		d.Strategy = cfg.Update.Strategy
	}

	if IsSet("github-token") {
		d.Source.Github.Token = cfg.Tokens.GitHub
	}

	if IsSet("os") {
		err = d.Platform.OS.Parse(cfg.OS)
		d.Platform.Extension = d.Platform.Extension.Default(d.Platform.OS)
		d.Platform.Library = d.Platform.Library.Default(d.Platform.OS, d.Platform.Distribution)
	}

	if IsSet("arch") {
		err = d.Platform.Architecture.Parse(cfg.Arch)
	}

	return err
}

// Load loads configuration defaults from a file or uses embedded defaults if not specified.
func (d *Defaults) Load(path string) error {
	if IsSet("defaults") {
		if err := d.FromFile(path); err != nil {
			return fmt.Errorf("loading defaults from %q: %w", path, err)
		}
	} else {
		if err := d.Default(); err != nil {
			return fmt.Errorf("setting defaults: %w", err)
		}
	}

	if err := d.Initialize(); err != nil {
		return fmt.Errorf("setting tool defaults: %w", err)
	}

	return nil
}
