package ascii

import (
	"strings"
	"testing"
)

func TestRenderChordBarreAllFretted(t *testing.T) {
	// F#m: 244222 in pitch order — all strings fretted, barre at 2
	frets := []int{2, 4, 4, 2, 2, 2}
	diagram := RenderChord("F#m", frets, true)
	for _, want := range []string{
		"e |---|-2------|",
		"B |---|-2------|",
		"G |---|-2------|",
		"D |---|-4------|",
		"A |---|-4------|",
		"E |---|-2------|",
	} {
		if !strings.Contains(diagram, want) {
			t.Errorf("barre diagram missing %q, got:\n%s", want, diagram)
		}
	}
}

func TestRenderChordMiniBarreTwoAdjacentStrings(t *testing.T) {
	// FMaj over C: x33211 in pitch order — mini barre at fret 1 (B and e only)
	frets := []int{-1, 3, 3, 2, 1, 1}
	diagram := RenderChord("FMaj", frets, true)
	for _, want := range []string{
		"e |---|-1------|",
		"B |---|-1------|",
		"G |-----2------|",
		"D |-----3------|",
		"A |-----3------|",
		"E |-----X------|",
	} {
		if !strings.Contains(diagram, want) {
			t.Errorf("mini-barre diagram missing %q, got:\n%s", want, diagram)
		}
	}
}

func TestRenderChordBarreMutedString(t *testing.T) {
	// Barre at 5 with muted low E: pitch order [-1, 5, 7, 7, 5, 5]
	frets := []int{-1, 5, 7, 7, 5, 5}
	diagram := RenderChord("", frets, true)
	if !strings.Contains(diagram, "E |-----X------|") {
		t.Errorf("muted string should use standard format (no barre pipe), got:\n%s", diagram)
	}
	if !strings.Contains(diagram, "A |---|-5------|") {
		t.Errorf("fretted A string should show barre pipe, got:\n%s", diagram)
	}
}

func TestRenderChordDWithSomeDoubleDigitFrets(t *testing.T) {
	// Issue 33: C# chord: barre at 9, play all strings
	frets := []int{9, 11, 11, 10, 9, 9}
	diagram := RenderChord("C#", frets, false)
	if !strings.Contains(diagram, "E |-----9------|") {
		t.Errorf("E string should still be aligned, got:\n%s", diagram)
	}
	if !strings.Contains(diagram, "G |-----10-----|") {
		t.Errorf("G string should correct for double digits, got:\n%s", diagram)
	}
	if !strings.Contains(diagram, "D |-----11-----|") {
		t.Errorf("D string should correct for double digits, got:\n%s", diagram)
	}
}

func TestRenderChordDWithAllDoubleDigitFrets(t *testing.T) {
	// Issue 33: E minor chord, full barre
	frets := []int{12, 14, 14, 12, 12, 12}
	diagram := RenderChord("Emin", frets, false)
	if !strings.Contains(diagram, "E |-----12-----|") {
		t.Errorf("E string should be aligned, got:\n%s", diagram)
	}
	if !strings.Contains(diagram, "A |-----14-----|") {
		t.Errorf("A string should be aligned, got:\n%s", diagram)
	}
	if !strings.Contains(diagram, "e |-----12-----|") {
		t.Errorf("e string should be aligned, got:\n%s", diagram)
	}
}

func TestRenderChordMutedAsX(t *testing.T) {
	// C chord: E muted, others have frets
	frets := []int{-1, 3, 2, 0, 1, 0}
	diagram := RenderChord("C", frets, true)
	if !strings.Contains(diagram, "E |-----X------|") {
		t.Errorf("muted E string should render as X, got:\n%s", diagram)
	}
	if !strings.Contains(diagram, "A |-----3------|") {
		t.Errorf("A string should show fret 3, got:\n%s", diagram)
	}
	if !strings.Contains(diagram, "e |-----0------|") {
		t.Errorf("e string should show 0, got:\n%s", diagram)
	}
}

func TestRenderChordName(t *testing.T) {
	frets := []int{0, 0, 0, 0, 0, 0}
	diagram := RenderChord("Em", frets, true)
	if !strings.HasPrefix(diagram, "    Em\n") {
		t.Errorf("diagram should start with chord name, got:\n%s", diagram)
	}
}

func TestRenderChordNoName(t *testing.T) {
	frets := []int{0, 0, 0, 0, 0, 0}
	diagram := RenderChord("", frets, true)
	if strings.HasPrefix(diagram, " ") {
		t.Errorf("nameless diagram should not start with spaces, got:\n%s", diagram)
	}
	if !strings.HasPrefix(diagram, "e") {
		t.Errorf("nameless diagram should start with 'e', got:\n%s", diagram)
	}
}
