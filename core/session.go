package core

import (
	"fmt"
	"strings"
)

// Session holds an in-progress chord progression and the operations on it.
// It is UI- and IO-free: any front end — the terminal TUI today, a native
// mobile app bound via gomobile later — drives the same progression logic
// through this one type.
type Session struct {
	Chords []Chord
}

// Len reports how many chords are in the progression.
func (s *Session) Len() int { return len(s.Chords) }

// Add parses fretInput in the given order, appends the resulting chord, and
// returns it. The progression is left unchanged if parsing fails.
func (s *Session) Add(name, fretInput string, order InputOrder) (Chord, error) {
	frets, err := ParseFrets(fretInput, order)
	if err != nil {
		return Chord{}, err
	}
	c := Chord{Name: strings.TrimSpace(name), Frets: frets}
	s.Chords = append(s.Chords, c)
	return c, nil
}

// Last returns the most recently added chord. ok is false if the progression
// is empty.
func (s *Session) Last() (chord Chord, ok bool) {
	if len(s.Chords) == 0 {
		return Chord{}, false
	}
	return s.Chords[len(s.Chords)-1], true
}

// Reset clears the progression.
func (s *Session) Reset() { s.Chords = nil }

// Rename changes the name of the chord at index i.
func (s *Session) Rename(i int, name string) error {
	if i < 0 || i >= len(s.Chords) {
		return fmt.Errorf("chord index %d out of range (have %d)", i, len(s.Chords))
	}
	s.Chords[i].Name = strings.TrimSpace(name)
	return nil
}

// EditFrets replaces the fret positions of the chord at index i, parsing
// fretInput in the given order. The chord is left unchanged if parsing fails.
func (s *Session) EditFrets(i int, fretInput string, order InputOrder) error {
	if i < 0 || i >= len(s.Chords) {
		return fmt.Errorf("chord index %d out of range (have %d)", i, len(s.Chords))
	}
	frets, err := ParseFrets(fretInput, order)
	if err != nil {
		return err
	}
	s.Chords[i].Frets = frets
	return nil
}
