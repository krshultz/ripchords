package chord

import (
	"fmt"
	"strconv"
	"strings"
)

const numStrings = 6

// InputOrder controls how fret positions are interpreted on input.
type InputOrder string

const (
	PitchOrder  InputOrder = "pitch"
	StringOrder InputOrder = "string_number"
)

// Chord holds a name and fret positions in pitch order (index 0 = low E, index 5 = high e).
type Chord struct {
	Name  string
	Frets []int
}

// stringNames in display order: highest pitch (e) first, lowest (E) last.
var stringNames = [numStrings]string{"e", "B", "G", "D", "A", "E"}

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

// ParseFrets parses a fret-position string and returns values in pitch order
// (index 0 = string 6 = low E, index 5 = string 1 = high e).
// -1 means muted/not played.
func ParseFrets(input string, order InputOrder) ([]int, error) {
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
		for i, j := 0, numStrings-1; i < j; i, j = i+1, j-1 {
			frets[i], frets[j] = frets[j], frets[i]
		}
	}
	return frets, nil
}

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

// RenderChord returns an ASCII tab diagram for a single chord. frets must be in pitch order.
func RenderChord(name string, frets []int) string {
	var sb strings.Builder
	if name != "" {
		sb.WriteString(fmt.Sprintf("    %s\n", name))
	}
	barre := detectBarre(frets)
	for display := 0; display < numStrings; display++ {
		pitchIdx := numStrings - 1 - display
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

func centerInWidth(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	total := w - len(s)
	left := total / 2
	right := total - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// RenderProgression renders chords side-by-side, wrapping rows to fit within width.
// Each segment is 14 chars wide ("|-----0------|") with a 1-space gap between chords.
// Total row width for N chords: 2 (label) + 14*N + (N-1) = 15N+1.
func RenderProgression(chords []Chord, width int) string {
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

		sb.WriteString("  ")
		for i, ch := range row {
			centered := centerInWidth(ch.Name, segWidth)
			sb.WriteString(centered)
			if i < len(row)-1 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")

		for display := 0; display < numStrings; display++ {
			pitchIdx := numStrings - 1 - display
			label := stringNames[display]
			sb.WriteString(label + " ")
			for i, ch := range row {
				fret := ch.Frets[pitchIdx]
				barre := detectBarre(ch.Frets)
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
