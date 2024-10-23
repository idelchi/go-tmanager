package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/pkg/pretty"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))
)

type model struct {
	items       []tools.Tool
	cursor      int
	selected    bool
	isLoading   bool
	loadingDots int
	scrollPos   int    // Track scroll position for detail view
	listScroll  int    // Track scroll position for list view
	content     string // Store the formatted content
	height      int    // Store terminal height
}

type tickMsg time.Time

func initialModel() model {
	return model{
		items:       make([]tools.Tool, 0),
		selected:    false,
		isLoading:   true,
		loadingDots: 0,
		scrollPos:   0,
		listScroll:  0,
		height:      getTerminalHeight(),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		return m, nil
	case tickMsg:
		if m.isLoading {
			m.loadingDots = (m.loadingDots + 1) % 4
			return m, tick()
		}
		return m, nil
	case []tools.Tool:
		m.items = msg
		m.isLoading = false
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.selected {
				if m.scrollPos > 0 {
					m.scrollPos--
				}
			} else {
				if m.cursor > 0 {
					m.cursor--
					// Scroll up if cursor moves above visible area
					if m.cursor < m.listScroll {
						m.listScroll = m.cursor
					}
				}
			}
		case "down", "j":
			if m.selected {
				lines := strings.Count(m.content, "\n")
				height := m.height - 5
				if m.scrollPos < lines-height {
					m.scrollPos++
				}
			} else {
				if m.cursor < len(m.items)-1 {
					m.cursor++
					// Scroll down if cursor moves below visible area
					visibleHeight := m.height - 5 // Account for UI chrome
					if m.cursor >= m.listScroll+visibleHeight {
						m.listScroll = m.cursor - visibleHeight + 1
					}
				}
			}
		case "enter":
			if !m.selected {
				m.selected = true
				m.scrollPos = 0
				m.content = pretty.YAMLMasked(m.items[m.cursor])
			} else {
				m.selected = false
			}
		case "esc":
			m.selected = false
			m.scrollPos = 0
		}
	}
	return m, nil
}

func (m model) listView() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Tools\n\n"))

	visibleHeight := m.height - 5 // Account for UI chrome
	end := m.listScroll + visibleHeight
	if end > len(m.items) {
		end = len(m.items)
	}

	// Only render visible items
	for i := m.listScroll; i < end; i++ {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		name := m.items[i].Name
		if i == m.cursor {
			s.WriteString(selectedStyle.Render(fmt.Sprintf("%s %s\n", cursor, name)))
		} else {
			s.WriteString(normalStyle.Render(fmt.Sprintf("%s %s\n", cursor, name)))
		}
	}

	// Add scroll indicators if needed
	if m.listScroll > 0 {
		s.WriteString("\n↑ More items above")
	}
	if end < len(m.items) {
		s.WriteString("\n↓ More items below")
	}

	s.WriteString("\nPress q to quit, enter to view details\n")
	return s.String()
}

func tick() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) View() string {
	if m.isLoading {
		return m.loadingView()
	}
	if m.selected {
		return m.detailView()
	}
	return m.listView()
}

func (m model) loadingView() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Collecting Tools\n\n"))
	dots := strings.Repeat(".", m.loadingDots)
	padding := strings.Repeat(" ", 3-m.loadingDots)
	s.WriteString(fmt.Sprintf("Loading%s%s\n", dots, padding))
	return s.String()
}

func (m model) detailView() string {
	var s strings.Builder
	s.WriteString(titleStyle.Render("Details\n\n"))

	// Split content into lines and apply scrolling
	lines := strings.Split(m.content, "\n")
	height := getTerminalHeight() - 5 // Account for UI chrome
	start := m.scrollPos
	end := start + height
	if end > len(lines) {
		end = len(lines)
	}

	visibleLines := lines[start:end]
	s.WriteString(strings.Join(visibleLines, "\n"))
	s.WriteString("\n\nPress ESC to go back, up/down to scroll\n")
	return s.String()
}

// getTerminalHeight returns the terminal height or a default value
func getTerminalHeight() int {
	// You might want to use termenv or similar to get actual terminal size
	return 20 // Default height, adjust as needed
}

func launchTUI(toolChan chan tools.Tool) error {
	p := tea.NewProgram(initialModel())

	go func() {
		var tools []tools.Tool
		for tool := range toolChan {
			tools = append(tools, tool)
		}
		p.Send(tools)
	}()

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("error running TUI: %v", err)
	}
	return nil
}
