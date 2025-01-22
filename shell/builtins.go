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
		err := s.cd(args[1:])
		s.setRC(err)
		return true, err
	case "export":
		err := s.export(args[1:])
		s.setRC(err)
		return true, err
	case "history":
		return true, s.history()
	case "pwd":
		fmt.Println(s.cwd)
		s.lastRC = 0
		return true, nil
	}

	return false, nil
}

func (s *Shell) setRC(err error) {
	if err != nil {
		s.lastRC = 1
	} else {
		s.lastRC = 0
	}
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

func (s *Shell) export(args []string) error {
	for _, arg := range args {
		if idx := strings.IndexByte(arg, '='); idx > 0 {
			key := arg[:idx]
			val := s.expand(arg[idx+1:])
			s.env[key] = val
			os.Setenv(key, val)
		} else {
			if val, ok := s.env[arg]; ok {
				os.Setenv(arg, val)
			}
		}
	}
	return nil
}

func (s *Shell) history() error {
	for i, cmd := range s.hist {
		fmt.Printf("%4d  %s\n", i+1, cmd)
	}
	return nil
}
