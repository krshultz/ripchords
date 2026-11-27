package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

// keyMap is the single source of truth for ripchords' editor hotkeys. The
// help bar, the ?/h help screen, and the -h usage text are all generated from
// it, so the help can't drift from the keys the handlers actually match.
type keyMap struct {
	Add     key.Binding
	Edit    key.Binding
	Last    key.Binding
	Save    key.Binding
	Reset   key.Binding
	Wipe    key.Binding // help text only: "rr" is two presses of Reset, detected by timing
	Options key.Binding
	Help    key.Binding
	Quit    key.Binding
}

var keys = keyMap{
	Add:     key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add chord")),
	Edit:    key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit chord")),
	Last:    key.NewBinding(key.WithKeys("l"), key.WithHelp("l", "last chord")),
	Save:    key.NewBinding(key.WithKeys("s"), key.WithHelp("s", "save")),
	Reset:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "reset progression")),
	Wipe:    key.NewBinding(key.WithKeys("rr"), key.WithHelp("rr", "reset options")),
	Options: key.NewBinding(key.WithKeys("o"), key.WithHelp("o", "options")),
	Help:    key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?/h", "help")),
	Quit:    key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
}

// ShortHelp is the one-line help bar on the progression screen. It uses
// shorter descriptions than FullHelp so the bar fits in 80 columns (help.Model
// overflows rather than truncating when even the "…" won't fit).
func (k keyMap) ShortHelp() []key.Binding {
	short := func(b key.Binding, desc string) key.Binding {
		b.SetHelp(b.Help().Key, desc)
		return b
	}
	return []key.Binding{
		short(k.Add, "add"),
		short(k.Edit, "edit"),
		k.Last,
		k.Save,
		short(k.Reset, "reset"),
		k.Options,
		k.Help,
		k.Quit,
	}
}

// FullHelp is the ?/h help screen, one inner slice per column.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Edit, k.Last, k.Save},
		{k.Reset, k.Wipe, k.Options, k.Help, k.Quit},
	}
}

// newHelp returns a help view styled to match the rest of the app: bold keys,
// plain descriptions. (The bubbles defaults are dim grays that are hard to
// read on some terminals.)
func newHelp() help.Model {
	h := help.New()
	bold := lipgloss.NewStyle().Bold(true)
	plain := lipgloss.NewStyle()
	h.Styles = help.Styles{
		Ellipsis:       plain,
		ShortKey:       bold,
		ShortDesc:      plain,
		ShortSeparator: plain,
		FullKey:        bold,
		FullDesc:       plain,
		FullSeparator:  plain,
	}
	h.ShortSeparator = "  "
	return h
}

// hotkeyUsage renders every binding from FullHelp as plain text for the -h
// usage screen, one per line.
func hotkeyUsage() string {
	var b strings.Builder
	for _, col := range keys.FullHelp() {
		for _, kb := range col {
			fmt.Fprintf(&b, "  %-4s %s\n", kb.Help().Key, kb.Help().Desc)
		}
	}
	return b.String()
}
