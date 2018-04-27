package main

import (
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

type View struct {
	startx int
	endx   int
	starty int
	endy   int
}

func NewView(x1, x2, y1, y2 int) *View {
	return &View{
		startx: x1,
		endx:   x2,
		starty: y1,
		endy:   y2,
	}
}

func (v *View) drawMainVew() {
	for y := v.starty; y < v.endy; y++ {
		for x := v.startx; x < v.endx; x++ {
			ch := '█'
			if y != 0 && y != v.endy-1 && x != 0 && x != v.endx-1 {
				continue
			}
			termbox.SetCell(x, y, ch, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault)
		}
	}
}

func (v *View) printString(x, y int, str string, fgcolor termbox.Attribute) {
	for _, runeValue := range str {
		termbox.SetCell(x, y, runeValue, fgcolor, termbox.ColorDefault)
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
			termbox.SetCell(x, y, ' ', termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault)
		}
	}
	render()
}
