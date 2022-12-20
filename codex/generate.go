package codex

import (
	"context"
	"fmt"
)

type CmdrCodex struct {
	codexApi OpenApiCodex
}

func NewCmdrCodex(codexApi OpenApiCodex) *CmdrCodex {
	return &CmdrCodex{
		codexApi: codexApi,
	}
}

func (c *CmdrCodex) Run(prompt string) {
	ctx := context.Background()
	result, err := c.codexApi.GetCodexSuggestion(ctx, fmt.Sprintf("cli command %s", prompt))
	if err != nil {
		panic(err)
	}
	for _, c := range result.Choices {
		println("\033[32m" + c.Text)
	}
}
