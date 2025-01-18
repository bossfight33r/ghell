// Package shell implements a simple Unix shell with support for pipes,
// I/O redirection, and a command history.
package shell

import (
	"os"
	"strings"
)

// Shell holds the runtime state of the shell session.
type Shell struct {
	cwd     string
	history []string
}

// New creates a Shell initialised to the current working directory.
func New() *Shell {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "?"
	}
	return &Shell{cwd: cwd}
}

// Cwd returns the shell's current working directory.
func (s *Shell) Cwd() string {
	return s.cwd
}

// Execute parses a raw input line and runs it.
func (s *Shell) Execute(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	s.history = append(s.history, input)

	// Split pipeline stages on '|'.
	stages := splitPipeline(input)

	if len(stages) == 1 {
		// No pipe — may still have redirection.
		return s.runSingle(stages[0])
	}

	return s.runPipeline(stages)
}
