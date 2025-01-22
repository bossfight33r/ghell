package shell

import (
	"os"
	"strconv"
	"strings"
	"unicode"
)

func (s *Shell) expand(token string) string {
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

		if token[i] == '?' {
			b.WriteString(strconv.Itoa(s.lastRC))
			i++
			continue
		}

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

		name := token[i:j]
		if val, ok := s.env[name]; ok {
			b.WriteString(val)
		} else {
			b.WriteString(os.Getenv(name))
		}
		i = j
	}
	return b.String()
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
