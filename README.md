# ripchords

> _Software for guitar players who want to write down what they just played_

---

## Overview
Hi, I'm Karl, the designer of Ripchords. I started this project because I needed it to solve some problems. I play the guitar, and I really enjoy improvising. Every now and then I stumble across a chord I really like, and want to write it down. Seems simple enough.

not for me though. I don't read music well at all, and when I do, I'm slow. I have trouble remembering the myriad ways to represent a chord visually - chord diagrams? Standard notation? Tablature? 

I find it easy to describe chords as a string of 6 characters, like `320003` for a simple G Major chord. So why not build something that can translate `320003` into something like this?

```
    G Major
e |-----3------|
B |-----0------|
G |-----0------|
D |-----0------|
A |-----2------|
E |-----3------|
```

## Installation
Put it wherever you want.

## Usage
`./ripchords` is a good place to start.

## Examples
```
➜  ripchords git:(single-chord) ./ripchords 
ripchords — entering frets in pitch order (E A D G B e)
Type 'quit' to exit, 'order' to change input order.

Chord name (or press Enter to skip): G Major
Fret positions: 320003

    G Major
e |-----3------|
B |-----0------|
G |-----0------|
D |-----0------|
A |-----2------|
E |-----3------|

Save to file? (Enter filename, or press Enter to skip): 

Chord name (or press Enter to skip): ^C
```

## Contributing

## License
