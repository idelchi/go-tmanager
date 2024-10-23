package main

import (
	"fmt"
	"reflect"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Item represents any struct with a Name field
type Item interface{}

type model struct {
	items    []Item
	cursor   int
	selected bool
}

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

func initialModel(items []Item) model {
	return model{
		items:    items,
		selected: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if !m.selected && m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if !m.selected && m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = !m.selected
		case "esc":
			m.selected = false
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.selected {
		return m.detailView()
	}
	return m.listView()
}

func (m model) listView() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("Items\n\n"))

	for i, item := range m.items {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		val := reflect.ValueOf(item)
		name := val.FieldByName("Name").String()

		if i == m.cursor {
			s.WriteString(selectedStyle.Render(fmt.Sprintf("%s %s\n", cursor, name)))
		} else {
			s.WriteString(normalStyle.Render(fmt.Sprintf("%s %s\n", cursor, name)))
		}
	}

	s.WriteString("\nPress q to quit\n")
	return s.String()
}

func (m model) detailView() string {
	var s strings.Builder
	item := m.items[m.cursor]
	val := reflect.ValueOf(item)
	typ := val.Type()

	s.WriteString(titleStyle.Render("Details\n\n"))

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		if field.Name != "Name" {
			value := val.Field(i)
			s.WriteString(fmt.Sprintf("%s: %v\n", field.Name, value.Interface()))
		}
	}

	s.WriteString("\nPress ESC to go back\n")
	return s.String()
}

func main() {
	// Example usage with sample data
	type Person struct {
		Name    string
		Age     int
		Email   string
		Country string
	}

	items := []Item{
		Person{Name: "Alice", Age: 30, Email: "alice@example.com", Country: "USA"},
		Person{Name: "Bob", Age: 25, Email: "bob@example.com", Country: "Canada"},
		Person{Name: "Charlie", Age: 35, Email: "charlie@example.com", Country: "UK"},
	}

	p := tea.NewProgram(initialModel(items))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}
