package shell

import (
	"fmt"
	"os"
	"sort"
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
		s.lastRC = 0
		if err != nil {
			s.lastRC = 1
		}
		return true, err
	case "export":
		s.export(args[1:])
		s.lastRC = 0
		return true, nil
	case "jobs":
		s.jstore.list()
		s.lastRC = 0
		return true, nil
	case "history":
		for i, c := range s.hist {
			fmt.Printf("%4d  %s\n", i+1, c)
		}
		s.lastRC = 0
		return true, nil
	case "pwd":
		fmt.Println(s.cwd)
		s.lastRC = 0
		return true, nil
	case "set":
		for _, arg := range args[1:] {
			switch arg {
			case "-e":
				s.exitOnError = true
			case "+e":
				s.exitOnError = false
			}
		}
		s.lastRC = 0
		return true, nil
	case "source", ".":
		if len(args) < 2 {
			return true, fmt.Errorf("source: filename argument required")
		}
		s.lastRC = 0
		return true, s.sourceFile(args[1])
	case "alias":
		if len(args) == 1 {
			keys := make([]string, 0, len(s.aliases))
			for k := range s.aliases {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Printf("alias %s='%s'\n", k, s.aliases[k])
			}
			s.lastRC = 0
			return true, nil
		}
		for _, arg := range args[1:] {
			idx := strings.IndexByte(arg, '=')
			if idx > 0 {
				s.aliases[arg[:idx]] = arg[idx+1:]
			}
		}
		s.lastRC = 0
		return true, nil
	case "unalias":
		for _, arg := range args[1:] {
			delete(s.aliases, arg)
		}
		s.lastRC = 0
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

func (s *Shell) export(args []string) {
	for _, arg := range args {
		idx := strings.IndexByte(arg, '=')
		if idx > 0 {
			key, val := arg[:idx], s.expand(arg[idx+1:])
			s.env[key] = val
			os.Setenv(key, val)
			continue
		}
		if val, ok := s.env[arg]; ok {
			os.Setenv(arg, val)
		}
	}
}
