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

func TestResolveVersionPrefersStamped(t *testing.T) {
	orig := version
	defer func() { version = orig }()

	version = "v9.9.9"
	if got := resolveVersion(); got != "v9.9.9" {
		t.Errorf("resolveVersion() = %q, want stamped value %q", got, "v9.9.9")
	}
}

func TestResolveVersionNeverEmpty(t *testing.T) {
	orig := version
	defer func() { version = orig }()

	version = "dev"
	if got := resolveVersion(); got == "" {
		t.Error("resolveVersion() returned empty string")
	}
}

func TestUsageText(t *testing.T) {
	got := usageText()
	for _, want := range []string{"ripchords", "Usage:", "--version", "--help"} {
		if !strings.Contains(got, want) {
			t.Errorf("usageText() missing %q\n---\n%s", want, got)
		}
	}
}
