package shell

import (
	"fmt"
	"os"
	"strings"
)

func (s *Shell) sourceFile(path string) error {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			path = home + path[1:]
		}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("source: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if err := s.Execute(line); err != nil {
			return err
		}
	}
	return nil
}
