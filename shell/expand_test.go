package shell

import "testing"

func TestExpand(t *testing.T) {
	s := newTestShell()
	s.env["NAME"] = "Ivan"
	s.env["DIR"] = "/tmp"
	s.lastRC = 42

	tests := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		{"$NAME", "Ivan"},
		{"hello $NAME", "hello Ivan"},
		{"$DIR/file.txt", "/tmp/file.txt"},
		{"$?", "42"},
		{"exit code: $?", "exit code: 42"},
		{"$UNDEFINED", ""},
		{"no dollar", "no dollar"},
		{"$", "$"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := s.expand(tt.input)
			if got != tt.want {
				t.Errorf("expand(%q): want %q, got %q", tt.input, tt.want, got)
			}
		})
	}
}

func TestIsAssignment(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"NAME=Ivan", true},
		{"_VAR=123", true},
		{"MY_VAR=hello world", true},
		{"PATH=/usr/bin:/bin", true},
		{"ls -la", false},
		{"=value", false},
		{"2BAD=nope", false},
		{"no equals", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := isAssignment(tt.input)
			if got != tt.want {
				t.Errorf("isAssignment(%q): want %v, got %v", tt.input, tt.want, got)
			}
		})
	}
}
