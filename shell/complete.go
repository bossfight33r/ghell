package shell

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
)

var builtinNames = []string{"cd", "pwd", "exit", "history", "export", "jobs", "set", "source", "alias", "unalias"}

type completer struct {
	pathCmds []string
}

func (s *Shell) Completer() readline.AutoCompleter {
	return &completer{}
}

func (c *completer) Do(line []rune, pos int) ([][]rune, int) {
	str := string(line[:pos])
	fields := strings.Fields(str)
	endsSpace := len(str) > 0 && (str[len(str)-1] == ' ' || str[len(str)-1] == '\t')

	prefix := ""
	if !endsSpace && len(fields) > 0 {
		prefix = fields[len(fields)-1]
	}

	var matches []string
	if len(fields) == 0 || (len(fields) == 1 && !endsSpace) {
		matches = c.commands(prefix)
	} else {
		matches = c.files(prefix)
	}

	res := make([][]rune, len(matches))
	for i, m := range matches {
		res[i] = []rune(m)
	}
	return res, len([]rune(prefix))
}

func (c *completer) commands(prefix string) []string {
	if c.pathCmds == nil {
		seen := map[string]bool{}
		for _, dir := range filepath.SplitList(os.Getenv("PATH")) {
			ents, err := os.ReadDir(dir)
			if err != nil {
				continue
			}
			for _, e := range ents {
				name := e.Name()
				if !e.IsDir() && !seen[name] {
					seen[name] = true
					c.pathCmds = append(c.pathCmds, name)
				}
			}
		}
	}

	var out []string
	for _, b := range builtinNames {
		if strings.HasPrefix(b, prefix) {
			out = append(out, b)
		}
	}
	for _, cmd := range c.pathCmds {
		if strings.HasPrefix(cmd, prefix) {
			out = append(out, cmd)
		}
	}
	return out
}

func (c *completer) files(prefix string) []string {
	p := prefix
	if strings.HasPrefix(p, "~/") || p == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			if p == "~" {
				p = home
			} else {
				p = home + p[1:]
			}
		}
	}

	dir, base := ".", p
	if strings.Contains(p, "/") {
		if strings.HasSuffix(p, "/") {
			dir, base = p, ""
		} else {
			dir, base = filepath.Dir(p), filepath.Base(p)
		}
	}

	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var out []string
	for _, e := range ents {
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
			if home, _ := os.UserHomeDir(); home != "" {
				full = "~" + strings.TrimPrefix(full, home)
			}
		}
		out = append(out, full)
	}
	return out
}
