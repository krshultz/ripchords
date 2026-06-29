package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/krshultz/ripchords/core"
	"github.com/krshultz/ripchords/render/ascii"
)

type appState int

const (
	stateFirstRun      appState = iota // first-run: choosing input order
	stateFirstRunBarre                 // first-run: choosing barre on/off
	stateChordName                     // waiting for chord name
	stateFrets                         // waiting for fret positions
	stateRendered                      // showing progression, hotkeys active
	stateSave                          // waiting for filename
	stateSettings                      // settings overlay
	stateLastChord                     // showing last chord modal
	stateConfirmWipe                   // confirming a full config wipe
	stateEditPick                      // choosing which chord to edit
	stateEditAction                    // choosing rename vs edit frets
	stateEditName                      // entering a new name for the chosen chord
	stateEditFrets                     // entering new fret positions for the chosen chord
)

type model struct {
	cfg         Config
	session     core.Session
	state       appState
	input       textinput.Model
	pendingName string
	cursor      int
	lastR       time.Time
	width       int
	height      int
	err         string
	prev        appState // state to return to from settings
	editIdx     int      // index of the chord being edited
}

func chordNamePrompt(n int) string {
	if n == 0 {
		return "Chord name: "
	}
	return "Next chord name: "
}

func newModel(cfg Config) model {
	ti := textinput.New()
	ti.Focus()
	m := model{
		cfg:    cfg,
		input:  ti,
		width:  80,
		height: 24,
	}
	if cfg.InputOrder == "" {
		m.state = stateFirstRun
	} else {
		m.state = stateChordName
		m.input.Prompt = chordNamePrompt(m.session.Len())
	}
	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	switch m.state {
	case stateChordName, stateFrets, stateSave, stateEditName, stateEditFrets:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.state {
	case stateFirstRun:
		return m.handleFirstRun(msg)
	case stateFirstRunBarre:
		return m.handleFirstRunBarre(msg)
	case stateConfirmWipe:
		return m.handleConfirmWipe(msg)
	case stateChordName:
		return m.handleChordName(msg)
	case stateFrets:
		return m.handleFrets(msg)
	case stateRendered:
		return m.handleRendered(msg)
	case stateSave:
		return m.handleSave(msg)
	case stateSettings:
		return m.handleSettings(msg)
	case stateLastChord:
		return m.handleLastChord(msg)
	case stateEditPick:
		return m.handleEditPick(msg)
	case stateEditAction:
		return m.handleEditAction(msg)
	case stateEditName:
		return m.handleEditName(msg)
	case stateEditFrets:
		return m.handleEditFrets(msg)
	}
	return m, nil
}

func (m model) handleFirstRun(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < 1 {
			m.cursor++
		}
	case tea.KeyEnter:
		if m.cursor == 0 {
			m.cfg.InputOrder = core.PitchOrder
		} else {
			m.cfg.InputOrder = core.StringOrder
		}
		// Advance to the barre choice; config is saved once both are picked.
		m.state = stateFirstRunBarre
		m.cursor = 0
		if !m.cfg.ShowBarre {
			m.cursor = 1
		}
		return m, nil
	}
	return m, nil
}

func (m model) handleFirstRunBarre(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < 1 {
			m.cursor++
		}
	case tea.KeyEnter:
		m.cfg.ShowBarre = m.cursor == 0
		saveConfig(m.cfg)
		m.state = stateChordName
		m.input.Prompt = chordNamePrompt(m.session.Len())
		m.input.SetValue("")
		m.cursor = 0
		return m, textinput.Blink
	}
	return m, nil
}

