package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ivan/ghell/shell"
)

func main() {
	sh := shell.New()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s $ ", sh.Cwd())

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if err := sh.Execute(input); err != nil {
			if err == shell.ErrExit {
				os.Exit(0)
			}
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
