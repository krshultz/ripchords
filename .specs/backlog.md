# Running list of goals for ripchords CLI

> Shipped items have been removed from this backlog (see git history / README for what's done).
> What remains below is open work or items that deviate from the original spec.

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

## Settings flow — remaining
* Storage location: the spec called for settings under `~/.ripchords/`. They currently
  live in `~/.config/ripchords/config.json`. Decide on the final location.
* New settings can be added over time; the toggle-based settings overlay already
  accommodates this. (Input-order and show-barre toggles are done.)

## Longer term goals
* Output formats beyond ASCII tab: JPG/PNG, and maybe an animated GIF of transitions.
  (ASCII tab + save-to-file is done.)
* Import a list of chords, perhaps as a CSV file.
* Internal database for reverse lookup: input fret positions, get a chord name back
  ("what chord is this?").
* Display the chords in a progression by name, e.g. "A minor --> C major --> DSus2".
