package shell

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

var builtinNames = []string{"cd", "pwd", "exit", "history", "export", "jobs"}

type completer struct {
	sh       *Shell
	pathCmds []string
}

func (s *Shell) Completer() readline.AutoCompleter {
	return &completer{sh: s}
}

func (c *completer) Do(line []rune, pos int) ([][]rune, int) {
	lineStr := string(line[:pos])
	fields := strings.Fields(lineStr)
	endsWithSpace := len(lineStr) > 0 && (lineStr[len(lineStr)-1] == ' ' || lineStr[len(lineStr)-1] == '\t')

	var prefix string
	if !endsWithSpace && len(fields) > 0 {
		prefix = fields[len(fields)-1]
	}

	firstWord := len(fields) == 0 || (len(fields) == 1 && !endsWithSpace)

	var candidates []string
	if firstWord {
		candidates = c.commands(prefix)
	} else {
		candidates = c.files(prefix)
	}

	result := make([][]rune, 0, len(candidates))
	for _, cand := range candidates {
		result = append(result, []rune(cand))
	}
	return result, len([]rune(prefix))
}

func (c *completer) commands(prefix string) []string {
	c.loadPathCmds()

	seen := make(map[string]bool)
	var out []string

	for _, b := range builtinNames {
		if strings.HasPrefix(b, prefix) && !seen[b] {
			out = append(out, b)
			seen[b] = true
		}
	}
	for _, cmd := range c.pathCmds {
		if strings.HasPrefix(cmd, prefix) && !seen[cmd] {
			out = append(out, cmd)
			seen[cmd] = true
		}
	}
	return out
}

func (c *completer) loadPathCmds() {
	if c.pathCmds != nil {
		return
	}
	seen := make(map[string]bool)
	for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && !seen[e.Name()] {
				c.pathCmds = append(c.pathCmds, e.Name())
				seen[e.Name()] = true
			}
		}
	}
}

func (c *completer) files(prefix string) []string {
	expanded := prefix
	if prefix == "~" || strings.HasPrefix(prefix, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			if prefix == "~" {
				expanded = home
			} else {
				expanded = home + prefix[1:]
			}
		}
	}

	dir := "."
	base := expanded

	if strings.Contains(expanded, "/") {
		if strings.HasSuffix(expanded, "/") {
			dir = expanded
			base = ""
		} else {
			dir = filepath.Dir(expanded)
			base = filepath.Base(expanded)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var out []string
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), base) {
			continue
		}
		full := e.Name()
		if dir != "." {
			full = filepath.Join(dir, e.Name())
		}
		if e.IsDir() {
			full += "/"
		}
		if strings.HasPrefix(prefix, "~/") {
			if home, err := os.UserHomeDir(); err == nil {
				full = "~" + strings.TrimPrefix(full, home)
			}
		}
		out = append(out, full)
	}
	return out
}
