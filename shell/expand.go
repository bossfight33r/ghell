package shell

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func (s *Shell) expand(token string) string {
	if token == "~" {
		if home, err := os.UserHomeDir(); err == nil {
			return home
		}
	}
	if strings.HasPrefix(token, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return home + token[1:]
		}
	}

	if !strings.ContainsRune(token, '$') {
		return token
	}

	var b strings.Builder
	i := 0
	for i < len(token) {
		if token[i] != '$' {
			b.WriteByte(token[i])
			i++
			continue
		}
		i++
		if i >= len(token) {
			b.WriteByte('$')
			break
		}

		// $(cmd) — command substitution
		if token[i] == '(' {
			depth := 1
			j := i + 1
			for j < len(token) && depth > 0 {
				if token[j] == '(' {
					depth++
				} else if token[j] == ')' {
					depth--
				}
				j++
			}
			b.WriteString(s.runCapture(token[i+1 : j-1]))
			i = j
			continue
		}

		// ${VAR}, ${VAR:-default}, ${VAR:+alt}, ${VAR:?msg}, ${#VAR}
		if token[i] == '{' {
			j := i + 1
			depth := 1
			for j < len(token) && depth > 0 {
				if token[j] == '{' {
					depth++
				} else if token[j] == '}' {
					depth--
				}
				j++
			}
			b.WriteString(s.expandBrace(token[i+1 : j-1]))
			i = j
			continue
		}

		// $?
		if token[i] == '?' {
			b.WriteString(strconv.Itoa(s.lastRC))
			i++
			continue
		}

		// $VAR
		j := i
		for j < len(token) {
			ch := rune(token[j])
			if j == i && !unicode.IsLetter(ch) && ch != '_' {
				break
			}
			if j > i && !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
				break
			}
			j++
		}
		if j == i {
			b.WriteByte('$')
			continue
		}
		b.WriteString(s.getVar(token[i:j]))
		i = j
	}
	return b.String()
}

func (s *Shell) getVar(name string) string {
	if val, ok := s.env[name]; ok {
		return val
	}
	return os.Getenv(name)
}

func (s *Shell) expandBrace(expr string) string {
	// ${#VAR} — string length
	if strings.HasPrefix(expr, "#") {
		return strconv.Itoa(len(s.getVar(expr[1:])))
	}

	// ${VAR:-default}, ${VAR:+alt}, ${VAR:?msg}
	for _, op := range []string{":-", ":+", ":?"} {
		idx := strings.Index(expr, op)
		if idx <= 0 {
			continue
		}
		name := expr[:idx]
		arg := expr[idx+len(op):]
		val := s.getVar(name)
		switch op {
		case ":-":
			if val == "" {
				return s.expand(arg)
			}
			return val
		case ":+":
			if val != "" {
				return s.expand(arg)
			}
			return ""
		case ":?":
			if val == "" {
				fmt.Fprintf(os.Stderr, "%s: %s\n", name, arg)
			}
			return val
		}
	}

	return s.getVar(expr)
}

func isAssignment(s string) bool {
	idx := strings.IndexByte(s, '=')
	if idx <= 0 {
		return false
	}
	key := s[:idx]
	if strings.ContainsAny(key, " \t") {
		return false
	}
	for i, ch := range key {
		if i == 0 && !unicode.IsLetter(ch) && ch != '_' {
			return false
		}
		if i > 0 && !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return false
		}
	}
	return true
}
