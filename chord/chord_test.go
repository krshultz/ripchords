package chord

import (
	"strings"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"x32010", []string{"x", "3", "2", "0", "1", "0"}},
		{"x 3 2 0 1 0", []string{"x", "3", "2", "0", "1", "0"}},
		{"X 3 2 0 1 0", []string{"X", "3", "2", "0", "1", "0"}},
		{"0 1 0 2 3 x", []string{"0", "1", "0", "2", "3", "x"}},
	}
	for _, tt := range tests {
		got := tokenize(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("tokenize(%q) len = %d, want %d", tt.input, len(got), len(tt.want))
			continue
		}
		for i := range tt.want {
			if got[i] != tt.want[i] {
				t.Errorf("tokenize(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestParseFretsPitchOrder(t *testing.T) {
	// C chord pitch order: x 3 2 0 1 0 → E=muted, A=3, D=2, G=0, B=1, e=0
	frets, err := ParseFrets("x 3 2 0 1 0", PitchOrder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{-1, 3, 2, 0, 1, 0}
	for i, w := range want {
		if frets[i] != w {
			t.Errorf("frets[%d] = %d, want %d", i, frets[i], w)
		}
	}
}

func TestParseFretsPitchOrderSpaceless(t *testing.T) {
	// Same as above but spaceless
	frets, err := ParseFrets("x32010", PitchOrder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{-1, 3, 2, 0, 1, 0}
	for i, w := range want {
		if frets[i] != w {
			t.Errorf("frets[%d] = %d, want %d", i, frets[i], w)
		}
	}
}

func TestParseFretsStringOrder(t *testing.T) {
	// C chord string-number order: 0 1 0 2 3 x → string1(e)=0, 2(B)=1, 3(G)=0, 4(D)=2, 5(A)=3, 6(E)=x
	// In pitch order: [E=-1, A=3, D=2, G=0, B=1, e=0]
	frets, err := ParseFrets("0 1 0 2 3 x", StringOrder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []int{-1, 3, 2, 0, 1, 0}
	for i, w := range want {
		if frets[i] != w {
			t.Errorf("frets[%d] = %d, want %d", i, frets[i], w)
		}
	}
}

func TestParseFretsErrors(t *testing.T) {
	tests := []struct {
		input     string
		wantInErr string
	}{
		{"1 2 3", "expected 6"},
		{"1 2 3 4 5 6 7", "expected 6"},
		{"1 2 z 4 5 6", "not a valid fret"},
		{"1 2 25 4 5 6", "out of range"},
	}
	for _, tt := range tests {
		_, err := ParseFrets(tt.input, PitchOrder)
		if err == nil {
			t.Errorf("ParseFrets(%q) expected error, got nil", tt.input)
			continue
		}
		if !strings.Contains(err.Error(), tt.wantInErr) {
			t.Errorf("ParseFrets(%q) error = %q, want it to contain %q", tt.input, err.Error(), tt.wantInErr)
		}
	}
}

func TestDetectBarre(t *testing.T) {
	tests := []struct {
		name     string
		frets    []int // pitch order: [E A D G B e]
		wantFret int
	}{
		{"F#m barre at 2", []int{2, 4, 4, 2, 2, 2}, 2},
		{"open C no barre", []int{-1, 3, 2, 0, 1, 0}, 0},
		{"all open no barre", []int{0, 0, 0, 0, 0, 0}, 0},
		{"three strings at min fret", []int{-1, 5, 7, 7, 5, 5}, 5},
		{"two non-adjacent strings at min fret no barre", []int{-1, 7, 9, 9, 9, 7}, 0},
		{"two adjacent strings at min fret mini barre", []int{-1, 3, 3, 2, 1, 1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectBarre(tt.frets)
			if got != tt.wantFret {
				t.Errorf("detectBarre(%v) = %d, want %d", tt.frets, got, tt.wantFret)
			}
		})
	}
}

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
