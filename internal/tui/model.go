// Package tui holds the Bubble Tea interactive program for the Buildium CLI.
//
// This file currently provides only a placeholder root model so the binary
// launches and quits cleanly. The real multi-step wizard (template select ->
// field entry -> confirm -> generate) replaces it in a later task.
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	hintStyle = lipgloss.NewStyle().
			Faint(true)
)

// Model is the root Bubble Tea model.
type Model struct{}

// New constructs the root model.
func New() Model {
	return Model{}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

// View implements tea.Model.
func (m Model) View() string {
	return titleStyle.Render("Buildium CLI") + "\n\n" +
		"Scaffold Buildium tutorials and solution templates.\n\n" +
		hintStyle.Render("Press q to quit.") + "\n"
}
