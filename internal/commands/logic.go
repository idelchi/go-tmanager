package commands

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/pretty"
)

// TODO(Idelchi): https://jerrynsh.com/3-easy-ways-to-add-version-flag-in-go/

// embedded struct to hold the embedded default and tools configuration.
type embedded struct {
	// defaults holds the embedded default configuration.
	defaults []byte

	// tools holds the embedded tools configuration.
	tools []byte
}

// App is the main application struct holding configuration and state.
type App struct {
	cfg             Config
	defaults        Defaults
	log             *logger.Logger
	toolsList       tools.Tools
	hasInstallError bool // Indicates if any tool failed to install

	collectedTools []tools.Tool

	version string

	embedded embedded
}

// result struct to hold the tool and any associated error.
type result struct {
	tool  *tools.Tool
	found file.File
	err   error
	msg   string
}

// NewApp initializes a new App instance.
func NewApp(version string, defaults, tools []byte) *App {
	return &App{
		version: version,
		embedded: embedded{
			defaults: defaults,
			tools:    tools,
		},
	}
}

// Run is the main entry point of the application.
func (app *App) Run() error {
	if err := app.initialize(); err != nil {
		return err
	}

	if app.cfg.DumpTools {
		pretty.PrintYAML(PrintYAMLBytes(app.embedded.tools))

		return nil
	}

	if app.cfg.Update.Update {
		return app.processUpdate()
	}

	app.log = logger.New(app.cfg.Log)

	app.logStartupInfo()

	if err := app.loadToolsList(); err != nil {
		return err
	}

	tags, withoutTags := splitTags(app.cfg.Tags)

	if err := app.processTools(tags, withoutTags); err != nil {
		return err
	}

	if app.hasInstallError {
		return errors.New("one or more tools failed to install")
	}

	return nil
}

// initialize parses and validates the configuration and loads defaults.
func (app *App) initialize() error {
	cfg, err := parseFlags(app.version, app.embedded.defaults)
	if err != nil {
		return fmt.Errorf("error parsing flags: %v", err)
	}
	app.cfg = cfg

	if err := app.cfg.Validate(); err != nil {
		return fmt.Errorf("error validating configuration: %v", err)
	}

	app.defaults = Defaults{}
	if err := app.defaults.Load(app.cfg.Defaults.Name(), app.embedded.defaults); err != nil {
		return fmt.Errorf("error loading defaults: %v", err)
	}
	if err := app.defaults.Merge(app.cfg); err != nil {
		return fmt.Errorf("error merging defaults: %v", err)
	}

	return nil
}

// processUpdate handles the update process based on the configuration.
func (app *App) processUpdate() error {
	updater := GodylUpdater{
		Strategy:    app.cfg.Update.Strategy,
		Defaults:    app.defaults.Defaults,
		NoVerifySSL: app.cfg.NoVerifySSL,
	}

	if err := updater.Update(app.version); err != nil {
		return fmt.Errorf("error updating: %v", err)
	}

	return nil
}

// loadToolsList loads the tools configuration from the given path.
func (app *App) loadToolsList() error {
	toolsList, err := app.loadTools(app.cfg.Tools)
	if err != nil {
		return fmt.Errorf("error loading tools: %v", err)
	}
	app.toolsList = toolsList
	return nil
}

// logStartupInfo logs initial startup information.
func (app *App) logStartupInfo() {
	app.log.Info("*** ***")
	app.log.Info("godyl running with:")
	app.log.Info("*** ***")
	app.log.Info("platform:")
	app.log.Info(pretty.YAML(app.defaults.Platform))
	app.log.Info("*** ***")
}

// loadTools loads the tools configuration from the given path.
func (app *App) loadTools(path string) (tools.Tools, error) {
	var toolsList tools.Tools

	if err := toolsList.Load(path); err != nil {
		return toolsList, fmt.Errorf("loading tools from %q: %w", path, err)
	}

	app.log.Info("loaded %d tools from %q", len(toolsList), path)

	return toolsList, nil
}

// processTools processes each tool in the tools list concurrently.
func (app *App) processTools(tags, withoutTags []string) error {
	processor := NewToolProcessor(app)
	return processor.Process(tags, withoutTags)
}

// ToolProcessor handles concurrent processing of tools.
type ToolProcessor struct {
	app       *App
	resultCh  chan result
	errGroup  *errgroup.Group
	waitGroup *sync.WaitGroup
	toolChan  chan tools.Tool
}

// NewToolProcessor creates a new ToolProcessor.
func NewToolProcessor(app *App) *ToolProcessor {
	return &ToolProcessor{
		app:      app,
		resultCh: make(chan result),
		errGroup: &errgroup.Group{},
		toolChan: make(chan tools.Tool),
	}
}

