package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idelchi/godyl/internal/tools"
)

var (
	// Styles
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#7D56F4"))

	statusStyle = lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1)

	existsStyle = statusStyle.Copy().
			Background(lipgloss.Color("#28A745")).
			Foreground(lipgloss.Color("#FFFFFF"))

	skipStyle = statusStyle.Copy().
			Background(lipgloss.Color("#FFC107")).
			Foreground(lipgloss.Color("#000000"))

	downloadStyle = statusStyle.Copy().
			Background(lipgloss.Color("#17A2B8")).
			Foreground(lipgloss.Color("#FFFFFF"))

	errorStyle = statusStyle.Copy().
			Background(lipgloss.Color("#DC3545")).
			Foreground(lipgloss.Color("#FFFFFF"))
)

// ToolStatus represents the status of a tool after dry run
type ToolStatus struct {
	Tool    *tools.Tool
	Status  string
	Message string
	Error   error
}

// Model represents the TUI state
type Model struct {
	tools      []ToolStatus
	cursor     int
	width      int
	height     int
	ready      bool
	quitting   bool
	processing bool
	app        *App // Reference to main app
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	// Run the dry run analysis
	return func() tea.Msg {
		// Set dry run flag
		m.app.cfg.Dry = true

		// Run the standard processing
		tags, withoutTags := splitTags(m.app.cfg.Tags)
		processor := NewToolProcessor(m.app)

		var toolStatuses []ToolStatus

		// Create a channel to collect results
		resultCh := make(chan result)

		// Start processing in a goroutine
		go func() {
			for _, tool := range m.app.toolsList {
				tool := tool // capture loop variable
				res := result{}

				tool.ApplyDefaults(m.app.defaults.Defaults)

				if err := tool.Resolve(tags, withoutTags); err != nil {
					res = result{tool: &tool, err: err}
				} else {
					_, found, err := tool.Download()
					res = result{tool: &tool, found: found, err: err}
				}

				resultCh <- res
			}
			close(resultCh)
		}()

		// Collect results
		for res := range resultCh {
			status := ToolStatus{Tool: res.tool}

			switch {
			case res.err == nil:
				status.Status = "download"
				status.Message = "Will be downloaded"
			case tools.IsErrAlreadyExists(res.err):
				status.Status = "exists"
				status.Message = "Already installed"
			case tools.IsErrSkipped(res.err):
				status.Status = "skip"
				status.Message = "Skipped"
			default:
				status.Status = "error"
				status.Message = res.err.Error()
				status.Error = res.err
			}

			toolStatuses = append(toolStatuses, status)
		}

		return toolStatuses
	}
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case []ToolStatus:
		m.tools = msg
		m.processing = false

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.tools)-1 {
				m.cursor++
			}

		case "enter":
			if m.cursor < len(m.tools) {
				// TODO: Implement action on selected tool
			}
		}
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	if m.processing {
		return "Processing tools..."
	}

	var b strings.Builder

	// Header
	b.WriteString(headerStyle.Render("godyl Tool Status\n\n"))

	// Tool list
	maxNameWidth := 0
	for _, t := range m.tools {
		if len(t.Tool.Name) > maxNameWidth {
			maxNameWidth = len(t.Tool.Name)
		}
	}

	for i, t := range m.tools {
		// Prepare the status badge
		var status string
		switch t.Status {
		case "exists":
			status = existsStyle.Render("EXISTS")
		case "skip":
			status = skipStyle.Render("SKIP")
		case "download":
			status = downloadStyle.Render("DOWNLOAD")
		case "error":
			status = errorStyle.Render("ERROR")
		}

		// Create the line
		line := fmt.Sprintf(
			"%-*s %s %s",
			maxNameWidth,
			t.Tool.Name,
			status,
			t.Message,
		)

		// Apply selection styling
		if i == m.cursor {
			b.WriteString(selectedItemStyle.Render(line))
		} else {
			b.WriteString(itemStyle.Render(line))
		}
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	b.WriteString(itemStyle.Render("↑/↓: Navigate • Enter: Select • q: Quit"))

	return b.String()
}

// StartTUI initializes and runs the TUI
func StartTUI(app *App) error {
	model := Model{
		app:        app,
		processing: true,
	}

	p := tea.NewProgram(model)
	_, err := p.Run()

	return err
}
