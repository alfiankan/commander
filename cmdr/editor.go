package cmdr

import (
	"fmt"
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
			fmt.Println("Error", err.Error())
		}
		return CmdrFinished{err}
	})
}

func (m EditorTuiViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	return m, tea.ExecProcess(exec.Command(m.editorbinary, m.chartPath), func(err error) tea.Msg {
		if err != nil {
			fmt.Println("Error", err.Error())
		}
		return CmdrFinished{err}
	})
}

func (m EditorTuiViewModel) View() string {
	var b strings.Builder

	return b.String()

}
