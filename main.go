package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"buildium_cli/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "buildium: error:", err)
		os.Exit(1)
	}
}
