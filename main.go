package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/ivan/ghell/shell"
)

func main() {
	shell.Init()

	sh := shell.New()

	rl, err := readline.NewEx(&readline.Config{
		AutoComplete: sh.Completer(),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		rl.SetPrompt(sh.Cwd() + " $ ")

		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		}
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if err := sh.Execute(line); err != nil {
			if err == shell.ErrExit {
				os.Exit(0)
			}
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
