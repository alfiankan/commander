package cmdr

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CmdrTuiViewModel struct {
	list          list.Model
	textInputs    map[string]*textinput.Model
	wizardState   int
	cursorMode    textinput.CursorMode
	focusIndex    int
	selectedChart int
	orderedKV     []string
	finalCommand  string
	unixshell     string
}

func NewCmdrTUI(chartsFilePath, unixshell string) CmdrTuiViewModel {
	var chartsItems []list.Item

	for _, chart := range readCharts(chartsFilePath) {
		for _, v := range chart.Charts {
			if v.Type == "snippet" {
				chartsItems = append(chartsItems, ChartItem{title: v.Usage, desc: "(Snippet)", chartPrompt: v.Prompt, tmplt: v.Cmdt})
			} else if v.Type == "cmd" {
				chartsItems = append(chartsItems, ChartItem{title: v.Usage, desc: fmt.Sprintf("(Command) %s", v.Cmdt), chartPrompt: v.Prompt, tmplt: v.Cmdt})
			}
		}
	}

	listItem := list.New(chartsItems, list.NewDefaultDelegate(), 0, 0)
	listItem.Styles.StatusBarFilterCount = titleStyle
	listItem.Styles.Title = titleStyle
	listItem.Title = "Commander charts v0.1.0"

	textInputs := map[string]*textinput.Model{}
	return CmdrTuiViewModel{listItem, textInputs, 0, textinput.CursorBlink, 0, 0, []string{}, "", unixshell}
}

func readCharts(chartsFilesPath string) (cmdCharts []CmdrChart) {
	entries, err := os.ReadDir(chartsFilesPath)
	if err != nil {
		panic(err)
	}
	for _, v := range entries {
		fmt.Println(v.Name())
		chartB, err := os.ReadFile(fmt.Sprintf("%s/%s", chartsFilesPath, v.Name()))
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

func (m CmdrTuiViewModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CmdrTuiViewModel) getKeyFromInputTexByInputIndex(idx int) string {
	needle := 0
	for k := range m.textInputs {
		if needle == idx {
			return k
		}
		needle++
	}
	return ""
}

func (m CmdrTuiViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" && m.wizardState == 0 {

			promptsChart := m.list.SelectedItem().(ChartItem)

			m.textInputs = map[string]*textinput.Model{}

			for _, v := range promptsChart.chartPrompt {
				txtIn := textinput.New()

				txtIn.Placeholder = fmt.Sprintf("%s default:(%s)", v.Label, v.DefaultValue)
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

		if msg.String() == "up" || msg.String() == "down" || msg.String() == "enter" || msg.String() == "tab" || msg.String() == "ctrl+d" {
			if m.wizardState == 1 {
				if m.focusIndex > len(m.textInputs) {
					m.focusIndex = 0
				} else if m.focusIndex < 0 {
					m.focusIndex = len(m.textInputs)
				}

				if msg.String() == "enter" && m.focusIndex == len(m.textInputs) {
					fmt.Println("Please wait...", m.parseCommandTemplate(m.list.SelectedItem().(ChartItem).desc))
					m.wizardState = 2

					return m, tea.ExecProcess(exec.Command(m.unixshell, "-c", m.parseCommandTemplate(m.list.SelectedItem().(ChartItem).tmplt)), func(err error) tea.Msg {
						if err != nil {
							fmt.Println("Error", err.Error())
						}
						return CmdrFinished{err}
					})
				}
				if msg.String() == "up" {
					m.focusIndex--
				} else if msg.String() == "down" {
					m.focusIndex++
				}

				if msg.String() == "ctrl+d" {
					m.textInputs[m.getKeyFromInputTexByInputIndex(m.focusIndex)].SetValue("")
				}

				cmds := make([]tea.Cmd, len(m.textInputs))
				for i, v := range m.orderedKV {
					if i == m.focusIndex {
						cmds[i] = m.textInputs[v].Focus()
						m.textInputs[v].PromptStyle = focusedStyle
						m.textInputs[v].TextStyle = focusedStyle
						continue
					}
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

func (m *CmdrTuiViewModel) updateInputs(msg tea.Msg) tea.Cmd {
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

func (m CmdrTuiViewModel) View() string {
	var b strings.Builder

	if m.wizardState == 0 {
		b.WriteString(appStyle.Render(m.list.View()))
	} else if m.wizardState == 1 {
		p := m.list.SelectedItem().(ChartItem)

		b.WriteString(titleStyle.Render("Commander") + " -> " + p.title)
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
		b.WriteString(helpStyle.Render("ctrl + d delete current prompt input"))
		b.WriteString("\n")
		b.WriteString(helpStyle.Render("hit enter on [ Next ] to execute command"))
	} else if m.wizardState == 2 {
		b.WriteString(titleStyle.Render("Execution completed ctrl+c to exit"))
	}
	return b.String()

}

func (m CmdrTuiViewModel) parseCommandTemplate(tmplt string) string {

	chart := m.list.SelectedItem().(ChartItem)

	for _, v := range chart.chartPrompt {
		if m.textInputs[v.Tmplt].Value() != "" {
			tmplt = strings.ReplaceAll(tmplt, fmt.Sprintf("{{%s}}", v.Tmplt), m.textInputs[v.Tmplt].Value())
		} else {
			tmplt = strings.ReplaceAll(tmplt, fmt.Sprintf("{{%s}}", v.Tmplt), v.DefaultValue)
		}
	}
	return tmplt
}
