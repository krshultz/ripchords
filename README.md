# ripchords

> _Software for guitar players who want to write down what they just played_

---

## What it does

Hi, I'm Karl. I play guitar and I love to improvise. Every now and then I stumble
across a chord I really like and want to write it down — but that's harder than it
sounds. I don't read music well, and I struggle to remember the "right" way to draw
a chord. Diagram? Standard notation? Tablature? Which way up? By the time I figure
it out, I've forgotten the chord.

What I *can* do easily is describe a chord as six numbers, one per string — like
`320003` for a G major. So ripchords does the boring part: you type the fret
positions, it draws the chord as [ASCII tab](https://en.wikipedia.org/wiki/ASCII_tab).

Ripchords was designed by me, to solve a problem for me. Now that it's open source,
I'll look at PRs, but the design direction is always going to be driven by my own
needs.

```
    G Major
e |-----3------|
B |-----0------|
G |-----0------|
D |-----0------|
A |-----2------|
E |-----3------|
```

It's a small terminal app. You build up a progression one chord at a time and save
it to a plain text file you can keep, paste, or share.

## Installation

```bash
go install github.com/krshultz/ripchords@latest
```

That drops a single self-contained binary on your `PATH`. No internet connection is
needed to run it — everything happens locally.

## Usage

Run `ripchords` with no arguments to start the editor.

**First run** walks you through a quick two-step setup:

1. **Input order** — how you like to type fret positions:
   - *Pitch order* (low to high): `E A D G B e` — e.g. `x 3 2 0 1 0` for C
   - *String-number order* (string 1–6): `e B G D A E` — e.g. `0 1 0 2 3 x` for C
2. **Barre chords** — whether to collapse a barre onto a single line where possible.

Your answers are saved, so you only do this once.

**Entering a chord** is two prompts — a name (optional), then the fret positions:

```
Chord name: C
Fret positions: x 3 2 0 1 0

    C
e |-----0------|
B |-----1------|
G |-----0------|
D |-----2------|
A |-----3------|
E |------------|
```

Keep going to build a progression. Press `s` to save it to a file when you're done.

## Hotkeys

From the progression view:

| Key  | Action                                                          |
|------|-----------------------------------------------------------------|
| `a`  | Add a chord (name, then fret positions)                         |
| `e`  | Edit a chord you already entered — rename it or fix its frets   |
| `l`  | Show the last chord on its own                                  |
| `s`  | Save the progression to a file                                  |
| `r`  | Reset the progression (clears chords, keeps your settings)      |
| `rr` | Double-tap `r` to wipe saved config and restart setup (asks first) |
| `?`  | Open settings (input order, barre rendering)                    |
| `q`  | Quit                                                            |

`l`, `?`, and `q` also work at the name/fret prompts when the input line is empty,
and `Esc` backs out of a prompt.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).

---

_Designed by Karl Shultz. Built in collaboration with [Claude Code](https://claude.com/claude-code) — I own the design and direction; the code was written with lots of AI assistance._
