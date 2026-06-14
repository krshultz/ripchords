package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"ripchords/chord"
)

// version is stamped at build time via -ldflags "-X main.version=...".
// When unset (e.g. a bare `go build`), resolveVersion falls back to the
// module version or VCS revision embedded by the Go toolchain.
var version = "dev"

// resolveVersion returns the most specific version string available.
func resolveVersion() string {
	if version != "dev" {
		return version
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return version
	}
	if v := info.Main.Version; v != "" && v != "(devel)" {
		return v
	}
	var rev, dirty string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			rev = s.Value
		case "vcs.modified":
			if s.Value == "true" {
				dirty = "-dirty"
			}
		}
	}
	if rev != "" {
		if len(rev) > 12 {
			rev = rev[:12]
		}
		return rev + dirty
	}
	return version
}

// usageText is the help/usage screen shown for -h, --help, -?, and invalid flags.
func usageText() string {
	return fmt.Sprintf(`ripchords %s — software for guitar players

Usage:
  ripchords [flags]

With no flags, ripchords starts an interactive editor: enter a chord name,
then its fret positions (e.g. "x 3 2 0 1 0" for C major) and it renders the
chord as ASCII tab. Build up a progression and save it to a file.

Flags:
  -v, --version   print the version and exit
  -h, --help      show this help and exit

Editor hotkeys:
  l   show the last chord        ?   open settings
  r   reset the progression      q   quit
`, resolveVersion())
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, path[2:])
		}
	}
	return path
}

type Config struct {
	InputOrder chord.InputOrder `json:"input_order,omitempty"`
	ShowBarre  bool             `json:"show_barre"`
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".ripchords.json"
	}
	return filepath.Join(home, ".ripchords", "config.json")
}

func loadConfig() Config {
	cfg := Config{ShowBarre: true}
	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, &cfg) //nolint
	return cfg
}

func saveConfig(cfg Config) {
	path := configPath()
	os.MkdirAll(filepath.Dir(path), 0755) //nolint
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(path, data, 0644) //nolint
}

func main() {
	flag.Usage = func() { fmt.Fprint(os.Stderr, usageText()) }

	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.BoolVar(&showVersion, "v", false, "print version and exit")
	var showHelp bool
	flag.BoolVar(&showHelp, "help", false, "show help and exit")
	flag.BoolVar(&showHelp, "h", false, "show help and exit")
	flag.BoolVar(&showHelp, "?", false, "show help and exit")
	flag.Parse()
	if showHelp {
		fmt.Print(usageText())
		return
	}
	if showVersion {
		fmt.Printf("Ripchords CLI %s - software for guitar players\n", resolveVersion())
		return
	}

	cfg := loadConfig()
	p := tea.NewProgram(newModel(cfg))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
