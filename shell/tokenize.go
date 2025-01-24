package shell

import "strings"

// tokenize splits a raw command string into tokens respecting single and
// double quotes. Variables are expanded in unquoted and double-quoted
// sections; single-quoted content is kept literal.
func (s *Shell) tokenize(raw string) []string {
	var tokens []string
	var cur strings.Builder
	i := 0

	for i < len(raw) {
		ch := raw[i]

		switch {
		case ch == '\'':
			i++
			for i < len(raw) && raw[i] != '\'' {
				cur.WriteByte(raw[i])
				i++
			}
			if i < len(raw) {
				i++
			}

		case ch == '"':
			i++
			var inner strings.Builder
			for i < len(raw) && raw[i] != '"' {
				inner.WriteByte(raw[i])
				i++
			}
			if i < len(raw) {
				i++
			}
			cur.WriteString(s.expand(inner.String()))

		case ch == ' ' || ch == '\t':
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
			i++

		default:
			j := i
			for j < len(raw) && raw[j] != ' ' && raw[j] != '\t' && raw[j] != '\'' && raw[j] != '"' {
				j++
			}
			cur.WriteString(s.expand(raw[i:j]))
			i = j
		}
	}

	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

type step struct {
	cmd string
	op  string
}

// splitAtOps splits input on &&, || and ; operators outside of quotes.
// Single | is left intact for splitPipeline to handle.
func splitAtOps(input string) []step {
	var steps []step
	var cur strings.Builder
	i := 0

	for i < len(input) {
		ch := input[i]

		if ch == '\'' || ch == '"' {
			quote := ch
			cur.WriteByte(ch)
			i++
			for i < len(input) && input[i] != quote {
				cur.WriteByte(input[i])
				i++
			}
			if i < len(input) {
				cur.WriteByte(input[i])
				i++
			}
			continue
		}

		if ch == '&' && i+1 < len(input) && input[i+1] == '&' {
			steps = append(steps, step{cmd: strings.TrimSpace(cur.String()), op: "&&"})
			cur.Reset()
			i += 2
			continue
		}

		if ch == '|' && i+1 < len(input) && input[i+1] == '|' {
			steps = append(steps, step{cmd: strings.TrimSpace(cur.String()), op: "||"})
			cur.Reset()
			i += 2
			continue
		}

		if ch == ';' {
			steps = append(steps, step{cmd: strings.TrimSpace(cur.String()), op: ";"})
			cur.Reset()
			i++
			continue
		}

		cur.WriteByte(ch)
		i++
	}

	if last := strings.TrimSpace(cur.String()); last != "" {
		steps = append(steps, step{cmd: last, op: ""})
	}
	return steps
}

// splitPipeline splits on | outside of quotes.
func splitPipeline(input string) []string {
	var stages []string
	var cur strings.Builder
	i := 0

	for i < len(input) {
		ch := input[i]

		if ch == '\'' || ch == '"' {
			quote := ch
			cur.WriteByte(ch)
			i++
			for i < len(input) && input[i] != quote {
				cur.WriteByte(input[i])
				i++
			}
			if i < len(input) {
				cur.WriteByte(input[i])
				i++
			}
			continue
		}

		if ch == '|' && (i+1 >= len(input) || input[i+1] != '|') {
			if s := strings.TrimSpace(cur.String()); s != "" {
				stages = append(stages, s)
			}
			cur.Reset()
			i++
			continue
		}

		cur.WriteByte(ch)
		i++
	}

	if last := strings.TrimSpace(cur.String()); last != "" {
		stages = append(stages, last)
	}
	return stages
}
