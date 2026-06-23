# Running list of goals for ripchords CLI

> When an item ships (merges to main), it moves to the "Recently shipped" section
> at the bottom — not deleted — so completed work stays visible in-file.
> Everything above that section is open work or items that deviate from the original spec.

## Core design principles
* Ripchords must be entirely self contained. No internet access required at all.

## Prompt enhancements still open
* `(h)` help hotkey inside the editor. The CLI usage screen now exists (shown for
  `-h`/`--help`/`-?`/invalid flags); what remains is surfacing it as an in-editor `h`
  hotkey alongside `l`, `?`, `r`, `q`.
* Settings hotkey binding: the original spec called for `(s)` to open settings. Today
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

## Editing & reuse (new)
* Cancel entering a chord — close the gaps. `Esc` already cancels at the name step, but:
  - There's no on-screen hint that `Esc` cancels.
  - At the fret step, `Esc` only steps *back* to the name prompt; it doesn't abort the add.
  - When entering the very first chord (empty progression), there's no way to back out.
* Reuse a chord already entered this session. When notating a progression that repeats
  chords (e.g. Bmin7 and DMaj several times each), the user shouldn't have to re-type the
  fret positions. Recall a prior chord by name. Open design questions:
  - Recall by typing the name (autocomplete) vs. picking from a list?
  - Session-only memory, or persisted to disk between runs?
  - Keep the TUI lightweight — avoid clutter.

## Settings flow — remaining
* New settings can be added over time; the toggle-based settings overlay already
  accommodates this. (Input-order and show-barre toggles are done.)

## Longer term goals
* Output formats beyond ASCII tab: JPG/PNG, and maybe an animated GIF of transitions.
  (ASCII tab + save-to-file is done.)
* Ingest chord fret positions as a runtime argument (pass a chord directly on the
  command line instead of entering it interactively). Ideally the same parsing/build
  logic backs both this runtime flow and the "import from CSV" flow below.
* Import a list of chords, perhaps as a CSV file.
* Internal database for reverse lookup: input fret positions, get a chord name back
  ("what chord is this?").
* Display the chords in a progression by name, e.g. "A minor --> C major --> DSus2".

## Recently shipped
* Edit a fret pattern instead of restarting (#13). From the rendered screen, `e` re-opens
  an already-submitted chord to replace its fret positions — no need to re-enter the whole
  progression to fix a typo.
* Rename a chord (#13). `e` → pick a chord → Rename changes the name of a chord already
  entered this session without re-entering fret positions.
