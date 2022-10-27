package cmdr

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"encoding/json"
	"fmt"
	"os"
)

type ChartPrompt struct {
	Tmplt string `json:"tmplt"`
	Label string `json:"label"`
	Kind  string `json:"kind"`
}

type Chart struct {
	Usage  string        `json:"usage"`
	Cmdt   string        `json:"cmdt"`
	Help   string        `json:"help"`
	Prompt []ChartPrompt `json:"prompt,omitempty"`
}

type CmdrChart struct {
	Kind        string  `json:"kind"`
	Description string  `json:"description"`
	Charts      []Chart `json:"charts"`
}

type Cmdr struct {
	dotFilePath string
}

func NewCmdr(path string) Cmdr {
	return Cmdr{
		dotFilePath: path,
	}
}

func (cr *Cmdr) readCharts() (cmdCharts []CmdrChart) {
	entries, err := os.ReadDir(cr.dotFilePath)
	if err != nil {
		panic(err)
	}
	for _, v := range entries {
		fmt.Println(v.Name())
		chartB, err := os.ReadFile(fmt.Sprintf("%s/%s", cr.dotFilePath, v.Name()))
		if err != nil {
			panic(err)
		}
		var cmdChart CmdrChart
		if err := json.Unmarshal(chartB, &cmdChart); err != nil {
			panic(err)
		}
		cmdCharts = append(cmdCharts, cmdChart)
	}
	return
}

func (cr *Cmdr) ListViewCharts() {
	var chartsItems []list.Item

	for _, chart := range cr.readCharts() {
		for _, v := range chart.Charts {
			chartsItems = append(chartsItems, item{title: v.Usage, desc: v.Cmdt, chartPrompt: v.Prompt})
		}
	}

	listItem := list.New(chartsItems, list.NewDefaultDelegate(), 0, 0)
	listItem.Styles.StatusBarFilterCount = titleStyle
	listItem.Styles.Title = titleStyle

	m := model{list: listItem}
	m.list.Title = "Commander Charts"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
