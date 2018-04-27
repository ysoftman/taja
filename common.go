package main

import termbox "github.com/nsf/termbox-go"

var (
	gameViewStartX = 0
	gameViewStartY = 0
	gameViewEndX   = 0
	gameViewEndY   = 0

	inputBoxStartX = 0
	inputBoxEndX   = 0
	inputBoxStartY = 0
	inputBoxEndY   = 0

	gameView *View
	ibox     *inputBox
)

func initConstValue() {
	mx, my := termbox.Size()

	gameViewStartX = 0
	gameViewStartY = 0
	gameViewEndX = mx - 1
	gameViewEndY = my - (my / 4)

	inputBoxStartX = 0
	inputBoxEndX = mx - 1
	inputBoxStartY = gameViewEndY
	inputBoxEndY = my - 1
}

func render() {
	gameView.drawMainVew()
	ibox.drawInputBox()

	termbox.Flush()
}
