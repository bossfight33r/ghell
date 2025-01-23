package shell

import (
	"reflect"
	"testing"
)

func newTestShell() *Shell {
	return &Shell{env: make(map[string]string), jstore: newJobStore()}
}

func TestParseCommand(t *testing.T) {
	s := newTestShell()

	tests := []struct {
		input   string
		args    []string
		inFile  string
		outFile string
		appFile string
	}{
		{
			input: "ls -la /tmp",
			args:  []string{"ls", "-la", "/tmp"},
		},
		{
			input:   "echo hello > out.txt",
			args:    []string{"echo", "hello"},
			outFile: "out.txt",
		},
		{
			input:   "echo hi >> log.txt",
			args:    []string{"echo", "hi"},
			appFile: "log.txt",
		},
		{
			input:  "cat < in.txt",
			args:   []string{"cat"},
			inFile: "in.txt",
		},
		{
			input:   "sort < in.txt > out.txt",
			args:    []string{"sort"},
			inFile:  "in.txt",
			outFile: "out.txt",
		},
		{
			input: "echo",
			args:  []string{"echo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cmd := s.parseCommand(tt.input)

			if !reflect.DeepEqual(cmd.args, tt.args) {
				t.Errorf("args: want %v, got %v", tt.args, cmd.args)
			}
			if cmd.inFile != tt.inFile {
				t.Errorf("inFile: want %q, got %q", tt.inFile, cmd.inFile)
			}
			if cmd.outFile != tt.outFile {
				t.Errorf("outFile: want %q, got %q", tt.outFile, cmd.outFile)
			}
			if cmd.appFile != tt.appFile {
				t.Errorf("appFile: want %q, got %q", tt.appFile, cmd.appFile)
			}
		})
	}
}

func TestSplitPipeline(t *testing.T) {
	tests := []struct {
		input  string
		stages []string
	}{
		{"ls", []string{"ls"}},
		{"ls | grep go", []string{"ls", "grep go"}},
		{"ls | grep go | wc -l", []string{"ls", "grep go", "wc -l"}},
		{"  ls  |  grep go  ", []string{"ls", "grep go"}},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := splitPipeline(tt.input)
			if !reflect.DeepEqual(got, tt.stages) {
				t.Errorf("want %v, got %v", tt.stages, got)
			}
		})
	}
}
