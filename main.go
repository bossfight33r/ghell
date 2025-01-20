package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ivan/ghell/shell"
)

func main() {
	sh := shell.New()
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s $ ", sh.Cwd())

		line, err := r.ReadString('\n')
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
