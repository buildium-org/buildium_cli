package tui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"buildium_cli/internal/generator"
)

var (
	appStyle = lipgloss.NewStyle().Margin(1, 2)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))

	focusedLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("205"))

	hintStyle = lipgloss.NewStyle().
			Faint(true)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9"))

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10"))
)

// View implements tea.Model.
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	switch m.step {
	case stepSelect:
		return m.list.View()
	case stepForm:
		return appStyle.Render(m.formView())
	case stepConfirm:
		return appStyle.Render(m.confirmView())
	case stepResult:
		return appStyle.Render(m.resultView())
	}
	return ""
}

func (m Model) formView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Configure: "+m.chosen.Label) + "\n\n")

	for i := range m.inputs {
		label := m.inputLabel(i)
		if i == m.focus {
			b.WriteString(focusedLabelStyle.Render("» "+label) + "\n")
		} else {
			b.WriteString(labelStyle.Render("  "+label) + "\n")
		}
		b.WriteString("  " + m.inputs[i].View() + "\n\n")
	}

	if m.formErr != "" {
		b.WriteString(errorStyle.Render("✗ "+m.formErr) + "\n\n")
	}
	b.WriteString(hintStyle.Render("tab/↑↓ move • enter next/submit • esc back • ctrl+c quit"))
	return b.String()
}

func (m Model) confirmView() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Review") + "\n\n")
	fmt.Fprintf(&b, "  %s %s\n", labelStyle.Render("Template:"), m.chosen.Label)
	fmt.Fprintf(&b, "  %s %s\n", labelStyle.Render("Directory:"), strings.TrimSpace(m.inputs[0].Value()))
	for i, f := range m.chosen.Fields {
		fmt.Fprintf(&b, "  %s %s\n", labelStyle.Render(f.Label+":"), strings.TrimSpace(m.inputs[i+1].Value()))
	}
	b.WriteString("\n" + hintStyle.Render("enter/y generate • esc/n back • ctrl+c quit"))
	return b.String()
}

func (m Model) resultView() string {
	if m.generating {
		return titleStyle.Render("Generating…")
	}
	if m.genErr != nil {
		detail := m.genErr.Error()
		if errors.Is(m.genErr, generator.ErrDestinationNotEmpty) {
			detail = "That directory already exists and isn't empty — choose another name."
		}
		return errorStyle.Render("✗ Generation failed") + "\n\n  " + detail +
			"\n\n" + hintStyle.Render("press any key to quit")
	}
	return successStyle.Render("✓ Created "+m.chosen.Label) +
		"\n\n  Output: " + m.destDir +
		"\n\n" + hintStyle.Render("press any key to quit")
}
