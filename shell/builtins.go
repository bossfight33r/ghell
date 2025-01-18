package shell

import (
	"fmt"
	"os"
	"strings"
)

// tryBuiltin checks whether the first token is a built-in command and, if so,
// runs it and returns (true, error). Otherwise it returns (false, nil).
func (s *Shell) tryBuiltin(args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}

	switch args[0] {
	case "exit":
		return true, ErrExit

	case "cd":
		return true, s.builtinCd(args[1:])

	case "history":
		return true, s.builtinHistory()

	case "pwd":
		fmt.Println(s.cwd)
		return true, nil
	}

	return false, nil
}

// builtinCd changes the shell's working directory.
func (s *Shell) builtinCd(args []string) error {
	var target string
	switch len(args) {
	case 0:
		// cd with no argument → go home.
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cd: %w", err)
		}
		target = home
	case 1:
		if args[0] == "-" {
			fmt.Println("cd -: not supported yet")
			return nil
		}
		target = args[0]
	default:
		return fmt.Errorf("cd: too many arguments")
	}

	// Expand ~ manually so users can write "cd ~/foo".
	if strings.HasPrefix(target, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("cd: %w", err)
		}
		target = home + target[1:]
	}

	if err := os.Chdir(target); err != nil {
		return fmt.Errorf("cd: %s", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cd: %w", err)
	}
	s.cwd = cwd
	return nil
}

// builtinHistory prints the command history.
func (s *Shell) builtinHistory() error {
	for i, cmd := range s.history {
		fmt.Printf("%4d  %s\n", i+1, cmd)
	}
	return nil
}
