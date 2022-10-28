package main

import (
	"fmt"

	"github.com/alfiankan/commander/cmdr"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	initialModel := cmdr.NewCmdrTUI("/Users/alfiankan/development/repack/commander/charts", "/bin/bash")
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
