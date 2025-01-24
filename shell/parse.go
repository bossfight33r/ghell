package shell

type command struct {
	args    []string
	inFile  string
	outFile string
	appFile string
}

func (s *Shell) parseCommand(raw string) command {
	tokens := s.tokenize(raw)
	var cmd command

	i := 0
	for i < len(tokens) {
		switch tokens[i] {
		case ">":
			if i+1 < len(tokens) {
				cmd.outFile = tokens[i+1]
				i += 2
			} else {
				i++
			}
		case ">>":
			if i+1 < len(tokens) {
				cmd.appFile = tokens[i+1]
				i += 2
			} else {
				i++
			}
		case "<":
			if i+1 < len(tokens) {
				cmd.inFile = tokens[i+1]
				i += 2
			} else {
				i++
			}
		default:
			cmd.args = append(cmd.args, tokens[i])
			i++
		}
	}

	return cmd
}
