package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("could not determine home dir: %v", err)
	}

	tests := []struct {
		input string
		want  string
	}{
		{"~/foo.txt", filepath.Join(home, "foo.txt")},
		{"~/a/b/c.txt", filepath.Join(home, "a/b/c.txt")},
		{"/absolute/path.txt", "/absolute/path.txt"},
		{"relative/path.txt", "relative/path.txt"},
		{"~notahome", "~notahome"},
	}

	for _, tt := range tests {
		got := expandPath(tt.input)
		if got != tt.want {
			t.Errorf("expandPath(%q) = %q, want %q", tt.input, got, tt.want)
		}
		if strings.HasPrefix(tt.input, "~/") && strings.HasPrefix(got, "~") {
			t.Errorf("expandPath(%q) still contains tilde: %q", tt.input, got)
		}
	}
}
