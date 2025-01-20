package shell

import (
	"os"
	"strings"
)

type Shell struct {
	cwd  string
	hist []string
}

func New() *Shell {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "?"
	}
	return &Shell{cwd: cwd}
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

	stages := splitPipeline(input)
	if len(stages) == 1 {
		return s.runSingle(stages[0])
	}

	return s.runPipeline(stages)
}
