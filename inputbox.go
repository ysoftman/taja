package main

import (
	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

func NewInputBox(x1, x2, y1, y2 int, fg, bg termbox.Attribute) *inputBox {
	return &inputBox{
		startx:  x1,
		endx:    x2,
		starty:  y1,
		endy:    y2,
		curpos:  x1 + 1,
		fgcolor: fg,
		bgcolor: bg,
	}
}

func (ib *inputBox) drawInputBox() {
	for y := ib.starty; y < ib.endy; y++ {
		for x := ib.startx; x < ib.endx; x++ {
			ch := '█'
			if y != ib.starty && y != ib.endy-1 && x != ib.startx && x != ib.endx-1 {
				continue
			}
			termbox.SetCell(x, y, ch, ib.fgcolor, ib.bgcolor)
		}
	}
}

func (ib *inputBox) clear() {
	for y := ib.starty; y < ib.endy; y++ {
		for x := ib.startx; x < ib.endx; x++ {
			termbox.SetCell(x, y, ' ', termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault)
		}
	}
}
func (ib *inputBox) setChar(runeValue rune) {
	if ib.curpos >= ib.endx {
		return
	}

	ib.inputstr += string(runeValue)

	termbox.SetCell(ib.curpos, ib.starty+1, runeValue, termbox.ColorDefault, termbox.ColorDefault)
	w := runewidth.RuneWidth(runeValue)
	if w == 0 || (w == 2 && runewidth.IsAmbiguousWidth(runeValue)) {
		w = 1
	}
	ib.curpos += w
	termbox.SetCursor(ib.curpos, ib.starty+1)
	render()
}

func (ib *inputBox) delChar() {
	if ib.curpos <= 1 {
		return
	}

	ib.inputstr = ib.inputstr[:1]

	termbox.SetCell(ib.curpos-1, ib.starty+1, ' ', termbox.ColorDefault, termbox.ColorDefault)
	ib.curpos--
	termbox.SetCursor(ib.curpos, ib.starty+1)
	render()
}

func (ib *inputBox) keyEnter() {
	ib.clear()
	ib.curpos = ib.startx + 1
	termbox.SetCursor(ib.startx+1, ib.starty+1)
	render()
}