func (m model) handleChordName(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		m.pendingName = strings.TrimSpace(m.input.Value())
		m.err = ""
		m.state = stateFrets
		m.input.Prompt = "Fret positions: "
		m.input.SetValue("")
		return m, textinput.Blink
	case tea.KeyEsc:
		if m.session.Len() > 0 {
			m.state = stateRendered
			m.input.SetValue("")
			return m, nil
		}
	}
	if m.input.Value() == "" {
		switch msg.String() {
		case "?":
			m.prev = m.state
			m.state = stateSettings
			m.cursor = 0
			return m, nil
		case "l":
			if m.session.Len() > 0 {
				m.prev = m.state
				m.state = stateLastChord
				return m, nil
			}
		case "q":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) handleFrets(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		val := strings.TrimSpace(m.input.Value())
		if _, err := m.session.Add(m.pendingName, val, m.cfg.InputOrder); err != nil {
			m.err = err.Error()
			m.input.SetValue("")
			return m, nil
		}
		m.err = ""
		m.state = stateRendered
		m.input.SetValue("")
		return m, nil
	case tea.KeyEsc:
		m.err = ""
		m.state = stateChordName
		m.input.Prompt = chordNamePrompt(m.session.Len())
		m.input.SetValue(m.pendingName)
		return m, textinput.Blink
	}
	if m.input.Value() == "" {
		switch msg.String() {
		case "?":
			m.prev = m.state
			m.state = stateSettings
			m.cursor = 0
			return m, nil
		case "l":
			if m.session.Len() > 0 {
				m.prev = m.state
				m.state = stateLastChord
				return m, nil
			}
		case "q":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) handleRendered(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "l":
		if m.session.Len() > 0 {
			m.prev = m.state
			m.state = stateLastChord
			return m, nil
		}
	case "a":
		m.state = stateChordName
		m.input.Prompt = chordNamePrompt(m.session.Len())
		m.input.SetValue("")
		m.err = ""
		return m, textinput.Blink
	case "e":
		if m.session.Len() > 0 {
			m.state = stateEditPick
			m.cursor = 0
			m.err = ""
		}
		return m, nil
	case "s":
		if m.session.Len() == 0 {
			return m, nil
		}
		m.state = stateSave
		m.input.Prompt = "Save to file: "
		m.input.SetValue("")
		m.err = ""
		return m, textinput.Blink
	case "r":
		now := time.Now()
		if !m.lastR.IsZero() && now.Sub(m.lastR) <= 500*time.Millisecond {
			// Second quick r: don't wipe yet — ask for confirmation first, so
			// an accidental double-tap can't destroy the saved config.
			m.lastR = time.Time{}
			m.state = stateConfirmWipe
		} else {
			m.lastR = now
			m.session.Reset()
		}
	case "?":
		m.prev = m.state
		m.state = stateSettings
		m.cursor = 0
	}
	return m, nil
}

func (m model) handleConfirmWipe(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}
	switch msg.String() {
	case "y", "Y":
		os.Remove(configPath()) //nolint
		m.cfg = Config{ShowBarre: true}
		m.session.Reset()
		m.state = stateFirstRun
		m.cursor = 0
	default: // n, Esc, or anything else cancels
		m.state = stateRendered
	}
	return m, nil
}

func (m model) handleSave(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		filename := expandPath(strings.TrimSpace(m.input.Value()))
		if filename != "" {
			f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				m.err = fmt.Sprintf("could not open file: %s", err)
				m.input.SetValue("")
				return m, nil
			}
			_, werr := fmt.Fprint(f, ascii.RenderProgression(m.session.Chords, 80, m.cfg.ShowBarre))
			if werr == nil {
				_, werr = fmt.Fprintln(f)
			}
			if cerr := f.Close(); werr == nil {
				werr = cerr
			}
			if werr != nil {
				m.err = fmt.Sprintf("could not write file: %s", werr)
			} else {
				m.err = fmt.Sprintf("saved to %s", filename)
			}
		}
		m.state = stateRendered
		m.input.SetValue("")
		return m, nil
	case tea.KeyEsc:
		m.err = ""
		m.state = stateRendered
		m.input.SetValue("")
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

var settingsLabels = []string{"Input order", "Show barre chords"}

func (m model) handleSettings(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc, tea.KeyEnter:
		saveConfig(m.cfg)
		m.state = m.prev
		if m.prev == stateChordName || m.prev == stateFrets || m.prev == stateSave {
			return m, textinput.Blink
		}
		return m, nil
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < len(settingsLabels)-1 {
			m.cursor++
		}
	case tea.KeySpace:
		switch m.cursor {
		case 0:
			if m.cfg.InputOrder == core.PitchOrder {
				m.cfg.InputOrder = core.StringOrder
			} else {
				m.cfg.InputOrder = core.PitchOrder
			}
		case 1:
			m.cfg.ShowBarre = !m.cfg.ShowBarre
		}
	}
	return m, nil
}

func (m model) handleLastChord(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc, tea.KeyEnter:
		m.state = m.prev
		if m.prev == stateChordName || m.prev == stateFrets || m.prev == stateSave {
			return m, textinput.Blink
		}
		return m, nil
	}
	return m, nil
}

func (m model) handleEditPick(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.state = stateRendered
		return m, nil
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < m.session.Len()-1 {
			m.cursor++
		}
	case tea.KeyEnter:
		m.editIdx = m.cursor
		m.state = stateEditAction
		m.cursor = 0
		return m, nil
	}
	return m, nil
}

var editActionLabels = []string{"Rename", "Edit frets"}

