package codex

import (
	"os"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

func StartCommanderCodex(query string) {
	w := wow.New(os.Stdout, spin.Get(spin.Clock), "Processing Codex")

	go func() {
		w.Start()
	}()

	codexApi := NewCodexApi(os.Getenv("OPENAI_API_KEY"))
	codexCmdr := NewCmdrCodex(codexApi)
	codexCmdr.Run(query)

}
