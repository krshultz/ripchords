package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/term"
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

type Chord struct {
	name  string
	frets []int
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

// renderChord returns an ASCII tab diagram for a single chord. frets must be in pitch order.
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

// centerInWidth centers s within a field of w chars, truncating if s is too long.
func centerInWidth(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	total := w - len(s)
	left := total / 2
	right := total - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// renderProgression renders chords side-by-side, wrapping rows to fit within width.
// Each segment is 14 chars wide ("|-----0------|") with a 1-space gap between chords.
// Total row width for N chords: 2 (label) + 14*N + (N-1) = 15N+1.
func renderProgression(chords []Chord, width int) string {
	const segWidth = 14 // "|-----0------|"
	const colWidth = 15 // segWidth + 1 space separator

	chordsPerRow := (width - 1) / colWidth
	if chordsPerRow < 1 {
		chordsPerRow = 1
	}

	var sb strings.Builder
	for start := 0; start < len(chords); start += chordsPerRow {
		end := start + chordsPerRow
		if end > len(chords) {
			end = len(chords)
		}
		row := chords[start:end]

		// Name row: center each name within its 14-char column.
		sb.WriteString("  ")
		for i, ch := range row {
			centered := centerInWidth(ch.name, segWidth)
			sb.WriteString(centered)
			if i < len(row)-1 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")

		// String rows.
		for display := 0; display < numStrings; display++ {
			pitchIdx := numStrings - 1 - display
			label := stringNames[display]
			sb.WriteString(label + " ")
			for i, ch := range row {
				fret := ch.frets[pitchIdx]
				barre := detectBarre(ch.frets)
				var marker string
				switch {
				case fret == -1:
					marker = "X"
				case fret == 0:
					marker = "0"
				default:
					marker = strconv.Itoa(fret)
				}
				var seg string
				if barre > 0 && fret > 0 {
					seg = "|---|-" + marker + "------|"
				} else {
					seg = "|-----" + marker + "------|"
				}
				sb.WriteString(seg)
				if i < len(row)-1 {
					sb.WriteString(" ")
				}
			}
			sb.WriteString("\n")
		}

		if end < len(chords) {
			sb.WriteString("\n")
		}
	}
	return sb.String()
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
		if strings.ToLower(input) == "quit" {
			return nil, false
		}
		frets, err := parseFrets(input, cfg.InputOrder)
		if err != nil {
			fmt.Printf("Invalid input: %s\nPlease try again.\n\n", err)
			continue
		}
		return frets, true
	}
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
		var progression []Chord

		// --- first chord ---
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

		frets, ok := promptFrets(scanner, cfg)
		if !ok {
			fmt.Println("Goodbye!")
			return
		}
		progression = append(progression, Chord{name: name, frets: frets})

		// --- additional chords ---
	moreChords:
		for {
			fmt.Println()
			answer := prompt(scanner, "Another chord? (chord positions, yes, or no): ")
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
				progression = append(progression, Chord{name: chordName, frets: f})
			default:
				// Treat direct fret input as "yes + chord" with no name.
				f, err := parseFrets(answer, cfg.InputOrder)
				if err != nil {
					fmt.Println("Not understood. Enter chord positions, 'yes', or 'no'.")
					continue
				}
				progression = append(progression, Chord{frets: f})
			}
		}

		// --- render and display ---
		width, height := terminalSize()
		if width > 80 {
			width = 80
		}
		output := renderProgression(progression, width)
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
				fmt.Fprint(f, renderProgression(progression, 80))
				fmt.Fprintln(f)
				f.Close()
				fmt.Printf("Appended to %s\n\n", saveFile)
			}
		} else {
			fmt.Println()
		}
	}
}
