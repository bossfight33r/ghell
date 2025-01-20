package shell

import (
	"fmt"
	"os"
	"strings"
)

func (s *Shell) tryBuiltin(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}

	switch args[0] {
	case "exit":
		return true, ErrExit
	case "cd":
		return true, s.cd(args[1:])
	case "history":
		return true, s.history()
	case "pwd":
		fmt.Println(s.cwd)
		return true, nil
	}

	return false, nil
}

func (s *Shell) cd(args []string) error {
	var target string

	switch len(args) {
	case 0:
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		target = home
	case 1:
		if args[0] == "-" {
			fmt.Println("cd -: not supported")
			return nil
		}
		target = args[0]
	default:
		return fmt.Errorf("cd: too many args")
	}

	if strings.HasPrefix(target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		target = home + target[1:]
	}

	if err := os.Chdir(target); err != nil {
		return fmt.Errorf("cd: %s", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	s.cwd = cwd
	return nil
}

func (s *Shell) history() error {
	for i, cmd := range s.hist {
		fmt.Printf("%4d  %s\n", i+1, cmd)
	}
	return nil
}
