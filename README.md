# ripchords

> _Software for guitar players who want to write down what they just played_

---

## Overview
Hi, I'm Karl, the designer of Ripchords. I started this project because I needed it to solve some problems. I play the guitar, and I really enjoy improvising. Every now and then I stumble across a chord I really like, and want to write it down. Seems simple enough, right?

It's not simple for me. I don't read music well at all, and when I do, I'm slow. I have trouble remembering the myriad ways to represent a chord visually - chord diagrams? Standard notation? Tablature? If I do a diagram, how is it oriented? There are so many ways.

I find it easy to describe chords as a string of 6 characters, like `320003` for a simple G Major chord. So why not build something that can translate `320003` into [ASCII tab](https://en.wikipedia.org/wiki/ASCII_tab)? 

```
    G Major
e |-----3------|
B |-----0------|
G |-----0------|
D |-----0------|
A |-----2------|
E |-----3------|
```

This way I'm using a somewhat standard format that I see on the internet all the time.  And maybe I'd be able to use the ASCII tab diagrams programmatically as well. Here's what they look like:

## Installation
Put it wherever you want. If you build it from source it'll just dump the binary right there onto `$PWD`.

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

Not necessary just yet.

## License

Not necessary just yet.
