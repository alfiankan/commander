package cmdr

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type PrompterTuiModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode textinput.CursorMode
}

func (m PrompterTuiModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PrompterTuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := m.updateInputs(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, cmd
}
func (m *PrompterTuiModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m PrompterTuiModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func NewPropterTuiModel(prompts []ChartPrompt) PrompterTuiModel {
	m := PrompterTuiModel{
		inputs: make([]textinput.Model, 3),
	}

	for i, v := range prompts {

		txtInput := textinput.New()
		txtInput.CursorStyle = cursorStyle
		txtInput.Placeholder = v.Label
		txtInput.CharLimit = 0
		if i == 0 {
			txtInput.Focus()
		}
		m.inputs = append(m.inputs, txtInput)
	}
	//m.inputs[0].Focus()
	return m
}

type Prompter struct {
	tmplt string
	promt []ChartPrompt
}

func NewPrompter(prompt []ChartPrompt) Prompter {
	return Prompter{
		promt: prompt,
	}
}

func (p *Prompter) Render() {
	fmt.Println(p.promt)
	if err := tea.NewProgram(NewPropterTuiModel(p.promt)).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
