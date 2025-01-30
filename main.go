package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/ivan/ghell/shell"
)

func resolveHeredoc(line string, rl *readline.Instance) (string, error) {
	i := 0
	for i < len(line)-1 {
		if line[i] != '<' || line[i+1] != '<' {
			i++
			continue
		}
		// skip <<<
		if i+2 < len(line) && line[i+2] == '<' {
			i++
			continue
		}

		// skip whitespace, find marker word
		j := i + 2
		for j < len(line) && (line[j] == ' ' || line[j] == '\t') {
			j++
		}
		k := j
		for k < len(line) && line[k] != ' ' && line[k] != '\t' {
			k++
		}
		if k == j {
			i++
			continue
		}
		marker := line[j:k]

		var buf strings.Builder
		for {
			rl.SetPrompt("> ")
			hline, err := rl.Readline()
			if err != nil || strings.TrimSpace(hline) == marker {
				break
			}
			buf.WriteString(hline)
			buf.WriteByte('\n')
		}

		tmp, err := os.CreateTemp("", "ghell-heredoc-*")
		if err != nil {
			return line, err
		}
		tmp.WriteString(buf.String())
		tmp.Close()

		return line[:i] + "< " + tmp.Name() + line[k:], nil
	}
	return line, nil
}

func main() {
	shell.Init()

	sh := shell.New()

	home, _ := os.UserHomeDir()
	histFile := filepath.Join(home, ".ghell_history")
	sh.LoadHistory(histFile)

	rl, err := readline.NewEx(&readline.Config{
		AutoComplete: sh.Completer(),
		HistoryFile:  histFile,
		HistoryLimit: 1000,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer rl.Close()

	for {
		rl.SetPrompt(sh.Prompt())

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

		line, err = resolveHeredoc(line, rl)
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
