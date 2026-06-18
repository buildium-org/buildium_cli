package tui

import (
	"fmt"

	"buildium_cli/internal/templates"
)

// templateItem adapts a templates.Template to the bubbles/list item interface.
type templateItem struct {
	t templates.Template
}

func (i templateItem) Title() string { return i.t.Label }

func (i templateItem) Description() string {
	return fmt.Sprintf("%d field(s) • key: %s", len(i.t.Fields), i.t.Key)
}

// FilterValue is unused (filtering is disabled) but required by list.Item.
func (i templateItem) FilterValue() string { return i.t.Label }
