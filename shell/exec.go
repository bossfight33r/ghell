package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

func (s *Shell) runSingle(raw string) error {
	cmd := s.parseCommand(raw)
	if len(cmd.args) == 0 {
		return nil
	}

	if ok, err := s.tryBuiltin(cmd.args); ok {
		return err
	}

	return s.runCmd(cmd, os.Stdin, os.Stdout)
}

func (s *Shell) runCmd(cmd command, stdin io.Reader, stdout io.Writer) error {
	p := exec.Command(cmd.args[0], cmd.args[1:]...)
	p.Stderr = os.Stderr

	if cmd.inFile != "" {
		f, err := os.Open(cmd.inFile)
		if err != nil {
			s.lastRC = 1
			return fmt.Errorf("%s: %w", cmd.args[0], err)
		}
		defer f.Close()
		p.Stdin = f
	} else {
		p.Stdin = stdin
	}

	if cmd.outFile != "" {
		f, err := os.Create(cmd.outFile)
		if err != nil {
			s.lastRC = 1
			return fmt.Errorf("%s: %w", cmd.args[0], err)
		}
		defer f.Close()
		p.Stdout = f
	} else if cmd.appFile != "" {
		f, err := os.OpenFile(cmd.appFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			s.lastRC = 1
			return fmt.Errorf("%s: %w", cmd.args[0], err)
		}
		defer f.Close()
		p.Stdout = f
	} else {
		p.Stdout = stdout
	}

	if err := p.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() {
					fmt.Println()
					s.lastRC = 128 + int(status.Signal())
				} else {
					s.lastRC = status.ExitStatus()
				}
			}
			return nil
		}
		s.lastRC = 1
		return fmt.Errorf("%s: %w", cmd.args[0], err)
	}

	s.lastRC = 0
	return nil
}
