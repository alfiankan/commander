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

func runCommanderTUI(shellPath string) {

	homePath := os.Getenv("HOME")
	chartsPath := fmt.Sprintf("%s/.commander", homePath)

	initialModel := cmdr.NewCmdrTUI(chartsPath, shellPath)
	p := tea.NewProgram(initialModel)
	if err := p.Start(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

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
		Long:  `Commander TUI, create, run, share, prompt commands and snippets with ease`,
		Run: func(cmd *cobra.Command, args []string) {

			shellPath, errFlag := cmd.Flags().GetString("shell")

			if errFlag != nil {
				fmt.Println(errFlag)
			}
			runCommanderTUI(shellPath)

		},
	}

	codexCmd := &cobra.Command{
		Use:     "codex",
		Short:   "openai codex",
		Long:    "get ai generated commands and snippets from openai codex",
		Example: "cmdr codex kubectl create simple job",
		Run: func(cmd *cobra.Command, args []string) {
			saveMode, errFlag := cmd.Flags().GetBool("save")
			if errFlag != nil {
				panic(errFlag)
			}
			fmt.Println(args, saveMode)

		},
	}
	codexCmd.Flags().BoolP("save", "s", false, "save generated command")
	rootCmd.AddCommand(codexCmd)

	chartsCmd := &cobra.Command{
		Use:     "get",
		Short:   "get charts online",
		Long:    fmt.Sprintf("get and download charts from %s", cmdr.ChartRepoHost),
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
	chartsCmd.Flags().BoolP("all", "a", false, fmt.Sprintf("download all chart available online from %s", cmdr.ChartRepoHost))
	chartsCmd.Flags().String("host", cmdr.ChartRepoHost, "update chart from defined host")

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
				_, err = f.WriteString(cmdr.ChartTemplate)
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