func (m model) handleEditAction(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.state = stateEditPick
		m.cursor = m.editIdx
		return m, nil
	case tea.KeyUp:
		if m.cursor > 0 {
			m.cursor--
		}
	case tea.KeyDown:
		if m.cursor < len(editActionLabels)-1 {
			m.cursor++
		}
	case tea.KeyEnter:
		if m.cursor == 0 {
			m.state = stateEditName
			m.input.Prompt = "New name: "
			m.input.SetValue(m.session.Chords[m.editIdx].Name)
			m.input.CursorEnd()
		} else {
			m.state = stateEditFrets
			m.input.Prompt = "New fret positions: "
			m.input.SetValue("")
		}
		return m, textinput.Blink
	}
	return m, nil
}

func (m model) handleEditName(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		if err := m.session.Rename(m.editIdx, m.input.Value()); err != nil {
			m.err = err.Error()
		} else {
			m.err = ""
		}
		m.state = stateRendered
		m.input.SetValue("")
		return m, nil
	case tea.KeyEsc:
		m.state = stateEditAction
		m.cursor = 0
		m.input.SetValue("")
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) handleEditFrets(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEnter:
		val := strings.TrimSpace(m.input.Value())
		if err := m.session.EditFrets(m.editIdx, val, m.cfg.InputOrder); err != nil {
			m.err = err.Error()
			m.input.SetValue("")
			return m, nil
		}
		m.err = ""
		m.state = stateRendered
		m.input.SetValue("")
		return m, nil
	case tea.KeyEsc:
		m.err = ""
		m.state = stateEditAction
		m.cursor = 0
		m.input.SetValue("")
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	switch m.state {
	case stateFirstRun:
		return m.viewFirstRun()
	case stateFirstRunBarre:
		return m.viewFirstRunBarre()
	case stateConfirmWipe:
		return m.viewConfirmWipe()
	case stateSettings:
		return m.viewSettings()
	case stateLastChord:
		return m.viewLastChord()
	case stateEditPick:
		return m.viewEditPick()
	case stateEditAction:
		return m.viewEditAction()
	default:
		return m.viewMain()
	}
}

func (m model) chordLabel(i int) string {
	name := m.session.Chords[i].Name
	if name == "" {
		name = "(unnamed)"
	}
	return name
}

func (m model) viewEditPick() string {
	var b strings.Builder
	b.WriteString(m.viewHeader())
	b.WriteString("\n")

	w := m.width
	if w > 80 {
		w = 80
	}
	b.WriteString(ascii.RenderProgression(m.session.Chords, w, m.cfg.ShowBarre))
	b.WriteString("\n")

	b.WriteString("  Edit which chord?\n")
	for i := range m.session.Chords {
		prefix := "    "
		if i == m.cursor {
			prefix = "  > "
		}
		fmt.Fprintf(&b, "%s%d. %s\n", prefix, i+1, m.chordLabel(i))
	}
	b.WriteString("\n  ↑/↓ select  Enter choose  Esc cancel\n")
	return b.String()
}

func (m model) viewEditAction() string {
	var b strings.Builder
	b.WriteString(m.viewHeader())
	b.WriteString("\n")

	last := m.session.Chords[m.editIdx]
	b.WriteString(ascii.RenderChord(last.Name, last.Frets, m.cfg.ShowBarre))
	b.WriteString("\n")

	fmt.Fprintf(&b, "  Editing chord %d: %s\n", m.editIdx+1, m.chordLabel(m.editIdx))
	for i, label := range editActionLabels {
		prefix := "    "
		if i == m.cursor {
			prefix = "  > "
		}
		b.WriteString(prefix + label + "\n")
	}
	b.WriteString("\n  ↑/↓ select  Enter choose  Esc back\n")
	return b.String()
}

func (m model) viewFirstRun() string {
	var b strings.Builder
	b.WriteString("Welcome to ripchords!\n\n")
	b.WriteString("How do you prefer to enter fret positions?\n\n")
	opts := []string{
		`  Pitch order (low to high): E A D G B e  —  e.g. "x 3 2 0 1 0" for C`,
		`  String-number order (string 1–6): e B G D A E  —  e.g. "0 1 0 2 3 x" for C`,
	}
	for i, opt := range opts {
		if i == m.cursor {
			b.WriteString("> " + opt[2:] + "\n")
		} else {
			b.WriteString(opt + "\n")
		}
	}
	b.WriteString("\n  ↑/↓ select  Enter confirm  Ctrl+C quit\n")
	return b.String()
}

