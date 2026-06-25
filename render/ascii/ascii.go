// Package ascii renders core chord models as ASCII tab diagrams for the
// terminal. It is one presentation layer over core; other front ends (e.g. a
// native mobile UI) would draw their own diagrams from the same core model and
// would not import this package.
package ascii

import (
	"fmt"
	"strconv"
	"strings"

	"ripchords/core"
)

// stringNames in display order: highest pitch (e) first, lowest (E) last.
var stringNames = [core.NumStrings]string{"e", "B", "G", "D", "A", "E"}

// barreInfo returns the barre fret (0 if none) and whether it is a mini barre
// (exactly two strings at the barre fret), honoring the showBarre toggle.
func barreInfo(frets []int, showBarre bool) (barre int, mini bool) {
	if !showBarre {
		return 0, false
	}
	barre = core.Barre(frets)
	if barre > 0 {
		n := 0
		for _, f := range frets {
			if f == barre {
				n++
			}
		}
		mini = n == 2
	}
	return barre, mini
}

// stringSegment renders the 14-char tab segment for a single string.
func stringSegment(fret, barre int, mini bool) string {
	var marker string
	switch fret {
	case -1:
		marker = "X"
	case 0:
		marker = "0"
	default:
		marker = strconv.Itoa(fret)
	}
	showMark := (mini && fret == barre) || (!mini && barre > 0 && fret > 0)
	if showMark {
		return "|---|-" + marker + "------|"
	}
	return "|-----" + marker + "------|"
}

// RenderChord returns an ASCII tab diagram for a single chord. frets must be in pitch order.
// showBarre controls whether barre chord markers (|) are rendered.
func RenderChord(name string, frets []int, showBarre bool) string {
	var sb strings.Builder
	if name != "" {
		fmt.Fprintf(&sb, "    %s\n", name)
	}
	barre, mini := barreInfo(frets, showBarre)
	for display := 0; display < core.NumStrings; display++ {
		pitchIdx := core.NumStrings - 1 - display
		label := stringNames[display]
		sb.WriteString(label + " " + stringSegment(frets[pitchIdx], barre, mini) + "\n")
	}
	return sb.String()
}

func centerInWidth(s string, w int) string {
	if len(s) >= w {
		return s[:w]
	}
	total := w - len(s)
	left := total / 2
	right := total - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// RenderProgression renders chords side-by-side, wrapping rows to fit within width.
// Each segment is 14 chars wide ("|-----0------|") with a 1-space gap between chords.
// Total row width for N chords: 2 (label) + 14*N + (N-1) = 15N+1.
// showBarre controls whether barre chord markers (|) are rendered.
func RenderProgression(chords []core.Chord, width int, showBarre bool) string {
	const segWidth = 14 // "|-----0------|"
	const colWidth = 15 // segWidth + 1 space separator

	chordsPerRow := (width - 1) / colWidth
	if chordsPerRow < 1 {
		chordsPerRow = 1
	}

	var sb strings.Builder
	for start := 0; start < len(chords); start += chordsPerRow {
		end := start + chordsPerRow
		if end > len(chords) {
			end = len(chords)
		}
		row := chords[start:end]

		sb.WriteString("  ")
		for i, ch := range row {
			centered := centerInWidth(ch.Name, segWidth)
			sb.WriteString(centered)
			if i < len(row)-1 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")

		barres := make([]int, len(row))
		minis := make([]bool, len(row))
		for i, ch := range row {
			barres[i], minis[i] = barreInfo(ch.Frets, showBarre)
		}

		for display := 0; display < core.NumStrings; display++ {
			pitchIdx := core.NumStrings - 1 - display
			label := stringNames[display]
			sb.WriteString(label + " ")
			for i, ch := range row {
				sb.WriteString(stringSegment(ch.Frets[pitchIdx], barres[i], minis[i]))
				if i < len(row)-1 {
					sb.WriteString(" ")
				}
			}
			sb.WriteString("\n")
		}

		if end < len(chords) {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
