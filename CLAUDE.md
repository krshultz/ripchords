# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**ripchords** is a Go TUI application that renders guitar chord diagrams in ASCII tab format. The user provides fret positions for each string; the app outputs a standard ASCII chord diagram.

## Project Status / Backlog

Open work and roadmap live in **`.specs/backlog.md`** (specs in general live under `.specs/`).
Check it at the start of a session to see what's outstanding. When an item ships (merges to
main) it moves to the "Recently shipped" section at the bottom of the backlog — not deleted —
so what remains above is open work or deviations from the original spec.

## Keep Documentation in Sync (required)

Every change to ripchords updates the docs it affects **in the same PR**, not a follow-up. A
PR isn't done until these match the new behavior:

- **`.specs/backlog.md`**: move shipped items to "Recently shipped" and fix any open item whose
  wording the change makes stale (e.g. "once #48 lands").
- **`README.md`**: usage, hotkeys table, flags, examples and rendered diagrams.
- **`CLAUDE.md`** (this file): domain notes, examples, and anything that describes how the code
  works.
- **`CONTRIBUTING.md`**: if tooling, branch rules or commit conventions change.
- **In-app text**: the `-h` usage text and help screen come from `keys.go`; prompts, header
  and modal footers in `tui.go` are hand-written and need checking by hand.

Before committing, grep the docs for terms the change touches (old key names, removed options,
renamed screens) and verify each claim against the code or the built binary. Don't edit
`CHANGELOG.md`; release-please writes it.

## Commands

```bash
go build ./...          # build
go test ./...           # run all tests
go test ./... -run Foo  # run a single test
go vet ./...            # lint/vet
```

## Releases & Versioning

Versioning is automated with [release-please](https://github.com/googleapis/release-please).
The version is **derived from commit messages**, so commits/PR titles must follow
[Conventional Commits](https://www.conventionalcommits.org/):

| Prefix              | Effect (pre-1.0)            | Example                          |
|---------------------|-----------------------------|----------------------------------|
| `fix:`              | patch bump (0.1.0 → 0.1.1)  | `fix: correct mini-barre render` |
| `feat:`             | minor bump (0.1.0 → 0.2.0)  | `feat: add capo support`         |
| `feat!:` / `fix!:`  | minor bump while < 1.0.0¹   | `feat!: change settings format`  |
| `chore:`/`docs:`/`refactor:`/`test:` | no release | `chore: tidy backlog`            |

¹ `bump-minor-pre-major` is on, so breaking changes bump the minor (not major) until 1.0.0.

**How a release happens:** every push to `main` runs the release-please workflow, which
maintains an open "release PR" (titled `chore: release X.Y.Z`) accumulating changes and a
CHANGELOG. Merging that PR creates the git tag (`vX.Y.Z`) and the GitHub Release. The binary's
version is stamped from the tag via ldflags (see version resolution in `main.go`), so the tag,
`--version` output, and GitHub Release stay in sync.

Config lives in `release-please-config.json`; current version is tracked in
`.release-please-manifest.json`. Squash-merge PRs so the PR title becomes the conventional commit.

## Domain Knowledge

### Guitar Strings

Standard tuning from lowest to highest pitch: `E A D G B e`

Strings have two numbering conventions:

| String number | Note | Pitch order |
|---------------|------|-------------|
| 6             | E    | 1st (lowest) |
| 5             | A    | 2nd |
| 4             | D    | 3rd |
| 3             | G    | 4th |
| 2             | B    | 5th |
| 1             | e    | 6th (highest) |

String numbers are "inverted" relative to pitch: string 6 is the lowest pitch. This is a source of confusion — the codebase should be explicit about which convention is in use at any given boundary.

### Input Order

Fret positions are entered in **pitch order** only (low to high): positions for strings 6, 5, 4, 3, 2, 1 in that sequence. C chord: `X 3 2 0 1 0`. `core.ParseFrets` returns them in the same order (index 0 = low E).

String-number order (string 1 to 6) was supported until #48 and has been removed. Older `~/.ripchords/config.json` files may still contain an `input_order` field; it is ignored on load and dropped on the next save.

### ASCII Tab Format

Chord diagrams render with the highest-pitch string (e) on top, lowest (E) on the bottom. Each string is a horizontal line of hyphens; fret numbers are placed inline. Muted strings are entered as `x` (either case) and render as `X`; open strings are `0`. Example (C chord):

```
    C
e |-----0------|
B |-----1------|
G |-----0------|
D |-----2------|
A |-----3------|
E |-----X------|
```

The low E string is not played in this voicing, so it renders as `X`.