func (m model) viewFirstRunBarre() string {
	var b strings.Builder
	b.WriteString("Welcome to ripchords!\n\n")
	b.WriteString("Show barre chords as a single barred line where possible?\n\n")
	opts := []string{
		"  Yes — render barres (recommended)",
		"  No — show every fret individually",
	}
	for i, opt := range opts {
		if i == m.cursor {
			b.WriteString("> " + opt[2:] + "\n")
		} else {
			b.WriteString(opt + "\n")
		}
	}
	b.WriteString("\n  ↑/↓ select  Enter confirm  Ctrl+C quit\n")
	return b.String()
}

func (m model) viewConfirmWipe() string {
	var b strings.Builder
	b.WriteString(m.viewHeader())
	b.WriteString("\n")
	b.WriteString("  Wipe all saved settings and restart setup?\n")
	b.WriteString("  This clears your input-order and barre preferences.\n")
	b.WriteString("\n  y wipe config   n / Esc cancel\n")
	return b.String()
}

func (m model) viewSettings() string {
	var b strings.Builder
	b.WriteString(m.viewHeader())
	b.WriteString("\n")

	if m.session.Len() > 0 {
		w := m.width
		if w > 80 {
			w = 80
		}
		b.WriteString(ascii.RenderProgression(m.session.Chords, w, m.cfg.ShowBarre))
		b.WriteString("\n")
	}

	b.WriteString("  Settings\n")
	b.WriteString("  " + strings.Repeat("─", 38) + "\n")

	orderVal := "pitch (E A D G B e)"
	if m.cfg.InputOrder == core.StringOrder {
		orderVal = "string (e B G D A E)"
	}
	barreVal := "off"
	if m.cfg.ShowBarre {
		barreVal = "on"
	}
	values := []string{orderVal, barreVal}

	for i, label := range settingsLabels {
		prefix := "    "
		if i == m.cursor {
			prefix = "  > "
		}
		fmt.Fprintf(&b, "%s%-22s %s\n", prefix, label, values[i])
	}
	b.WriteString("\n  ↑/↓ select  Space toggle  Enter/Esc close\n")
	return b.String()
}

func (m model) viewLastChord() string {
	var b strings.Builder
	b.WriteString(m.viewHeader())
	b.WriteString("\n")
	last, _ := m.session.Last()
	b.WriteString(ascii.RenderChord(last.Name, last.Frets, m.cfg.ShowBarre))
	b.WriteString("\n  Esc to dismiss\n")
	return b.String()
}

func (m model) viewHeader() string {
	order := "pitch"
	if m.cfg.InputOrder == core.StringOrder {
		order = "string"
	}
	return fmt.Sprintf("ripchords — %s order  |  l last chord  |  ? settings  |  r reset  |  q quit\n", order)
}

func (m model) viewMain() string {
	var b strings.Builder
	b.WriteString(m.viewHeader())
	b.WriteString("\n")

	w := m.width
	if w > 80 {
		w = 80
	}

	if m.session.Len() > 0 {
		b.WriteString(ascii.RenderProgression(m.session.Chords, w, m.cfg.ShowBarre))
		b.WriteString("\n")
	}

	if m.err != "" {
		b.WriteString("  " + m.err + "\n\n")
	}

	switch m.state {
	case stateChordName:
		fmt.Fprintf(&b, "  Num chords: %d\n", m.session.Len())
		b.WriteString(m.input.View() + "\n")
	case stateFrets:
		fmt.Fprintf(&b, "  Num chords: %d\n", m.session.Len())
		order := "E A D G B e"
		if m.cfg.InputOrder == core.StringOrder {
			order = "e B G D A E"
		}
		b.WriteString("  (" + order + ")\n")
		b.WriteString(m.input.View() + "\n")
	case stateRendered:
		if m.session.Len() == 0 {
			b.WriteString("  (progression cleared)\n\n")
		}
		b.WriteString("  a add chord  e edit chord  l last chord  s save  r reset  rr wipe config  ? settings  q quit\n")
	case stateSave:
		b.WriteString(m.input.View() + "\n")
	case stateEditName:
		fmt.Fprintf(&b, "  Editing chord %d name\n", m.editIdx+1)
		b.WriteString(m.input.View() + "\n")
	case stateEditFrets:
		fmt.Fprintf(&b, "  Editing chord %d frets\n", m.editIdx+1)
		order := "E A D G B e"
		if m.cfg.InputOrder == core.StringOrder {
			order = "e B G D A E"
		}
		b.WriteString("  (" + order + ")\n")
		b.WriteString(m.input.View() + "\n")
	}

	return b.String()
}
