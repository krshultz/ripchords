package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"ripchords/chord"
)

type appState int

const (
	stateFirstRun   appState = iota
	stateChordName           // waiting for chord name
	stateFrets               // waiting for fret positions
	stateRendered            // showing progression, hotkeys active
	stateSave                // waiting for filename
	stateSettings            // settings overlay
	stateLastChord           // showing last chord modal
)

type model struct {
	cfg         Config
	progression []chord.Chord
	state       appState
	input       textinput.Model
	pendingName string
	cursor      int
	lastR       time.Time
	width       int
	height      int
	err         string
	prev        appState // state to return to from settings
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
		m.input.Prompt = "Chord name: "
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
	if m.state == stateChordName || m.state == stateFrets || m.state == stateSave {
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
			m.cfg.InputOrder = chord.PitchOrder
		} else {
			m.cfg.InputOrder = chord.StringOrder
		}
		saveConfig(m.cfg)
		m.state = stateChordName
		m.input.Prompt = "Chord name: "
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
		if len(m.progression) > 0 {
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
			if len(m.progression) > 0 {
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
		frets, err := chord.ParseFrets(val, m.cfg.InputOrder)
		if err != nil {
			m.err = err.Error()
			m.input.SetValue("")
			return m, nil
		}
		m.err = ""
		m.progression = append(m.progression, chord.Chord{Name: m.pendingName, Frets: frets})
		m.state = stateRendered
		m.input.SetValue("")
		return m, nil
	case tea.KeyEsc:
		m.err = ""
		m.state = stateChordName
		m.input.Prompt = "Chord name: "
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
			if len(m.progression) > 0 {
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
		if len(m.progression) > 0 {
			m.prev = m.state
			m.state = stateLastChord
			return m, nil
		}
	case "a":
		m.state = stateChordName
		m.input.Prompt = "Chord name: "
		m.input.SetValue("")
		m.err = ""
		return m, textinput.Blink
	case "s":
		if len(m.progression) == 0 {
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
			os.Remove(configPath()) //nolint
			m.cfg = Config{ShowBarre: true}
			m.progression = nil
			m.state = stateFirstRun
			m.cursor = 0
			m.lastR = time.Time{}
		} else {
			m.lastR = now
			m.progression = nil
		}
	case "?":
		m.prev = m.state
		m.state = stateSettings
		m.cursor = 0
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
			fmt.Fprint(f, chord.RenderProgression(m.progression, 80, m.cfg.ShowBarre))
			fmt.Fprintln(f)
			f.Close()
			m.err = fmt.Sprintf("saved to %s", filename)
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
			if m.cfg.InputOrder == chord.PitchOrder {
				m.cfg.InputOrder = chord.StringOrder
			} else {
				m.cfg.InputOrder = chord.PitchOrder
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

func (m model) View() string {
	switch m.state {
	case stateFirstRun:
		return m.viewFirstRun()
	case stateSettings:
		return m.viewSettings()
	case stateLastChord:
		return m.viewLastChord()
	default:
		return m.viewMain()
	}
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

func (m model) viewSettings() string {
	var b strings.Builder
	b.WriteString(m.viewHeader())
	b.WriteString("\n")

	if len(m.progression) > 0 {
		w := m.width
		if w > 80 {
			w = 80
		}
		b.WriteString(chord.RenderProgression(m.progression, w, m.cfg.ShowBarre))
		b.WriteString("\n")
	}

	b.WriteString("  Settings\n")
	b.WriteString("  " + strings.Repeat("─", 38) + "\n")

	orderVal := "pitch (E A D G B e)"
	if m.cfg.InputOrder == chord.StringOrder {
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
		b.WriteString(fmt.Sprintf("%s%-22s %s\n", prefix, label, values[i]))
	}
	b.WriteString("\n  ↑/↓ select  Space toggle  Enter/Esc close\n")
	return b.String()
}

func (m model) viewLastChord() string {
	var b strings.Builder
	b.WriteString(m.viewHeader())
	b.WriteString("\n")
	last := m.progression[len(m.progression)-1]
	b.WriteString(chord.RenderChord(last.Name, last.Frets, m.cfg.ShowBarre))
	b.WriteString("\n  Esc to dismiss\n")
	return b.String()
}

func (m model) viewHeader() string {
	order := "pitch"
	if m.cfg.InputOrder == chord.StringOrder {
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

	if len(m.progression) > 0 {
		b.WriteString(chord.RenderProgression(m.progression, w, m.cfg.ShowBarre))
		b.WriteString("\n")
	}

	if m.err != "" {
		b.WriteString("  " + m.err + "\n\n")
	}

	switch m.state {
	case stateChordName:
		b.WriteString(m.input.View() + "\n")
	case stateFrets:
		order := "E A D G B e"
		if m.cfg.InputOrder == chord.StringOrder {
			order = "e B G D A E"
		}
		b.WriteString("  (" + order + ")\n")
		b.WriteString(m.input.View() + "\n")
	case stateRendered:
		if len(m.progression) == 0 {
			b.WriteString("  (progression cleared)\n\n")
		}
		b.WriteString("  a add chord  l last chord  s save  r reset  rr wipe config  ? settings  q quit\n")
	case stateSave:
		b.WriteString(m.input.View() + "\n")
	}

	return b.String()
}
