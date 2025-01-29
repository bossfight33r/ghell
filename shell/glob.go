package shell

import (
	"path/filepath"
	"strings"
)

func expandGlobs(args []string) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if !strings.ContainsAny(a, "*?[") {
			out = append(out, a)
			continue
		}
		matches, err := filepath.Glob(a)
		if err != nil || len(matches) == 0 {
			out = append(out, a)
		} else {
			out = append(out, matches...)
		}
	}
	return out
}
