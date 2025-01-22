package shell

import "strings"

type command struct {
	args    []string
	inFile  string
	outFile string
	appFile string
}

func (s *Shell) parseCommand(raw string) command {
	tokens := strings.Fields(raw)
	var cmd command

	i := 0
	for i < len(tokens) {
		switch tokens[i] {
		case ">":
			if i+1 < len(tokens) {
				cmd.outFile = s.expand(tokens[i+1])
				i += 2
			} else {
				i++
			}
		case ">>":
			if i+1 < len(tokens) {
				cmd.appFile = s.expand(tokens[i+1])
				i += 2
			} else {
				i++
			}
		case "<":
			if i+1 < len(tokens) {
				cmd.inFile = s.expand(tokens[i+1])
				i += 2
			} else {
				i++
			}
		default:
			cmd.args = append(cmd.args, s.expand(tokens[i]))
			i++
		}
	}

	return cmd
}

func splitPipeline(input string) []string {
	var out []string
	for _, s := range strings.Split(input, "|") {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
