package shell

import "strings"

// command is a single executable unit: args + optional redirection targets.
type command struct {
	args       []string // argv[0] is the program name
	inputFile  string   // path for stdin redirection  (<)
	outputFile string   // path for stdout redirection (>)
	appendFile string   // path for stdout append      (>>)
}

// parseCommand converts a single-stage string (no pipes) into a command.
// It handles the redirection operators <, >, and >>.
func parseCommand(raw string) command {
	tokens := strings.Fields(raw)
	cmd := command{}

	i := 0
	for i < len(tokens) {
		tok := tokens[i]
		switch tok {
		case ">":
			if i+1 < len(tokens) {
				cmd.outputFile = tokens[i+1]
				i += 2
			} else {
				i++
			}
		case ">>":
			if i+1 < len(tokens) {
				cmd.appendFile = tokens[i+1]
				i += 2
			} else {
				i++
			}
		case "<":
			if i+1 < len(tokens) {
				cmd.inputFile = tokens[i+1]
				i += 2
			} else {
				i++
			}
		default:
			cmd.args = append(cmd.args, tok)
			i++
		}
	}

	return cmd
}

// splitPipeline splits an input line on unquoted '|' characters.
// Proper quote handling is left as a future exercise — this naive version
// splits on every '|'.
func splitPipeline(input string) []string {
	var stages []string
	for _, s := range strings.Split(input, "|") {
		stage := strings.TrimSpace(s)
		if stage != "" {
			stages = append(stages, stage)
		}
	}
	return stages
}
