package shell

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Shell struct {
	cwd         string
	prevDir     string
	hist        []string
	env         map[string]string
	aliases     map[string]string
	lastRC      int
	exitOnError bool
	jstore      *jobStore
}

func New() *Shell {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "?"
	}
	return &Shell{
		cwd:     cwd,
		env:     map[string]string{},
		aliases: map[string]string{},
		jstore:  newJobStore(),
	}
}

func (s *Shell) Cwd() string { return s.cwd }

func (s *Shell) LoadHistory(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			s.hist = append(s.hist, line)
		}
	}
}

func (s *Shell) Prompt() string {
	ps1 := s.env["PS1"]
	if ps1 == "" {
		ps1 = os.Getenv("PS1")
	}
	if ps1 == "" {
		return s.cwd + " $ "
	}

	home, _ := os.UserHomeDir()
	cwd := s.cwd
	if home != "" && strings.HasPrefix(cwd, home) {
		cwd = "~" + cwd[len(home):]
	}

	result := ps1
	result = strings.ReplaceAll(result, `\w`, cwd)
	result = strings.ReplaceAll(result, `\W`, filepath.Base(cwd))
	if u, err := user.Current(); err == nil {
		result = strings.ReplaceAll(result, `\u`, u.Username)
	}
	if h, err := os.Hostname(); err == nil {
		result = strings.ReplaceAll(result, `\h`, strings.SplitN(h, ".", 2)[0])
	}
	result = strings.ReplaceAll(result, `\$`, "$")
	return s.expand(result)
}

func (s *Shell) Execute(input string) error {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}
	s.hist = append(s.hist, input)

	skip := false
	for _, st := range splitAtOps(input) {
		if st.cmd != "" && !skip {
			if err := s.run(st.cmd); err != nil {
				if err == ErrExit {
					return err
				}
			}
			if s.exitOnError && s.lastRC != 0 {
				return ErrExit
			}
		}
		switch st.op {
		case "&&":
			skip = s.lastRC != 0
		case "||":
			skip = s.lastRC == 0
		default:
			skip = false
		}
	}
	return nil
}

func (s *Shell) run(input string) error {
	if isAssignment(input) {
		k, v, _ := strings.Cut(input, "=")
		s.env[k] = s.expand(v)
		s.lastRC = 0
		return nil
	}

	input = s.expandAlias(input)

	bg := strings.HasSuffix(input, "&")
	if bg {
		input = strings.TrimSpace(input[:len(input)-1])
	}

	stages := splitPipeline(input)
	if len(stages) == 1 {
		return s.runSingle(stages[0], bg)
	}
	return s.runPipeline(stages, bg)
}

func (s *Shell) expandAlias(input string) string {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return input
	}
	expansion, ok := s.aliases[fields[0]]
	if !ok {
		return input
	}
	return expansion + input[len(fields[0]):]
}
