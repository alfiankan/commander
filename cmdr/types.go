package cmdr

type ChartPrompt struct {
	Tmplt        string `json:"tmplt"`
	Label        string `json:"label"`
	Kind         string `json:"kind"`
	DefaultValue string `json:"default"`
}

type Chart struct {
	Usage  string        `json:"usage"`
	Cmdt   string        `json:"cmdt"`
	Type   string        `json:"type"`
	Prompt []ChartPrompt `json:"prompt,omitempty"`
}

type CmdrChart struct {
	Kind        string  `json:"kind"`
	Description string  `json:"description"`
	Charts      []Chart `json:"charts"`
}

type CmdrFinished struct {
	err error
}

type ChartItem struct {
	title, desc string
	chartPrompt []ChartPrompt
	tmplt       string
}

func (i ChartItem) Title() string       { return i.title }
func (i ChartItem) Description() string { return i.desc }
func (i ChartItem) FilterValue() string { return i.title }
