// Package core holds the pure, UI- and IO-free domain logic for ripchords:
// the chord model, fret-position parsing, and music-theory helpers. It depends
// only on the standard library so it can be reused unchanged by any front end
// (the terminal TUI today, a gomobile-bound native app later).
package core

import (
	"fmt"
	"strconv"
	"strings"
)

// NumStrings is the number of strings on a standard guitar.
const NumStrings = 6

// Chord holds a name and fret positions in pitch order (index 0 = low E, index 5 = high e).
type Chord struct {
	Name  string
	Frets []int
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

// ParseFrets parses a fret-position string given in pitch order, low to high
// (E A D G B e), and returns the values in that same order
// (index 0 = string 6 = low E, index 5 = string 1 = high e).
// -1 means muted/not played.
func ParseFrets(input string) ([]int, error) {
	tokens := tokenize(input)
	if len(tokens) != NumStrings {
		if len(tokens) > NumStrings {
			return nil, fmt.Errorf(
				"%d strings is too many. Frets above 10 use spaces (e.g., \"9 11 11 19 9 9\") so they process as %d",
				len(tokens), NumStrings,
			)
		}
		return nil, fmt.Errorf(
			"expected %d string positions, got %d\n  (e.g. \"x 3 2 0 1 0\" or \"x32010\" for C major)",
			NumStrings, len(tokens),
		)
	}
	frets := make([]int, NumStrings)
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
	return frets, nil
}

// Barre returns the fret at which a barre is held, or 0 if the chord has no barre.
// A barre is detected when two or more adjacent strings share the lowest fretted fret.
func Barre(frets []int) int {
	minFret := 0
	for _, f := range frets {
		if f > 0 && (minFret == 0 || f < minFret) {
			minFret = f
		}
	}
	if minFret == 0 {
		return 0
	}
	prev := false
	for _, f := range frets {
		at := f == minFret
		if at && prev {
			return minFret
		}
		prev = at
	}
	return 0
}
