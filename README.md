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

Run `ripchords` with no arguments to start the editor. There's no setup step:
it opens straight at the first chord prompt.

Fret positions are always typed **low to high pitch**, one per string:
`E A D G B e`. So `x 3 2 0 1 0` is a C chord.

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
E |-----X------|
```

Type `x` for a string you don't play; it shows as `X` in the diagram. Fret positions
can be space-separated (`x 3 2 0 1 0`) or run together (`x32010`). Frets of 10 and
above need the spaces (`x 10 12 12 11 10`).

Keep going to build a progression. Press `s` to save it to a file when you're done.
Saving **appends** to the file if it already exists, so you can collect several
progressions in one file. A leading `~/` in the filename means your home directory.

## Command-line flags

| Flag                  | Action                       |
|-----------------------|------------------------------|
| `-h`, `--help`, `-?`  | Show usage and exit          |
| `-v`, `--version`     | Print the version and exit   |

## Hotkeys

From the progression view:

| Key  | Action                                                          |
|------|-----------------------------------------------------------------|
| `a`  | Add a chord (name, then fret positions)                         |
| `e`  | Edit a chord you already entered — rename it or fix its frets   |
| `l`  | Show the last chord on its own                                  |
| `s`  | Save the progression to a file                                  |
| `r`  | Reset the progression (clears chords, keeps your settings)      |
| `rr` | Double-tap `r` to reset settings to their defaults (asks first) |
| `?`  | Open settings (barre rendering)                                 |
| `q`  | Quit                                                            |

`l`, `?`, and `q` also work at the name/fret prompts when the input line is empty,
and `Esc` backs out of a prompt.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT — see [LICENSE](LICENSE).

---

_Designed by Karl Shultz. Built in collaboration with [Claude Code](https://claude.com/claude-code) — I own the design and direction; the code was written with lots of AI assistance._
