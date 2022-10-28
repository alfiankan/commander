package charts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gojek/heimdall/httpclient"
)

const (
	padding  = 2
	maxWidth = 80
)

var p *tea.Program

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type ChartsRepo struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

func NewDownloader(chartsName []ChartsRepo, cdnHost string) {
	homePath := os.Getenv("HOME")
	chartsPath := fmt.Sprintf("%s/.commander", homePath)

	if len(chartsName) == 0 {
		var onlineChartsList []ChartsRepo

		timeout := 1000 * time.Millisecond
		client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

		res, err := client.Get(fmt.Sprintf("%s/repo.chart.json", cdnHost), nil)
		if err != nil {
			panic(err)
		}
		if res.StatusCode != 200 {
			panic("can't connect to host chart repo")
		}
		body, err := ioutil.ReadAll(res.Body)
		if err := json.Unmarshal(body, &onlineChartsList); err != nil {
			panic(err)
		}
		chartsName = onlineChartsList
	}

	m := DownloaderTuiModel{
		progress:   progress.New(progress.WithDefaultGradient()),
		chartsList: chartsName,
		cdnHost:    cdnHost,
		info:       "Getting charts metadata",
	}

	p = tea.NewProgram(m)

	go func() {
		time.Sleep(5 * time.Second)

		for i, v := range m.chartsList {
			time.Sleep(200 * time.Millisecond)

			timeout := 1000 * time.Millisecond
			client := httpclient.NewClient(httpclient.WithHTTPTimeout(timeout))

			res, err := client.Get(fmt.Sprintf("%s/%s.chart.json", cdnHost, v.Name), nil)
			if err != nil {
				panic(err)
			}

			body, err := ioutil.ReadAll(res.Body)
			if res.StatusCode == 200 {
				os.WriteFile(fmt.Sprintf("%s/%s.chart.json", chartsPath, v.Name), body, os.ModePerm)
			}
			p.Send(progressMsg{
				percent: (float64(100) / float64(len(m.chartsList)-i)) / float64(100),
				info:    fmt.Sprintf("Downloading %s chart", v.Name),
			})
		}
		p.Send(progressMsg{
			percent: 1.0,
			info:    "download completed press any key to exit",
		})

	}()

	if err := p.Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}

}

type progressMsg struct {
	percent float64
	info    string
}

type DownloaderTuiModel struct {
	progress   progress.Model
	chartsList []ChartsRepo
	cdnHost    string
	info       string
}

func (m DownloaderTuiModel) Init() tea.Cmd {

	return nil
}

func (m DownloaderTuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case progressMsg:
		var cmds []tea.Cmd
		m.info = msg.info
		cmds = append(cmds, m.progress.SetPercent(float64(msg.percent)))
		return m, tea.Batch(cmds...)
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd
	}

	return m, nil

}

func (m DownloaderTuiModel) View() string {
	pad := strings.Repeat(" ", padding)
	return pad + helpStyle(m.info) + "\n" +
		pad + m.progress.View() + "\n\n" +
		pad + helpStyle(fmt.Sprintf("Press any key to abort "))
}
