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
}

func New() *Shell {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "?"
	}
	return &Shell{
		cwd: cwd,
		env: make(map[string]string),
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

	if isAssignment(input) {
		parts := strings.SplitN(input, "=", 2)
		s.env[parts[0]] = s.expand(parts[1])
		s.lastRC = 0
		return nil
	}

	stages := splitPipeline(input)
	if len(stages) == 1 {
		return s.runSingle(stages[0])
	}

	return s.runPipeline(stages)
}
