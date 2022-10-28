package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
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
	textInputs    map[string]*textinput.Model
	wizardState   int
	cursorMode    textinput.CursorMode
	focusIndex    int
	selectedChart int
	orderedKV     []string
	finalCommand  string
}

func (m MainViewModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m MainViewModel) getKeyFromInputTexInputtIndex(idx int) string {
	needle := 0
	for k, _ := range m.textInputs {
		if needle == idx {
			return k
		}
		needle++
	}
	return ""
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

			m.textInputs = map[string]*textinput.Model{}

			for _, v := range promptsChart.chartPrompt {
				txtIn := textinput.New()

				txtIn.Placeholder = v.Label
				txtIn.CharLimit = 0
				txtIn.PromptStyle = focusedStyle
				txtIn.CursorStyle = cursorStyle
				m.textInputs[v.Tmplt] = &txtIn
				m.orderedKV = append(m.orderedKV, v.Tmplt)
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
					fmt.Println(m.parseCommandTemplate(m.list.SelectedItem().(ChartItem).desc))
					cmd := "/bin/bash"
					args := []string{"-c", m.parseCommandTemplate(m.list.SelectedItem().(ChartItem).desc)}
					c := exec.Command(cmd, args...)
					c.Env = os.Environ()

					f, _ := c.CombinedOutput()
					fmt.Println(string(f))

					return m, tea.Quit
				}
				if msg.String() == "up" {
					m.focusIndex--
				} else if msg.String() == "down" {
					m.focusIndex++
				}

				cmds := make([]tea.Cmd, len(m.textInputs))
				for i, v := range m.orderedKV {
					if i == m.focusIndex {
						// Set focused state
						cmds[i] = m.textInputs[v].Focus()
						m.textInputs[v].PromptStyle = focusedStyle
						m.textInputs[v].TextStyle = focusedStyle
						continue
					}
					// Remove focused state
					m.textInputs[v].Blur()
					m.textInputs[v].PromptStyle = noStyle
					m.textInputs[v].TextStyle = noStyle
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

	needle := 0
	for k, textInput := range m.textInputs {
		var mdl textinput.Model
		mdl, cmds[needle] = textInput.Update(msg)
		m.textInputs[k] = &mdl
		needle++
	}

	return tea.Batch(cmds...)
}

func (m MainViewModel) View() string {
	var b strings.Builder

	if m.wizardState == 0 {
		b.WriteString(appStyle.Render(m.list.View()))
	} else if m.wizardState == 1 {
		p := m.list.SelectedItem().(ChartItem)

		b.WriteString(titleStyle.Render("Fill flags or args for") + " -> " + p.title)
		b.WriteString("\n\n")

		m.finalCommand = m.parseCommandTemplate(p.desc)
		b.WriteString(m.finalCommand)

		b.WriteString("\n\n")
		for i, v := range m.orderedKV {
			b.WriteString(m.textInputs[v].View())
			if i < len(m.textInputs)-1 {
				b.WriteRune('\n')
			}
		}

		button := &blurredButton
		if m.focusIndex == len(m.textInputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)

		b.WriteString(helpStyle.Render("move ↑ or ↓"))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("ctrl + d delete current tect input"))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("hit enter on [ Next ] to execute command"))
	}
	return b.String()

}

func (m MainViewModel) parseCommandTemplate(tmplt string) string {

	chart := m.list.SelectedItem().(ChartItem)

	for _, v := range chart.chartPrompt {
		if m.textInputs[v.Tmplt].Value() != "" {
			tmplt = strings.ReplaceAll(tmplt, fmt.Sprintf("{{%s}}", v.Tmplt), m.textInputs[v.Tmplt].Value())
		}
	}
	return tmplt
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

	textInputs := map[string]*textinput.Model{}

	initialModel := MainViewModel{listItem, textInputs, 0, textinput.CursorBlink, 0, 0, []string{}, ""}
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("could not start program:", err)
	}
}
