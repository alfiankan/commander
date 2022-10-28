package cmdr

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
	noStyle             = lipgloss.NewStyle()
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("#25A065"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedButton       = focusedStyle.Copy().Render("[ Next ]")
	blurredButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Next"))
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	cursorStyle         = focusedStyle.Copy()
)
