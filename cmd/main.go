package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alfiankan/commander/cmdr"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func readCharts() (cmdCharts []cmdr.CmdrChart) {
	entries, err := os.ReadDir("/Users/alfiankan/development/repack/commander/charts")
	if err != nil {
		panic(err)
	}
	for _, v := range entries {
		fmt.Println(v.Name())
		chartB, err := os.ReadFile(fmt.Sprintf("%s/%s", "/Users/alfiankan/development/repack/commander/charts", v.Name()))
		if err != nil {
			panic(err)
		}
		var cmdChart cmdr.CmdrChart
		if err := json.Unmarshal(chartB, &cmdChart); err != nil {
			panic(err)
		}
		cmdCharts = append(cmdCharts, cmdChart)
	}
	return
}

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

type ChartItem struct {
	title, desc string
	chartPrompt []cmdr.ChartPrompt
}

func (i ChartItem) Title() string       { return i.title }
func (i ChartItem) Description() string { return i.desc }
func (i ChartItem) FilterValue() string { return i.title }

type MainViewModel struct {
	list          list.Model
	textInputs    []textinput.Model
	wizardState   int
	cursorMode    textinput.CursorMode
	focusIndex    int
	selectedChart int
}

func (m MainViewModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m MainViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" && m.wizardState == 0 {

			//ShowPropmter(m.list.SelectedItem(), m)
			promptsChart := m.list.SelectedItem().(ChartItem)

			m.textInputs = []textinput.Model{}

			for _, v := range promptsChart.chartPrompt {
				txtIn := textinput.New()

				txtIn.Placeholder = v.Label
				txtIn.PromptStyle = focusedStyle
				txtIn.CursorStyle = cursorStyle
				m.textInputs = append(m.textInputs, txtIn)
			}
			m.selectedChart = m.list.Index()
			m.wizardState = 1
			m.focusIndex = 0
		}

		if msg.String() == "up" || msg.String() == "down" || msg.String() == "enter" || msg.String() == "tab" {
			if m.wizardState == 1 {
				if m.focusIndex > len(m.textInputs) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.textInputs)
				}

				if msg.String() == "enter" && m.focusIndex == len(m.textInputs) {
					return m, tea.Quit
				}
				if msg.String() == "up" {
					m.focusIndex--
				} else if msg.String() == "down" {
					m.focusIndex++
				} else if msg.String() == "enter" && m.textInputs[m.focusIndex].Value() != "" {
					m.focusIndex++
				}

				if msg.String() == "tab" {
					if m.focusIndex >= 0 {
						m.textInputs[m.focusIndex].SetValue("")
					}
				}

				cmds := make([]tea.Cmd, len(m.textInputs))
				for i := 0; i <= len(m.textInputs)-1; i++ {
					if i == m.focusIndex {
						// Set focused state
						cmds[i] = m.textInputs[i].Focus()
						m.textInputs[i].PromptStyle = focusedStyle
						m.textInputs[i].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					m.textInputs[i].Blur()
					m.textInputs[i].PromptStyle = noStyle
					m.textInputs[i].TextStyle = noStyle
				}
				return m, tea.Batch(cmds...)
			}
		}

	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}
	var cmd tea.Cmd

	if m.wizardState == 1 {
		cmd = m.updateInputs(msg)
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *MainViewModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textInputs))

	for i := range m.textInputs {
		m.textInputs[i], cmds[i] = m.textInputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m MainViewModel) View() string {
	if m.wizardState == 0 {
		return appStyle.Render(m.list.View())
	} else if m.wizardState == 1 {
		var b strings.Builder

		for i := range m.textInputs {
			b.WriteString(m.textInputs[i].View())
			if i < len(m.textInputs)-1 {
				b.WriteRune('\n')
			}
		}
		button := &blurredButton
		if m.focusIndex == len(m.textInputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)

		for _, v := range m.textInputs {
			b.WriteString(helpStyle.Render(v.Value()))
		}

		b.WriteString(helpStyle.Render("cursor mode is "))
		b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
		p := m.list.SelectedItem().(ChartItem)
		b.WriteString(p.desc)
		return b.String()
	}
	return "Not Found"
}

func main() {

	var chartsItems []list.Item

	for _, chart := range readCharts() {
		for _, v := range chart.Charts {
			chartsItems = append(chartsItems, ChartItem{title: v.Usage, desc: v.Cmdt, chartPrompt: v.Prompt})
		}
	}

	listItem := list.New(chartsItems, list.NewDefaultDelegate(), 0, 0)
	listItem.Styles.StatusBarFilterCount = titleStyle
	listItem.Styles.Title = titleStyle
	listItem.Title = "Commander charts"

	textInputs := []textinput.Model{}

	initialModel := MainViewModel{listItem, textInputs, 0, textinput.CursorBlink, 0, 0}
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
