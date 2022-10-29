package cmdr

import (
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type EditorTuiViewModel struct {
	editorbinary string
	chartPath    string
}

func NewEditorTUI(chartFilePath, editor string) EditorTuiViewModel {

	return EditorTuiViewModel{
		editorbinary: editor,
		chartPath:    chartFilePath,
	}
}
func (m EditorTuiViewModel) Init() tea.Cmd {
	return tea.ExecProcess(exec.Command(m.editorbinary, m.chartPath), func(err error) tea.Msg {
		if err != nil {
			panic(err)
		}
		return tea.Quit
	})
}

func (m EditorTuiViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m EditorTuiViewModel) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Press ctrl+c key to quit"))
	return b.String()

}
