package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alfiankan/commander/charts"
	"github.com/alfiankan/commander/cmdr"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

const (
	chartRepoHost = "http://localhost:1313/commander-charts/charts"
	chartTemplate = `{
  "kind": "mychart",
  "description": "my personal chart",
  "charts": [
    {
      "usage": "git show log and statistic",
      "cmdt": "git log --stat",
      "type": "cmd",
      "prompt": []
    },
      {
      "usage": "load test apache benchmark",
      "cmdt": "ab -n {{total_req}} -c {{total_concurrent}} {{target_url}} ",
      "type": "cmd",
      "prompt": [
        {
          "tmplt": "total_req",
          "label": "total request",
          "default": "10"
        },
        {
          "tmplt": "total_concurrent",
          "label": "total concurrent",
          "default": "2"
        },
        {
          "tmplt": "target_url",
          "label": "url load test target",
          "default": "https://github.com/"
        }
      ]
    }
	]
}`
)

func runCommanderTUI(shellPath string) {

	homePath := os.Getenv("HOME")
	chartsPath := fmt.Sprintf("%s/.commander", homePath)

	initialModel := cmdr.NewCmdrTUI(chartsPath, shellPath)
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

const (
	greenTemplate = "\033[1;32m%s\033[0m"
	redTemplate   = "\033[1;31m%s\033[0m"
)

func main() {
	homePath := os.Getenv("HOME")
	chartsPath := fmt.Sprintf("%s/.commander", homePath)

	if _, err := os.Stat(chartsPath); os.IsNotExist(err) {
		if err := os.MkdirAll(chartsPath, os.ModePerm); err != nil {
			panic(err)
		}
	}
	rootCmd := &cobra.Command{
		Use:   "cmdr",
		Short: "Commander TUI v0.1.0",
		Long:  `Commander TUI, create, run, share commands and snippets with ease`,
		Run: func(cmd *cobra.Command, args []string) {

			shellPath, errFlag := cmd.Flags().GetString("shell")

			if errFlag != nil {
				fmt.Println(errFlag)
			}
			runCommanderTUI(shellPath)

		},
	}

	chartsCmd := &cobra.Command{
		Use:     "get",
		Short:   "get charts online",
		Long:    fmt.Sprintf("get and download charts from %s", chartRepoHost),
		Example: "cmdr get bitnami-container",
		Run: func(cmd *cobra.Command, args []string) {

			definedChartRepoHost, errFlag := cmd.Flags().GetString("host")
			if errFlag != nil {
				panic(errFlag)
			}

			isAllChart, errFlag := cmd.Flags().GetBool("all")
			if errFlag != nil {
				panic(errFlag)
			}
			if isAllChart {
				charts.NewDownloader([]charts.ChartsRepo{}, definedChartRepoHost)
			} else {
				var chartsL []charts.ChartsRepo
				for _, v := range args {
					chartsL = append(chartsL, charts.ChartsRepo{
						Name: v,
					})
				}
				charts.NewDownloader(chartsL, definedChartRepoHost)
			}
		},
	}
	chartsCmd.Flags().BoolP("all", "a", false, fmt.Sprintf("download all chart available online from %s", chartRepoHost))
	chartsCmd.Flags().String("host", chartRepoHost, "update chart from defined host")

	rootCmd.AddCommand(chartsCmd)

	myChartCmd := &cobra.Command{
		Use:     "mychart",
		Short:   "init or edit your own chart",
		Example: "cmdr mychart",
		Run: func(cmd *cobra.Command, args []string) {
			chartPath := fmt.Sprintf("%s/%s.chart.json", chartsPath, "mychart")
			if _, err := os.Stat(chartPath); os.IsNotExist(err) {
				f, err := os.Create(chartPath)
				if err != nil {
					panic(err)
				}
				_, err = f.WriteString(chartTemplate)
				if err != nil {
					panic(err)
				}
				f.Close()
				fmt.Println("your chart created at ", chartPath)
			} else {
				editor, errFlag := cmd.Flags().GetString("editor")
				if errFlag != nil {
					panic(errFlag)
				}
				if editor != "" {
					editor := cmdr.NewEditorTUI(chartPath, editor)
					p := tea.NewProgram(editor)
					if err := p.Start(); err != nil {
						fmt.Println("could not start program:", err)
					}
				}
			}
		},
	}
	myChartCmd.Flags().StringP("editor", "e", "vim", "open your own chart on editor (code, vim, nvim, nano .etc)")
	rootCmd.AddCommand(myChartCmd)
	rootCmd.AddCommand(&cobra.Command{
		Use:   "tojson",
		Short: "convert multiline text to json for snippet",
		Run: func(cmd *cobra.Command, args []string) {
			text := args[0]
			textJson, _ := json.Marshal(text)
			fmt.Println(string(textJson))
		},
	})

	rootCmd.Flags().StringP("shell", "s", "/bin/bash", "set shell executor path")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Something went wrong", err.Error())
	}
}
