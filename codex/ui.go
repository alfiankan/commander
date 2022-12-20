package codex

import (
	"fmt"
	"os"

	"github.com/gernest/wow"
	"github.com/gernest/wow/spin"
)

func StartCommanderCodex(query string) bool {

	if os.Getenv("OPENAI_API_KEY") == "" {
		fmt.Println("Please set the OPENAI_API_KEY environment variable on your bash profile (.bashrc, .zshrc, .etc) then refresh the terminal")
		fmt.Println()
		fmt.Println("You can get your API key from https://beta.openai.com/account/api-keys")
		fmt.Println()
		fmt.Println("\033[32mExample : OPENAI_API_KEY=myopenaitoken")

		return false
	}

	w := wow.New(os.Stdout, spin.Get(spin.Clock), "Processing Codex")

	go func() {
		w.Start()
	}()

	codexApi := NewCodexApi(os.Getenv("OPENAI_API_KEY"))
	codexCmdr := NewCmdrCodex(codexApi)
	codexCmdr.Run(query)
	return true
}
