package codex

import (
	"context"
	"encoding/json"
	"time"
)

type CodexApiMock struct {
	apiKey string
}

func NewCodexApiMock(apiKey string) OpenApiCodex {
	return &CodexApi{
		apiKey: apiKey,
	}
}

func (c *CodexApiMock) GetCodexSuggestion(ctx context.Context, prompt string) (suggestion CodexSuggestion, err error) {

	time.Sleep(2 * time.Second)
	fakeResponse := `{"id":"cmpl-6PRDOx6HwHiL20GqkDuxgsbDB4gq8","object":"text_completion","created":1671520834,"model":"text-davinci-003","choices":[{"text":"\n\nkubectl create job <jobname> --image=<imagename> --replicas=<numberofreplicas>","index":0,"logprobs":null,"finish_reason":"stop"}],"usage":{"prompt_tokens":8,"completion_tokens":30,"total_tokens":38}`

	err = json.Unmarshal([]byte(fakeResponse), &suggestion)

	return
}
