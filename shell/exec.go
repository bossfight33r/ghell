package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// runSingle runs one command (no pipeline) with optional I/O redirection.
func (s *Shell) runSingle(raw string) error {
	cmd := parseCommand(raw)
	if len(cmd.args) == 0 {
		return nil
	}

	// Check built-ins first.
	if ok, err := s.tryBuiltin(cmd.args); ok {
		return err
	}

	return s.runExternal(cmd, os.Stdin, os.Stdout)
}

// runExternal forks an external process, wiring up the supplied stdin/stdout.
// Explicit redirection in the command struct takes precedence over the passed
// readers/writers.
func (s *Shell) runExternal(cmd command, stdin io.Reader, stdout io.Writer) error {
	c := exec.Command(cmd.args[0], cmd.args[1:]...)
	c.Stderr = os.Stderr

	// --- stdin ---
	if cmd.inputFile != "" {
		f, err := os.Open(cmd.inputFile)
		if err != nil {
			return fmt.Errorf("%s: %w", cmd.args[0], err)
		}
		defer f.Close()
		c.Stdin = f
	} else {
		c.Stdin = stdin
	}

	// --- stdout ---
	if cmd.outputFile != "" {
		f, err := os.Create(cmd.outputFile)
		if err != nil {
			return fmt.Errorf("%s: %w", cmd.args[0], err)
		}
		defer f.Close()
		c.Stdout = f
	} else if cmd.appendFile != "" {
		f, err := os.OpenFile(cmd.appendFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("%s: %w", cmd.args[0], err)
		}
		defer f.Close()
		c.Stdout = f
	} else {
		c.Stdout = stdout
	}

	if err := c.Run(); err != nil {
		// Exit errors are normal (e.g. grep finds nothing → exit 1).
		if _, ok := err.(*exec.ExitError); ok {
			return nil
		}
		return fmt.Errorf("%s: %w", cmd.args[0], err)
	}
	return nil
}
