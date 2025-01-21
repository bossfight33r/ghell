package shell

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func (s *Shell) runPipeline(stages []string) error {
	cmds := make([]command, len(stages))
	for i, stage := range stages {
		cmds[i] = parseCommand(stage)
	}

	readers := make([]*os.File, len(cmds))
	writers := make([]*os.File, len(cmds))

	readers[0] = os.Stdin
	for i := 1; i < len(cmds); i++ {
		r, w, err := os.Pipe()
		if err != nil {
			return fmt.Errorf("pipe: %w", err)
		}
		writers[i-1] = w
		readers[i] = r
	}
	writers[len(cmds)-1] = nil

	procs := make([]*exec.Cmd, len(cmds))
	var toClose []*os.File

	for i, cmd := range cmds {
		if len(cmd.args) == 0 {
			continue
		}

		if ok, _ := s.tryBuiltin(cmd.args); ok {
			fmt.Fprintln(os.Stderr, "builtins in pipelines not supported")
			continue
		}

		p := exec.Command(cmd.args[0], cmd.args[1:]...)
		p.Stderr = os.Stderr

		if cmd.inFile != "" && i == 0 {
			f, err := os.Open(cmd.inFile)
			if err != nil {
				return fmt.Errorf("%s: %w", cmd.args[0], err)
			}
			toClose = append(toClose, f)
			p.Stdin = f
		} else {
			p.Stdin = readers[i]
			if readers[i] != os.Stdin {
				toClose = append(toClose, readers[i])
			}
		}

		if cmd.outFile != "" && i == len(cmds)-1 {
			f, err := os.Create(cmd.outFile)
			if err != nil {
				return fmt.Errorf("%s: %w", cmd.args[0], err)
			}
			toClose = append(toClose, f)
			p.Stdout = f
		} else if cmd.appFile != "" && i == len(cmds)-1 {
			f, err := os.OpenFile(cmd.appFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("%s: %w", cmd.args[0], err)
			}
			toClose = append(toClose, f)
			p.Stdout = f
		} else if writers[i] != nil {
			p.Stdout = writers[i]
			toClose = append(toClose, writers[i])
		} else {
			p.Stdout = os.Stdout
		}

		procs[i] = p
	}

	for _, p := range procs {
		if p == nil {
			continue
		}
		if err := p.Start(); err != nil {
			return fmt.Errorf("%s: %w", p.Path, err)
		}
	}

	for _, f := range toClose {
		f.Close()
	}

	for _, p := range procs {
		if p == nil {
			continue
		}
		if err := p.Wait(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
					fmt.Println()
				}
			} else {
				return fmt.Errorf("%s: %w", p.Path, err)
			}
		}
	}

	return nil
}
