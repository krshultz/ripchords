# Running list of goals for ripchords CLI

## Core design principles
* Ripchords must be entirely self contained. No internet access required at all.
* I want to start using semantic versioning.

## Things ripchords CLI should always do:
* It should always ask for chord name first, fret position second.
* Running it with --version should return a version number and a message. The version number shoiuld ideally be updated automatically as ripchords versions are released. The --version output should look like this. :
   ```
   Ripchords CLI v0.0.1 - software for guitar players
   ```
* Running it with -?, --help, -h, with an invalid option, or other common help flags should produce a usage screen. The usage screen should look like this to start with:
  ```
  ripchords v0.0.1
  ```

## Prompt enhancements for ripchords CLI
* The prompt needs to let the user configure their settings on the fly by entering `settings`
  * Each setting is presented as a toggle that the user can spacebar to select
  * When the settings changes are complete, they should take effect immediately, with no need to relaunch the app
  * Settings flow is dismissed by the `Esc` key, just like claude. Possibly might want a letter, if there are common letters to use like `b` to go back a level in a CLI menu

* Show the chords in the progression as part of the prompt. 
  * "Num. chords: x" where `x` == the number of chords so far.
  * The prompt should ahve a hotkey to display the prior chord.
  * Single-letter shortcuts to display the usage screen or the help (once the help exists)
  * A Start Over option. Useful for if the user views the progression so far and wants to start over, they can do so easily.

In both of these examples, the lines above "Chord name? (Enter to skip):" or "Next chord name? (Enter to skip):" should be the sticky headers. 

Example prompt on bringup:

```
Ripchords v0.0.1
Hotkeys: 
--> (s) opens settings
--> (l) last chord in progression
--> (h) help
Num chords: {displays the number of chords in the progression}
Chord name? (Enter to skip):
```

Example prompt once the user has added a single chord:

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

### Hotkey output
`(s)` key
Opens a bubbletea radio button flow to configure settings. I don't know how to represent that since it seems like it'll be a bubbletea thing and we don't have that yet

`(l)` key
Displays the last chord entered. In this example the user supplied the name `G Major` when they entered `32XXX3`:
```
    G Major
e |-----3------|
B |-----0------|
G |-----0------|
D |-----0------|
A |-----2------|
E |-----3------|

(Esc to dismiss)
```

`(p)` hotkey
Displays the progression in its current state, exactly like the program would print out when the user is finished. In this example there are three chords, only one of which was namned by the user:
```
     C Major                                  
e |-----0------| |---|-1------| |-----3------|
B |-----1------| |---|-1------| |-----3------|
G |-----0------| |---|-2------| |-----0------|
D |-----2------| |---|-3------| |-----0------|
A |-----3------| |---|-3------| |-----2------|
E |-----X------| |---|-1------| |-----3------|

(Esc to dismiss)
```

## Settings flow

### The basic idea

I want to add a Settings flow to ripchords which will let the user customize how it works. The number of customizations in the CLI will grow over time, so the codebase needs to accomodate that for easy improvements.

The settings should be stored in the user's `~/.ripchords` directory as a .json file.

### Examples of settings I want the user to be able to control
* Marking the barre with `|` should be configurable. If the user wants to be able to turn that feature off, they can.
* Entering in pitch order (the default, much more common) versus top-to-bottom

## Longer term goals

* Ability to output ASCII tab or JPGs, maybe even generate an animated gif of transitions?
* Ability to import a list of chords, perhaps as a CSV file
* Include an internal database to use as a reverse lookup, so users can input fret positions and ask "what chord is this?" And get a chord name back.
* Ways to display the chords in a progression by name. Like "A minor --> C major --> DSus2"