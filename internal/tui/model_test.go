package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

var (
	enterKey = tea.KeyMsg{Type: tea.KeyEnter}
	downKey  = tea.KeyMsg{Type: tea.KeyDown}
)

func send(t *testing.T, m Model, msg tea.Msg) (Model, tea.Cmd) {
	t.Helper()
	nm, cmd := m.Update(msg)
	return nm.(Model), cmd
}

// selectGo advances the list to the "go" solution template (index 1) and
// selects it, returning a model parked on the form step.
func selectGo(t *testing.T) Model {
	t.Helper()
	m := New()
	m, _ = send(t, m, downKey) // tutorial -> go
	m, _ = send(t, m, enterKey)
	if m.step != stepForm {
		t.Fatalf("expected stepForm after select, got %v", m.step)
	}
	if m.chosen.Key != "go" {
		t.Fatalf("expected chosen template 'go', got %q", m.chosen.Key)
	}
	// inputs: destination + 3 fields
	if len(m.inputs) != 4 {
		t.Fatalf("expected 4 inputs, got %d", len(m.inputs))
	}
	return m
}

func TestFormRejectsEmptyRequiredFields(t *testing.T) {
	m := selectGo(t)
	// Walk to the last input and submit with everything blank.
	for i := 0; i < len(m.inputs); i++ {
		m, _ = send(t, m, enterKey)
	}
	if m.step != stepForm {
		t.Fatalf("blank form should not advance past stepForm, got %v", m.step)
	}
	if m.formErr == "" {
		t.Fatal("expected a validation error for blank required fields")
	}
}

func TestFullFlowGeneratesProject(t *testing.T) {
	dest := filepath.Join(t.TempDir(), "proj")

	m := selectGo(t)
	m.inputs[0].SetValue(dest)
	m.inputs[1].SetValue("proj-1")
	m.inputs[2].SetValue("my-image")
	m.inputs[3].SetValue("my_harness")
	m.focus = len(m.inputs) - 1 // so the next enter submits rather than advances

	m, _ = send(t, m, enterKey)
	if m.step != stepConfirm {
		t.Fatalf("filled form should advance to confirm, got %v (err=%q)", m.step, m.formErr)
	}

	// Confirm -> kicks off generation; run the returned command synchronously.
	m, cmd := send(t, m, enterKey)
	if m.step != stepResult {
		t.Fatalf("confirm should advance to result, got %v", m.step)
	}
	if cmd == nil {
		t.Fatal("expected a generate command from confirm")
	}
	m, _ = send(t, m, cmd())
	if m.genErr != nil {
		t.Fatalf("generation failed: %v", m.genErr)
	}

	// The project should exist with substituted values.
	data, err := os.ReadFile(filepath.Join(dest, "meta.json"))
	if err != nil {
		t.Fatalf("read generated meta.json: %v", err)
	}
	if !strings.Contains(string(data), "proj-1") {
		t.Errorf("meta.json missing project id: %s", data)
	}
	if !strings.Contains(m.resultView(), "✓") {
		t.Errorf("result view should report success, got: %s", m.resultView())
	}
}

func TestDestinationExistsErrorSurfaced(t *testing.T) {
	dest := t.TempDir() // already exists, and we make it non-empty
	if err := os.WriteFile(filepath.Join(dest, "keep"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	m := selectGo(t)
	m.inputs[0].SetValue(dest)
	m.inputs[1].SetValue("p")
	m.inputs[2].SetValue("i")
	m.inputs[3].SetValue("h")
	m.focus = len(m.inputs) - 1

	m, _ = send(t, m, enterKey)  // -> confirm
	m, cmd := send(t, m, enterKey) // -> result, with generate cmd
	m, _ = send(t, m, cmd())

	if m.genErr == nil {
		t.Fatal("expected a generation error for non-empty destination")
	}
	view := m.resultView()
	if !strings.Contains(view, "already exists") {
		t.Errorf("result view should surface the dest-exists error, got: %s", view)
	}
}
