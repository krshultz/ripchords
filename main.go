package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const numStrings = 6

type InputOrder string

const (
	PitchOrder  InputOrder = "pitch"
	StringOrder InputOrder = "string_number"
)

type Config struct {
	InputOrder InputOrder `json:"input_order,omitempty"`
}

// stringNames in display order: highest pitch (e) first, lowest (E) last.
var stringNames = [numStrings]string{"e", "B", "G", "D", "A", "E"}

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

// tokenize splits input into per-string tokens.
// If the input contains spaces it splits on whitespace; otherwise each character is one token.
func tokenize(input string) []string {
	input = strings.TrimSpace(input)
	if strings.ContainsAny(input, " \t") {
		return strings.Fields(input)
	}
	tokens := make([]string, 0, len(input))
	for _, ch := range input {
		tokens = append(tokens, string(ch))
	}
	return tokens
}

// parseFrets parses a fret-position string and returns values in pitch order
// (index 0 = string 6 = low E, index 5 = string 1 = high e).
// -1 means muted/not played.
func parseFrets(input string, order InputOrder) ([]int, error) {
	tokens := tokenize(input)
	if len(tokens) != numStrings {
		return nil, fmt.Errorf(
			"expected %d string positions, got %d\n  (e.g. \"x 3 2 0 1 0\" or \"x32010\" for C major)",
			numStrings, len(tokens),
		)
	}
	frets := make([]int, numStrings)
	for i, tok := range tokens {
		if strings.ToLower(tok) == "x" {
			frets[i] = -1
			continue
		}
		n, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("position %d: %q is not a valid fret number or 'x'", i+1, tok)
		}
		if n < 0 || n > 24 {
			return nil, fmt.Errorf("position %d: fret %d is out of range (valid: 0–24)", i+1, n)
		}
		frets[i] = n
	}
	if order == StringOrder {
		// String-number order arrives as [string1…string6] = [e…E]; reverse to pitch order [E…e].
		for i, j := 0, numStrings-1; i < j; i, j = i+1, j-1 {
			frets[i], frets[j] = frets[j], frets[i]
		}
	}
	return frets, nil
}

// detectBarre returns the barre fret if 3+ strings share the minimum fretted position, else 0.
func detectBarre(frets []int) int {
	minFret := 0
	for _, f := range frets {
		if f > 0 && (minFret == 0 || f < minFret) {
			minFret = f
		}
	}
	if minFret == 0 {
		return 0
	}
	count := 0
	for _, f := range frets {
		if f == minFret {
			count++
		}
	}
	if count >= 3 {
		return minFret
	}
	return 0
}

// renderChord returns an ASCII tab diagram. frets must be in pitch order.
func renderChord(name string, frets []int) string {
	var sb strings.Builder
	if name != "" {
		sb.WriteString(fmt.Sprintf("    %s\n", name))
	}
	barre := detectBarre(frets)
	for display := 0; display < numStrings; display++ {
		pitchIdx := numStrings - 1 - display // display 0 (e) = pitch index 5
		label := stringNames[display]
		fret := frets[pitchIdx]
		var marker string
		switch {
		case fret == -1:
			marker = "X"
		case fret == 0:
			marker = "0"
		default:
			marker = strconv.Itoa(fret)
		}
		if barre > 0 && fret > 0 {
			sb.WriteString(label + " |---|-" + marker + "------|\n")
		} else {
			sb.WriteString(label + " |-----" + marker + "------|\n")
		}
	}
	return sb.String()
}

func prompt(scanner *bufio.Scanner, msg string) string {
	fmt.Print(msg)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func main() {
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
				cfg.InputOrder = PitchOrder
			case "2":
				cfg.InputOrder = StringOrder
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
	if cfg.InputOrder == StringOrder {
		orderLabel = "string-number order (string 1–6: e B G D A E)"
	}
	fmt.Printf("ripchords — entering frets in %s\n", orderLabel)
	fmt.Println("Type 'quit' to exit, 'order' to change input order.")
	fmt.Println()

	for {
		name := prompt(scanner, "Chord name (or press Enter to skip): ")
		switch strings.ToLower(name) {
		case "quit":
			fmt.Println("Goodbye!")
			return
		case "order":
			if cfg.InputOrder == PitchOrder {
				cfg.InputOrder = StringOrder
				fmt.Println("Switched to string-number order (string 1–6: e B G D A E)")
			} else {
				cfg.InputOrder = PitchOrder
				fmt.Println("Switched to pitch order (E A D G B e)")
			}
			saveConfig(cfg)
			fmt.Println()
			continue
		}

		var frets []int
		for {
			input := prompt(scanner, "Fret positions: ")
			if strings.ToLower(input) == "quit" {
				fmt.Println("Goodbye!")
				return
			}
			var err error
			frets, err = parseFrets(input, cfg.InputOrder)
			if err != nil {
				fmt.Printf("Invalid input: %s\nPlease try again.\n\n", err)
				continue
			}
			break
		}

		diagram := renderChord(name, frets)
		fmt.Println()
		fmt.Print(diagram)
		fmt.Println()

		saveFile := prompt(scanner, "Save to file? (Enter filename, or press Enter to skip): ")
		if saveFile != "" {
			f, err := os.OpenFile(saveFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("Could not open file: %s\n\n", err)
			} else {
				fmt.Fprint(f, diagram)
				fmt.Fprintln(f)
				f.Close()
				fmt.Printf("Appended to %s\n\n", saveFile)
			}
		} else {
			fmt.Println()
		}
	}
}
