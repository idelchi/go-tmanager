package commands

import (
	"fmt"
	"slices"

	"github.com/go-playground/validator/v10"

	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
)

// Update holds the configuration options for updating the built binary itself.
type Update struct {
	// Strategy to use for updating tools
	Strategy tools.Strategy `mapstructure:"strategy"`
	// Update the tools
	Update bool `mapstructure:"update"`
}

// Tokens holds the configuration options for authentication tokens.
type Tokens struct {
	// GitHub token for authentication
	GitHub string `mapstructure:"github-token" mask:"fixed"`

	// GitLab token for authentication
	GitLab string `mapstructure:"gitlab-token" mask:"fixed"`

	// URL token for authentication
	URL string `mapstructure:"url-token" mask:"fixed"`
}

// Show holds the configuration options for showing various configurations.
type Show struct {
	// Show the parsed configuration and exit
	Config bool `mapstructure:"show-config"`

	// Show the parsed environment variables and exit
	Env bool `mapstructure:"show-env"`

	// Show the parsed default configuration and exit
	Defaults bool `mapstructure:"show-defaults"`

	// Detect the platform and exit
	Platform bool `mapstructure:"show-platform"`
}

// Config holds all the configuration options for godyl.
type Config struct {
	// Show help message and exit
	Help bool
	// Show version information and exit
	Version bool

	// Path to .env file
	DotEnv file.File `mapstructure:"dot-env"`

	// Path to defaults file
	Defaults file.File

	// Show various configurations
	Show Show `mapstructure:",squash"`

	// DumpTools dump out default tools.yml as stdout
	DumpTools bool `mapstructure:"dump-tools"`

	// Update the tool itself
	Update Update `mapstructure:",squash"`

	// Run without making any changes (dry run)
	Dry bool

	// Log level (DEBUG, INFO, WARN, ERROR)
	Log logger.Level

	// Number of parallel downloads (>= 0)
	Parallel int `validate:"gte=0"`

	// Skip SSL verification
	NoVerifySSL bool `mapstructure:"no-verify-ssl"`

	// Path to tools configuration file
	Tools string

	// Output path for the downloaded tools
	Output string

	// Tags to filter tools by
	Tags []string

	// Source from which to install the tools
	Source sources.Type

	// Strategy to use for updating tools
	Strategy tools.Strategy `mapstructure:"strategy"`

	// Operating system to install the tools for
	OS string `mapstructure:"os"`

	// Architecture to install the tools for
	Arch string `mapstructure:"arch"`

	// Tokens for authentication
	Tokens Tokens `mapstructure:",squash"`
}

// Validate checks the configuration for errors.
func (c *Config) Validate() error {
	allowedUpdateStrategies := []tools.Strategy{tools.None, tools.Upgrade, tools.Force}
	if !slices.Contains(allowedUpdateStrategies, c.Update.Strategy) {
		return fmt.Errorf("%w: unknown update strategy: %q: allowed are %v", ErrUsage, c.Update.Strategy, allowedUpdateStrategies)
	}

	if IsSet("config") && !c.Defaults.Exists() {
		return fmt.Errorf("%w: defaults file %q does not exist", ErrUsage, c.Defaults)
	}

	if IsSet("dot-env") && !c.DotEnv.Exists() {
		return fmt.Errorf("%w: dot-env file %q does not exist", ErrUsage, c.DotEnv)
	}

	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
}

var ErrUsage = fmt.Errorf("usage error")
