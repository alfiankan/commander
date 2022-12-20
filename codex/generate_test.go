package codex

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCommand(t *testing.T) {

	api := NewCodexApiMock(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	result, err := api.GetCodexSuggestion(ctx, "cli command create kubernetes job")
	fmt.Println(result.Choices[0].Text)
	assert.Nil(t, err)
}
