package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/term"

	"ripchords/chord"
)

var version = "dev"

type Config struct {
	InputOrder chord.InputOrder `json:"input_order,omitempty"`
}

func configPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".ripchords.json"
	}
	return filepath.Join(home, ".config", "ripchords", "config.json")
}

func loadConfig() Config {
	data, err := os.ReadFile(configPath())
	if err != nil {
		return Config{}
	}
	var cfg Config
	json.Unmarshal(data, &cfg) //nolint
	return cfg
}

func saveConfig(cfg Config) {
	path := configPath()
	os.MkdirAll(filepath.Dir(path), 0755) //nolint
	data, _ := json.MarshalIndent(cfg, "", "  ")
	os.WriteFile(path, data, 0644) //nolint
}

func terminalSize() (width, height int) {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || w <= 0 {
		return 80, 24
	}
	return w, h
}

// paginate prints text one screenful at a time when output is a terminal.
func paginate(text string, linesPerPage int, scanner *bufio.Scanner) {
	lines := strings.Split(strings.TrimRight(text, "\n"), "\n")
	if !term.IsTerminal(int(os.Stdout.Fd())) || linesPerPage < 1 || len(lines) <= linesPerPage {
		fmt.Print(text)
		if !strings.HasSuffix(text, "\n") {
			fmt.Println()
		}
		return
	}
	start := 0
	for start < len(lines) {
		end := start + linesPerPage
		if end > len(lines) {
			end = len(lines)
		}
		fmt.Println(strings.Join(lines[start:end], "\n"))
		start = end
		if start < len(lines) {
			fmt.Print("-- more -- (Enter to continue, q to quit) ")
			scanner.Scan()
			if strings.ToLower(strings.TrimSpace(scanner.Text())) == "q" {
				fmt.Println()
				break
			}
			fmt.Println()
		}
	}
}

func prompt(scanner *bufio.Scanner, msg string) string {
	fmt.Print(msg)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

// promptFrets loops until the user enters valid fret positions or "quit".
// Returns (frets, true) on success, (nil, false) if the user typed "quit".
func promptFrets(scanner *bufio.Scanner, cfg Config) ([]int, bool) {
	for {
		input := prompt(scanner, "Fret positions: ")
		if input == "" || strings.ToLower(input) == "quit" {
			return nil, false
		}
		frets, err := chord.ParseFrets(input, cfg.InputOrder)
		if err != nil {
			fmt.Printf("Invalid input: %s\nPlease try again.\n\n", err)
			continue
		}
		return frets, true
	}
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

	scanner := bufio.NewScanner(os.Stdin)
	cfg := loadConfig()

	if cfg.InputOrder == "" {
		fmt.Println("Welcome to ripchords!")
		fmt.Println()
		fmt.Println("How do you prefer to enter fret positions?")
		fmt.Println("  1) Pitch order (low to high): E A D G B e  —  e.g. \"x 3 2 0 1 0\" for C")
		fmt.Println("  2) String-number order (string 1–6):  e B G D A E  —  e.g. \"0 1 0 2 3 x\" for C")
		for {
			choice := prompt(scanner, "Enter 1 or 2: ")
			switch choice {
			case "1":
				cfg.InputOrder = chord.PitchOrder
			case "2":
				cfg.InputOrder = chord.StringOrder
			default:
				fmt.Println("Please enter 1 or 2.")
				continue
			}
			break
		}
		saveConfig(cfg)
		fmt.Println()
	}

	orderLabel := "pitch order (E A D G B e)"
	if cfg.InputOrder == chord.StringOrder {
		orderLabel = "string-number order (string 1–6: e B G D A E)"
	}
	fmt.Printf("ripchords — entering frets in %s\n", orderLabel)
	fmt.Println("Type 'quit' to exit, 'order' to change input order.")
	fmt.Println()

	for {
		var progression []chord.Chord

		// --- first chord ---
		name := prompt(scanner, "Chord name (or press Enter to skip): ")
		switch strings.ToLower(name) {
		case "quit":
			fmt.Println("Goodbye!")
			return
		case "order":
			if cfg.InputOrder == chord.PitchOrder {
				cfg.InputOrder = chord.StringOrder
				fmt.Println("Switched to string-number order (string 1–6: e B G D A E)")
			} else {
				cfg.InputOrder = chord.PitchOrder
				fmt.Println("Switched to pitch order (E A D G B e)")
			}
			saveConfig(cfg)
			fmt.Println()
			continue
		}

		frets, ok := promptFrets(scanner, cfg)
		if !ok {
			fmt.Println("Goodbye!")
			return
		}
		progression = append(progression, chord.Chord{Name: name, Frets: frets})

		// --- additional chords ---
	moreChords:
		for {
			fmt.Println()
			answer := prompt(scanner, "Another chord? (chord name, fret positions, or no): ")
			lower := strings.ToLower(strings.TrimSpace(answer))

			switch lower {
			case "quit":
				fmt.Println("Goodbye!")
				return
			case "", "n", "no", "done":
				break moreChords
			case "y", "yes":
				chordName := prompt(scanner, "Chord name (or press Enter to skip): ")
				if strings.ToLower(chordName) == "quit" {
					fmt.Println("Goodbye!")
					return
				}
				f, ok := promptFrets(scanner, cfg)
				if !ok {
					fmt.Println("Goodbye!")
					return
				}
				progression = append(progression, chord.Chord{Name: chordName, Frets: f})
			default:
				// Direct fret positions: ask for name, then add.
				// Chord name typed directly: use it as the name, then ask for frets.
				f, err := chord.ParseFrets(answer, cfg.InputOrder)
				var chordName string
				if err != nil {
					// Not fret positions — treat the input as a chord name.
					chordName = answer
				} else {
					chordName = prompt(scanner, "Chord name (or press Enter to skip): ")
					if strings.ToLower(chordName) == "quit" {
						fmt.Println("Goodbye!")
						return
					}
				}
				if f == nil {
					f, ok = promptFrets(scanner, cfg)
					if !ok {
						fmt.Println("Goodbye!")
						return
					}
				}
				progression = append(progression, chord.Chord{Name: chordName, Frets: f})
			}
		}

		// --- render and display ---
		width, height := terminalSize()
		if width > 80 {
			width = 80
		}
		output := chord.RenderProgression(progression, width, true)
		fmt.Println()
		paginate(output, height-3, scanner)
		fmt.Println()

		// --- save to file ---
		saveFile := prompt(scanner, "Save to file? (Enter filename, or press Enter to skip): ")
		if saveFile != "" {
			f, err := os.OpenFile(saveFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Could not open file: %s\n\n", err)
			} else {
				fmt.Fprint(f, chord.RenderProgression(progression, 80, true))
				fmt.Fprintln(f)
				f.Close()
				fmt.Printf("Appended to %s\n\n", saveFile)
			}
		} else {
			fmt.Println()
		}
	}
}
