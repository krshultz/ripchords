package core

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

func TestBarre(t *testing.T) {
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
			got := Barre(tt.frets)
			if got != tt.wantFret {
				t.Errorf("Barre(%v) = %d, want %d", tt.frets, got, tt.wantFret)
			}
		})
	}
}
