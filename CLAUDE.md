# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**ripchords** is a Go TUI application that renders guitar chord diagrams in ASCII tab format. The user provides fret positions for each string; the app outputs a standard ASCII chord diagram.

## Commands

```bash
go build ./...          # build
go test ./...           # run all tests
go test ./... -run Foo  # run a single test
go vet ./...            # lint/vet
```

## Domain Knowledge

### Guitar Strings

Standard tuning from lowest to highest pitch: `E A D G B e`

Strings have two numbering conventions that both need to be supported:

| String number | Note | Pitch order |
|---------------|------|-------------|
| 6             | E    | 1st (lowest) |
| 5             | A    | 2nd |
| 4             | D    | 3rd |
| 3             | G    | 4th |
| 2             | B    | 5th |
| 1             | e    | 6th (highest) |

String numbers are "inverted" relative to pitch: string 6 is the lowest pitch. This is a source of confusion — the codebase should be explicit about which convention is in use at any given boundary.

### Input Modes

Users can enter fret positions in two orderings:

- **Pitch order** (low to high): positions for strings 6, 5, 4, 3, 2, 1 in that sequence. C chord: `X 3 2 0 1 0`
- **String-number order** (string 1 to 6): positions for strings 1, 2, 3, 4, 5, 6. C chord: `0 1 0 2 3 X`

### ASCII Tab Format

Chord diagrams render with the highest-pitch string (e) on top, lowest (E) on the bottom. Each string is a horizontal line of hyphens; fret numbers are placed inline. Muted strings are `x`; open strings are `0`. Example (C chord):

```
    C
e |-----0------|
B |-----1------|
G |-----0------|
D |-----2------|
A |-----3------|
E |------------|
```

The low E string is not played in this voicing — it appears as a blank line (or can be marked `x`).
