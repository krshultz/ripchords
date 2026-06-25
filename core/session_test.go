package core

import "testing"

// mustAdd adds a chord to s in pitch order and fails the test if it errors.
// Used for known-good setup where a parse failure means the test itself is broken.
func mustAdd(t *testing.T, s *Session, name, frets string) {
	t.Helper()
	if _, err := s.Add(name, frets, PitchOrder); err != nil {
		t.Fatalf("setup Add(%q, %q): %v", name, frets, err)
	}
}

func TestSessionAddAndLast(t *testing.T) {
	var s Session
	if s.Len() != 0 {
		t.Fatalf("new session Len = %d, want 0", s.Len())
	}
	if _, ok := s.Last(); ok {
		t.Errorf("Last() on empty session ok = true, want false")
	}

	c, err := s.Add("  C  ", "x 3 2 0 1 0", PitchOrder)
	if err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	if c.Name != "C" {
		t.Errorf("Add chord name = %q, want %q (trimmed)", c.Name, "C")
	}
	if s.Len() != 1 {
		t.Errorf("Len after one Add = %d, want 1", s.Len())
	}
	last, ok := s.Last()
	if !ok || last.Name != "C" {
		t.Errorf("Last() = (%v, %v), want C chord", last, ok)
	}
}

func TestSessionAddParseErrorLeavesProgressionUnchanged(t *testing.T) {
	var s Session
	if _, err := s.Add("bad", "1 2 3", PitchOrder); err == nil {
		t.Fatal("Add with too few positions: expected error, got nil")
	}
	if s.Len() != 0 {
		t.Errorf("Len after failed Add = %d, want 0", s.Len())
	}
}

func TestSessionReset(t *testing.T) {
	var s Session
	mustAdd(t, &s, "C", "x32010")
	mustAdd(t, &s, "G", "320003")
	s.Reset()
	if s.Len() != 0 {
		t.Errorf("Len after Reset = %d, want 0", s.Len())
	}
}

func TestSessionRename(t *testing.T) {
	var s Session
	mustAdd(t, &s, "C", "x32010")
	if err := s.Rename(0, "  Cmaj  "); err != nil {
		t.Fatalf("Rename returned error: %v", err)
	}
	if got := s.Chords[0].Name; got != "Cmaj" {
		t.Errorf("after Rename name = %q, want %q (trimmed)", got, "Cmaj")
	}
	if err := s.Rename(5, "x"); err == nil {
		t.Error("Rename out of range: expected error, got nil")
	}
}

func TestSessionEditFrets(t *testing.T) {
	var s Session
	mustAdd(t, &s, "C", "x32010")

	if err := s.EditFrets(0, "x 3 2 0 1 3", PitchOrder); err != nil {
		t.Fatalf("EditFrets returned error: %v", err)
	}
	want := []int{-1, 3, 2, 0, 1, 3}
	for i, w := range want {
		if s.Chords[0].Frets[i] != w {
			t.Errorf("after EditFrets frets[%d] = %d, want %d", i, s.Chords[0].Frets[i], w)
		}
	}

	// Bad input leaves the chord unchanged.
	if err := s.EditFrets(0, "nope", PitchOrder); err == nil {
		t.Error("EditFrets with bad input: expected error, got nil")
	}
	if s.Chords[0].Frets[5] != 3 {
		t.Errorf("failed EditFrets mutated chord: frets[5] = %d, want 3", s.Chords[0].Frets[5])
	}

	if err := s.EditFrets(-1, "x32010", PitchOrder); err == nil {
		t.Error("EditFrets out of range: expected error, got nil")
	}
}
