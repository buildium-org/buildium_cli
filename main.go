package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"buildium_cli/internal/tui"
	"buildium_cli/internal/version"
)

func main() {
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v", "version":
			fmt.Println(version.String())
			return
		}
	}

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "buildium: error:", err)
		os.Exit(1)
	}
}
