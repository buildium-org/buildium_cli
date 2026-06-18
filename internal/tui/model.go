// Package tui holds the Bubble Tea interactive program for the Buildium CLI:
// a wizard that walks the user through selecting a template, filling in its
// fields, confirming, and generating the project.
package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"buildium_cli/internal/generator"
	"buildium_cli/internal/templates"
)

type step int

const (
	stepSelect step = iota // choose a template
	stepForm               // fill in destination + fields
	stepConfirm            // review summary
	stepResult             // show generation outcome
)

// Model is the root wizard model.
type Model struct {
	step step

	// select
	list   list.Model
	chosen templates.Template

	// form: inputs[0] is the destination directory, inputs[i+1] is field i.
	inputs  []textinput.Model
	focus   int
	formErr string

	// result
	destDir    string
	generating bool
	genErr     error

	quitting bool
}

// genResultMsg carries the outcome of a generator.Generate run.
type genResultMsg struct{ err error }

// New constructs the wizard, seeding the template list from the catalog.
func New() Model {
	items := make([]list.Item, 0, len(templates.Catalog()))
	for _, t := range templates.Catalog() {
		items = append(items, templateItem{t})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a template to generate"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return Model{step: stepSelect, list: l}
}

// Init implements tea.Model.
func (m Model) Init() tea.Cmd { return nil }

// Update implements tea.Model.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.list.SetSize(ws.Width, ws.Height-2)
		return m, nil
	}

	switch m.step {
	case stepSelect:
		return m.updateSelect(msg)
	case stepForm:
		return m.updateForm(msg)
	case stepConfirm:
		return m.updateConfirm(msg)
	case stepResult:
		return m.updateResult(msg)
	}
	return m, nil
}

func (m Model) updateSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if it, ok := m.list.SelectedItem().(templateItem); ok {
				m.chosen = it.t
				cmd := m.buildForm()
				m.step = stepForm
				return m, cmd
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			m.step = stepSelect
			return m, nil
		case "tab", "down":
			return m, m.setFocus(m.focus + 1)
		case "shift+tab", "up":
			return m, m.setFocus(m.focus - 1)
		case "enter":
			if m.focus < len(m.inputs)-1 {
				return m, m.setFocus(m.focus + 1)
			}
			if err := m.validate(); err != "" {
				m.formErr = err
				return m, nil
			}
			m.formErr = ""
			m.step = stepConfirm
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
	return m, cmd
}

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc", "n":
			m.step = stepForm
			return m, m.setFocus(m.focus)
		case "enter", "y":
			m.destDir = strings.TrimSpace(m.inputs[0].Value())
			m.step = stepResult
			m.generating = true
			return m, m.runGenerate()
		}
	}
	return m, nil
}

func (m Model) updateResult(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case genResultMsg:
		m.generating = false
		m.genErr = msg.err
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "enter", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// buildForm constructs the input fields for the chosen template and focuses the
// first one, returning the focus command.
func (m *Model) buildForm() tea.Cmd {
	inputs := make([]textinput.Model, 0, len(m.chosen.Fields)+1)

	dir := textinput.New()
	dir.Prompt = "› "
	dir.Placeholder = m.chosen.Key + "-project"
	inputs = append(inputs, dir)

	for _, f := range m.chosen.Fields {
		ti := textinput.New()
		ti.Prompt = "› "
		if f.Help != "" {
			ti.Placeholder = f.Help
		}
		if f.Default != "" {
			ti.SetValue(f.Default)
		}
		inputs = append(inputs, ti)
	}

	m.inputs = inputs
	m.formErr = ""
	return m.setFocus(0)
}

// setFocus moves keyboard focus to input i (clamped), blurring the rest, and
// returns the cursor-blink command for the newly focused input.
func (m *Model) setFocus(i int) tea.Cmd {
	switch {
	case i < 0:
		i = 0
	case i > len(m.inputs)-1:
		i = len(m.inputs) - 1
	}
	m.focus = i

	var cmd tea.Cmd
	for j := range m.inputs {
		if j == i {
			cmd = m.inputs[j].Focus()
		} else {
			m.inputs[j].Blur()
		}
	}
	return cmd
}

// validate returns the first validation error message, or "" if the form is
// complete. The destination is always required; field requirements come from
// the catalog schema.
func (m Model) validate() string {
	if strings.TrimSpace(m.inputs[0].Value()) == "" {
		return "Destination directory is required"
	}
	for i, f := range m.chosen.Fields {
		if f.Required && strings.TrimSpace(m.inputs[i+1].Value()) == "" {
			return f.Label + " is required"
		}
	}
	return ""
}

// values collects the field inputs into the map the generator expects, keyed by
// templates.Field.Key.
func (m Model) values() map[string]string {
	v := make(map[string]string, len(m.chosen.Fields))
	for i, f := range m.chosen.Fields {
		v[f.Key] = strings.TrimSpace(m.inputs[i+1].Value())
	}
	return v
}

// runGenerate returns a command that materializes the chosen template.
func (m Model) runGenerate() tea.Cmd {
	tmpl, vals, dest := m.chosen, m.values(), m.destDir
	return func() tea.Msg {
		return genResultMsg{err: generator.Generate(tmpl, vals, dest)}
	}
}

// inputLabel returns the human label for input i (0 = destination directory).
func (m Model) inputLabel(i int) string {
	if i == 0 {
		return "Destination directory"
	}
	return m.chosen.Fields[i-1].Label
}
