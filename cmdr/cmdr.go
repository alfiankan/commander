package cmdr

import (
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
	Prompt []interface{} `json:"prompt"`
}

type CmdrChart struct {
	Kind        string  `json:"kind"`
	Description string  `json:"description"`
	Charts      []Chart `json:"charts"`
}

type Cmdr struct {
}

func NewCmdr() Cmdr {
	return Cmdr{}
}

func (cr *Cmdr) readCharts(path string) (cmdCharts []CmdrChart) {
	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, v := range entries {
		fmt.Println(v.Name())
		chartB, err := os.ReadFile(fmt.Sprintf("%s/%s", path, v.Name()))
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
