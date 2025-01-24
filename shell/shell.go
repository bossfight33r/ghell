package shell

import (
	"os"
	"strings"
)

type Shell struct {
	cwd    string
	hist   []string
	env    map[string]string
	lastRC int
	jstore *jobStore
}

func New() *Shell {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "?"
	}
	return &Shell{
		cwd:    cwd,
		env:    make(map[string]string),
		jstore: newJobStore(),
	}
}

func (s *Shell) Cwd() string {
	return s.cwd
}

func (s *Shell) Execute(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	s.hist = append(s.hist, input)

	steps := splitAtOps(input)
	skip := false

	for _, st := range steps {
		if st.cmd != "" && !skip {
			if err := s.runStep(st.cmd); err != nil {
				if err == ErrExit {
					return err
				}
			}
		}

		switch st.op {
		case "&&":
			skip = s.lastRC != 0
		case "||":
			skip = s.lastRC == 0
		case ";", "":
			skip = false
		}
	}

	return nil
}

func (s *Shell) runStep(input string) error {
	if isAssignment(input) {
		parts := strings.SplitN(input, "=", 2)
		s.env[parts[0]] = s.expand(parts[1])
		s.lastRC = 0
		return nil
	}

	bg := strings.HasSuffix(input, "&")
	if bg {
		input = strings.TrimSpace(input[:len(input)-1])
	}

	stages := splitPipeline(input)
	if len(stages) == 1 {
		return s.runSingle(stages[0], bg)
	}
	return s.runPipeline(stages, bg)
}
