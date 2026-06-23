package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"ripchords/core"
)

// feed drives a sequence of key messages through the model's Update loop and
// returns the resulting model.
func feed(m model, keys ...tea.KeyMsg) model {
	for _, k := range keys {
		next, _ := m.Update(k)
		m = next.(model)
	}
	return m
}

func runes(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

func key(t tea.KeyType) tea.KeyMsg { return tea.KeyMsg{Type: t} }

// renderedModel returns a model with two chords already entered, sitting at the
// rendered (hotkeys-active) state.
func renderedModel(t *testing.T) model {
	t.Helper()
	m := newModel(Config{InputOrder: core.PitchOrder, ShowBarre: true})
	if _, err := m.session.Add("C", "x 3 2 0 1 0", core.PitchOrder); err != nil {
		t.Fatalf("seed C: %v", err)
	}
	if _, err := m.session.Add("G", "3 2 0 0 0 3", core.PitchOrder); err != nil {
		t.Fatalf("seed G: %v", err)
	}
	m.state = stateRendered
	return m
}

func TestEditRenameFlow(t *testing.T) {
	m := renderedModel(t)

	// e -> picker, Down to second chord, Enter -> action, Enter (Rename).
	// The field is prefilled with the current name ("G") and the cursor sits at
	// the end, so typing "maj" appends to give "Gmaj". Enter -> back to rendered.
	m = feed(m,
		runes("e"),
		key(tea.KeyDown),
		key(tea.KeyEnter),
		key(tea.KeyEnter),
		runes("maj"),
		key(tea.KeyEnter),
	)

	if m.state != stateRendered {
		t.Errorf("state = %d, want stateRendered (%d)", m.state, stateRendered)
	}
	if got := m.session.Chords[1].Name; got != "Gmaj" {
		t.Errorf("chord[1].Name = %q, want %q", got, "Gmaj")
	}
	if got := m.session.Chords[0].Name; got != "C" {
		t.Errorf("chord[0].Name = %q, want %q (unrelated chord changed)", got, "C")
	}
}

func TestEditFretsFlow(t *testing.T) {
	m := renderedModel(t)

	// e -> picker (first chord selected), Enter -> action, Down to "Edit frets",
	// Enter, type new frets, Enter.
	m = feed(m,
		runes("e"),
		key(tea.KeyEnter),
		key(tea.KeyDown),
		key(tea.KeyEnter),
		runes("x 0 2 2 2 0"), // A
		key(tea.KeyEnter),
	)

	if m.state != stateRendered {
		t.Errorf("state = %d, want stateRendered (%d)", m.state, stateRendered)
	}
	want := []int{-1, 0, 2, 2, 2, 0}
	got := m.session.Chords[0].Frets
	if len(got) != len(want) {
		t.Fatalf("frets len = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("frets = %v, want %v", got, want)
			break
		}
	}
}

func TestEditFretsParseErrorKeepsEditing(t *testing.T) {
	m := renderedModel(t)
	orig := append([]int(nil), m.session.Chords[0].Frets...)

	m = feed(m,
		runes("e"),
		key(tea.KeyEnter),
		key(tea.KeyDown),
		key(tea.KeyEnter),
		runes("nonsense"),
		key(tea.KeyEnter),
	)

	if m.state != stateEditFrets {
		t.Errorf("state = %d, want stateEditFrets (%d) on parse error", m.state, stateEditFrets)
	}
	if m.err == "" {
		t.Error("expected an error message on parse failure, got none")
	}
	for i := range orig {
		if m.session.Chords[0].Frets[i] != orig[i] {
			t.Errorf("frets mutated on parse error: %v, want %v", m.session.Chords[0].Frets, orig)
			break
		}
	}
}

func TestEditPickEscReturnsToRendered(t *testing.T) {
	m := renderedModel(t)
	m = feed(m, runes("e"), key(tea.KeyEsc))
	if m.state != stateRendered {
		t.Errorf("state = %d, want stateRendered (%d) after Esc from picker", m.state, stateRendered)
	}
}
