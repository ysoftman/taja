package main

import (
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

func NewView(x1, x2, y1, y2 int, fg, bg termbox.Attribute) *View {
	return &View{
		startx:  x1,
		endx:    x2,
		starty:  y1,
		endy:    y2,
		fgcolor: fg,
		bgcolor: bg,
	}
}

func (v *View) drawFrame() {
	for y := v.starty; y < v.endy; y++ {
		for x := v.startx; x < v.endx; x++ {
			ch := ' '
			if y != v.starty && y != v.endy-1 && x != v.startx && x != v.endx-1 {
				continue
			}
			termbox.SetCell(x, y, ch, v.fgcolor, v.bgcolor)
		}
	}
}

func (v *View) printString(x, y int, str string, fgcolor termbox.Attribute) {
	for _, runeValue := range str {
		termbox.SetCell(v.startx+x, v.starty+y, runeValue, fgcolor, termbox.ColorDefault)
		w := runewidth.RuneWidth(runeValue)
		if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(runeValue)) {
			w = 1
		}
		x += w
	}
	render()
}

func (v *View) clear() {
	for y := v.starty; y < v.endy; y++ {
		for x := v.startx; x < v.endx; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
		}
	}
	render()
}

func (v *View) clearLine(y int) {
	for x := v.startx; x < v.endx; x++ {
		termbox.SetCell(x, v.starty+y, ' ', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}
	render()
}

func (v *View) clearWord(w word) {
	for x := 0; x < len(w.str); x++ {
		termbox.SetCell(w.x+x, w.y, ' ', termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}
	render()
}
