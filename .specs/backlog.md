# Running list of goals for ripchords CLI

> When an item ships (merges to main), it moves to the "Recently shipped" section
> at the bottom — not deleted — so completed work stays visible in-file.
> Everything above that section is open work or items that deviate from the original spec.

## Core design principles
* Ripchords must be entirely self contained. No internet access required at all.

## Bugs / open issues
* (none currently — see GitHub issues for anything new)

## Prompt enhancements still open
* `(h)` help hotkey inside the editor (#51). The CLI usage screen now exists (shown for
  `-h`/`--help`/`-?`/invalid flags); what remains is surfacing it as an in-editor `h`
  hotkey alongside `l`, `?`, `r`, `q`.
* Settings hotkey binding (#52): the original spec called for `(s)` to open settings. Today
  `?` opens settings and `s` is bound to "save." Decide on the final bindings.

  Example prompt the spec was aiming for, once a chord has been added:
  ```
  Ripchords v0.0.1
  Hotkeys:
  --> (s) opens settings
  --> (l) last chord in progression
  --> (p) prints current progression
  --> (h) help
  --> (r) reset progression and start over
  Num chords: 1
  Next chord name? (Enter to skip):
  ```

  Note: the progression is already shown live above the prompt, so a dedicated `(p)`
  "print progression" hotkey may be redundant — decide whether it's still wanted.
* Show which chord is being worked on when adding or editing (#37). Name the pending chord
  on the add-frets prompt; on the edit-frets prompt show the current positions in place of
  the `(E A D G B e)` guide line once #48 lands.
* Make Rename discoverable (#46). Rename already works (`e` → pick chord → Rename) but the
  help bar only says `e edit chord`. Surface it in the help bar, `-h` usage text, and README.
* Reprint the fret positions the user entered (#44). Show a chord's positions as text
  (e.g. `x 3 2 0 1 0`) in the `l` last-chord view and the `e` picker. Needs a canonical
  fret formatter built from `Chord.Frets`; the same formatter should back #37 and #45.

## Editing & reuse (new)
* Cancel entering a chord — close the gaps (#53). `Esc` already cancels at the name step, but:
  - There's no on-screen hint that `Esc` cancels.
  - At the fret step, `Esc` only steps *back* to the name prompt; it doesn't abort the add.
  - When entering the very first chord (empty progression), there's no way to back out.
* Reuse a chord already entered this session (#49). When notating a progression that repeats
  chords (e.g. Bmin7 and DMaj several times each), the user shouldn't have to re-type the
  fret positions. Recall a prior chord by name. Open design questions:
  - Recall by typing the name (autocomplete) vs. picking from a list?
  - ~~Session-only memory, or persisted to disk between runs?~~ Persisted — #49 proposes a
    JSON chord library in `~/.ripchords/` that every entered chord is added to.
  - Keep the TUI lightweight — avoid clutter.
* Duplicate a chord in a progression (#54). Copy an existing chord in place so a repeated
  chord can be added without re-entering its fret positions.
* Move a chord in a progression (#42). Reorder chords within the progression (e.g. shift a
  chord earlier or later in the sequence) via a Move action in the `e` edit menu.
* Remove a chord from a progression (#41). Drop a single chord via the `e` edit menu
  instead of resetting the whole progression.
* Start a new progression from the current one — template / "Save As" (#43). Needs a
  loadable progression file format (JSON or CSV; no network) and a load path, since `s`
  today only writes ASCII output. Shares a parse/build path with the CSV-import and
  runtime-argument items under Longer term goals (#57, #58).

## Output
* Print a single chord diagram (#38). Select one chord from the progression and print
  only that diagram.
* Configurable display width (#40). Today diagrams print side by side in rows; add an
  option to print one chord per row for easier copy/paste of a single chord. Leaning
  toward a toggle (one per row vs. in series) over full customization.

## Settings flow — remaining
* Remove string-number input order — keep pitch order only (#48). String-number ordering
  (string 1→6) is rarely used; standard notation is pitch order (E A D G B e). Drop the
  `StringOrder` option, the first-run "choose input order" step, and the settings toggle,
  leaving pitch order as the sole behavior. Touches `core/chord.go`
  (`InputOrder`/`StringOrder`, `ParseFrets`), `core/session.go` (`Add`/`EditFrets`
  signatures), `main.go` (`Config.InputOrder` / `input_order` JSON), and `tui.go`
  (`stateFirstRun`, `settingsLabels`, the order toggle, and display strings). Consider how
  existing persisted `input_order` config values are handled on load (migrate/ignore).
* Show fret positions below each diagram (#45). A settings toggle (off by default) that
  prints each chord's positions under its diagram so pasted output is self-documenting.
  Depends on the fret formatter from #44.
* New settings can be added over time; the toggle-based settings overlay already
  accommodates this. (Show-barre toggle is done; the input-order toggle is being removed
  — see above.)

## Longer term goals
* Show the relative minor/major of a chord (#55). Added 07-13-2026: display the relative
  minor/major of a single chord, or maybe show it in a separate part of the app?
* Output formats beyond ASCII tab (#56): JPG/PNG, and maybe an animated GIF of transitions.
  (ASCII tab + save-to-file is done.)
* Copy output straight to the clipboard (#39), so diagrams can be pasted into an email or
  forum post without going through a saved file.
* Ingest chord fret positions as a runtime argument (#57) (pass a chord directly on the
  command line instead of entering it interactively). Ideally the same parsing/build
  logic backs both this runtime flow and the "import from CSV" flow below.
* Import a list of chords, perhaps as a CSV file (#58). (Same file format / load path as #43.)
* Internal database for reverse lookup (#59): input fret positions, get a chord name back
  ("what chord is this?"). The chord library in #49 is the natural first data source.
* Display the chords in a progression by name (#60), e.g. "A minor --> C major --> DSus2".

## Recently shipped
* Misaligned output when frets mix single- and double-digit numbers (#33, fixed in #35,
  released in 0.2.4). Every fret cell is now padded to the width of the widest fret in the
  diagram so the closing `|` aligns across strings.
* Clear error for frets >9 in compact input (#26). Spaceless input like `101212111010`
  is ambiguous for two-digit frets; it now emits a targeted message telling the user to
  space-separate (e.g. `x 10 12 12 11 10`) instead of a confusing token-count error.
* Module-path rename (#24). `go.mod` module became `github.com/krshultz/ripchords`
  (was bare `ripchords`) so `go install github.com/krshultz/ripchords@latest` works.
* Edit a fret pattern instead of restarting (#13). From the rendered screen, `e` re-opens
  an already-submitted chord to replace its fret positions — no need to re-enter the whole
  progression to fix a typo.
* Rename a chord (#13). `e` → pick a chord → Rename changes the name of a chord already
  entered this session without re-entering fret positions.
* Error messages wrap to the terminal width (#30, released in 0.2.3) instead of running off
  the edge of the screen.
* README rewritten for the TUI, MIT LICENSE added (#23), and README updated now that the
  project is open source (#29).
* First-run setup prompts for the barre option, and `rr` asks for confirmation before
  wiping config instead of relying on a 500ms double-tap (#17, #18; fixed in #20, released
  in 0.2.1).
* Dev tooling (#19): golangci-lint, CI and PR-check workflows (Conventional Commits PR
  titles, branch names), a pre-commit hook, and CONTRIBUTING.md.