// Process starts processing tools with the given tags.
func (tp *ToolProcessor) Process(tags, withoutTags []string) error {
	tp.setupConcurrencyLimit()

	tp.waitGroup = &sync.WaitGroup{}
	tp.waitGroup.Add(1)

	go tp.collectResults()

	for _, tool := range tp.app.toolsList {
		tp.errGroup.Go(func() error {
			return tp.processTool(&tool, tags, withoutTags)
		})
	}

	if err := tp.errGroup.Wait(); err != nil {
		return fmt.Errorf("error processing tools: %v", err)
	}

	close(tp.resultCh)
	tp.waitGroup.Wait()

	if tp.app.hasInstallError {
		return errors.New("one or more tools failed to install")
	}

	return nil
}

// setupConcurrencyLimit sets the concurrency limit if specified in the config.
func (tp *ToolProcessor) setupConcurrencyLimit() {
	if tp.app.cfg.Parallel > 0 {
		tp.errGroup.SetLimit(tp.app.cfg.Parallel)
		tp.app.log.Info("running with %d parallel downloads", tp.app.cfg.Parallel)
	}
}

// collectResults reads results from the result channel and processes them.
func (tp *ToolProcessor) collectResults() {
	defer tp.waitGroup.Done()
	defer close(tp.toolChan)

	for res := range tp.resultCh {
		tp.app.processResult(res)
	}
}

// processTool processes an individual tool.
func (tp *ToolProcessor) processTool(tool *tools.Tool, tags, withoutTags []string) error {
	tool.ApplyDefaults(tp.app.defaults.Defaults)

	if err := tool.Resolve(tags, withoutTags); err != nil {
		tp.resultCh <- result{tool: tool, err: err}

		return nil
	}

	if tp.app.cfg.Dry {
		tp.resultCh <- result{tool: tool, err: nil}
		return nil
	}

	if tp.app.cfg.NoVerifySSL {
		tool.NoVerifySSL = true
	}

	msg, found, err := tool.Download()
	tp.resultCh <- result{tool: tool, found: found, err: err, msg: msg}

	if err != nil {
		return nil
	}

	output, _, err := tool.Post.Install(common.InstallData{Env: tool.Env})
	if err != nil {
		tp.resultCh <- result{tool: tool, err: err, msg: output}
	}

	return nil
}

// processResult processes the result from a tool operation.
func (app *App) processResult(res result) {
	tool := res.tool
	err := res.err
	msg := res.msg
	found := res.found
	app.log.Info("")
	app.log.Always(tool.Name)
	app.log.Debug("configuration:")
	app.log.Debug("-------")
	app.log.Debug(pretty.YAMLMasked(tool))
	app.log.Debug("-------")
	if err != nil {
		app.handleToolError(tool, err, msg)
	} else {
		app.logToolSuccess(tool, found)
	}
}

// handleToolError logs errors encountered during tool processing.
func (app *App) handleToolError(tool *tools.Tool, err error, msg string) {
	if errors.Is(err, tools.ErrAlreadyExists) ||
		errors.Is(err, tools.ErrDoesNotHaveTags) ||
		errors.Is(err, tools.ErrDoesHaveTags) ||
		errors.Is(err, tools.ErrSkipped) {
		app.log.Warn("  %v", err)
	} else {
		app.hasInstallError = true // Set the flag if a tool fails to install
		app.log.Error("  failed to install")
		app.log.Debug("configuration:")
		app.log.Debug("-------")
		app.log.Debug(pretty.JSONMasked(tool))
		app.log.Debug("-------")
		app.log.Error("  %v", err)
		app.log.Error("  %s", msg)
	}
}

// logToolSuccess logs information about successfully processed tools.
func (app *App) logToolSuccess(tool *tools.Tool, found file.File) {
	if tool.Version != "" {
		app.log.Info("  version: %s", tool.Version)
	}
	app.log.Info("  picked download %q", filepath.Base(tool.Path))
	if tool.Mode == "find" {
		app.log.Info("  picked file %q", found)
		app.log.Info("  installed successfully at %q", filepath.Join(tool.Output, tool.Exe.Name))
		app.logToolAliases(tool)
	} else {
		app.log.Info("  extracted to %q", tool.Output)
	}
}

// logToolAliases logs any aliases for the tool.
func (app *App) logToolAliases(tool *tools.Tool) {
	if tool.Aliases != nil {
		app.log.Info("  symlinks:")
		for _, alias := range tool.Aliases {
			app.log.Info("    - %q", filepath.Join(tool.Output, alias))
		}
	}
}
