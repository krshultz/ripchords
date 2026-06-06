package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"

	"ripchords/chord"
)

var version = "dev"

type Config struct {
	InputOrder chord.InputOrder `json:"input_order,omitempty"`
	ShowBarre  bool             `json:"show_barre"`
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".ripchords.json"
	}
	return filepath.Join(home, ".config", "ripchords", "config.json")
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
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "print version and exit")
	flag.BoolVar(&showVersion, "v", false, "print version and exit")
	flag.Parse()
	if showVersion {
		fmt.Printf("Ripchords CLI %s - software for guitar players\n", version)
		return
	}

	cfg := loadConfig()
	p := tea.NewProgram(newModel(cfg))
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
