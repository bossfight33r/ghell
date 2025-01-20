package shell

import (
	"fmt"
	"os"
	"os/exec"
)

// runPipeline connects multiple commands via in-memory pipes.
//
//	ls -la | grep go | wc -l
//
// Each stage's stdout becomes the next stage's stdin. The first stage reads
// from os.Stdin and the last stage writes to os.Stdout.
func (s *Shell) runPipeline(stages []string) error {
	cmds := make([]command, len(stages))
	for i, stage := range stages {
		cmds[i] = parseCommand(stage)
	}

	// Build the chain of os.Pipe() pairs.
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
	writers[len(cmds)-1] = nil // last stage writes to os.Stdout

	// Build and start all processes.
	procs := make([]*exec.Cmd, len(cmds))
	closeAfterStart := make([]*os.File, 0, len(cmds)*2)

	for i, cmd := range cmds {
		if len(cmd.args) == 0 {
			continue
		}

		// Built-ins in a pipeline are not supported in this version.
		if ok, _ := s.tryBuiltin(cmd.args); ok {
			fmt.Fprintf(os.Stderr, "built-ins inside pipelines are not supported\n")
			continue
		}

		c := exec.Command(cmd.args[0], cmd.args[1:]...)
		c.Stderr = os.Stderr

		// stdin
		if cmd.inputFile != "" && i == 0 {
			f, err := os.Open(cmd.inputFile)
			if err != nil {
				return fmt.Errorf("%s: %w", cmd.args[0], err)
			}
			closeAfterStart = append(closeAfterStart, f)
			c.Stdin = f
		} else {
			c.Stdin = readers[i]
			if readers[i] != os.Stdin {
				closeAfterStart = append(closeAfterStart, readers[i])
			}
		}

		// stdout
		if cmd.outputFile != "" && i == len(cmds)-1 {
			f, err := os.Create(cmd.outputFile)
			if err != nil {
				return fmt.Errorf("%s: %w", cmd.args[0], err)
			}
			closeAfterStart = append(closeAfterStart, f)
			c.Stdout = f
		} else if cmd.appendFile != "" && i == len(cmds)-1 {
			f, err := os.OpenFile(cmd.appendFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return fmt.Errorf("%s: %w", cmd.args[0], err)
			}
			closeAfterStart = append(closeAfterStart, f)
			c.Stdout = f
		} else if writers[i] != nil {
			c.Stdout = writers[i]
			closeAfterStart = append(closeAfterStart, writers[i])
		} else {
			c.Stdout = os.Stdout
		}

		procs[i] = c
	}

	// Start all processes.
	for _, c := range procs {
		if c == nil {
			continue
		}
		if err := c.Start(); err != nil {
			return fmt.Errorf("%s: %w", c.Path, err)
		}
	}

	// Close the pipe ends that belong to the parent so children get EOF.
	for _, f := range closeAfterStart {
		f.Close()
	}

	// Wait for all processes.
	for _, c := range procs {
		if c == nil {
			continue
		}
		if err := c.Wait(); err != nil {
			if _, ok := err.(*exec.ExitError); !ok {
				return fmt.Errorf("%s: %w", c.Path, err)
			}
		}
	}

	return nil
}
