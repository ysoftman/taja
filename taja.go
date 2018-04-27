// title : taja 타자 연습 게임
// author : ysoftman
// desc : 도스시절 한메 타자 산성비 게임을 생각하며...
package main

import (
	"sync"
	"time"

	"github.com/fatih/color"
	termbox "github.com/nsf/termbox-go"
)

var done = make(chan struct{})

func main() {
	startGame()
}

func GetColorString(cl, str string) string {
	switch cl {
	case "yellow":
		yellow := color.New(color.FgYellow).SprintFunc()
		return yellow(str)
	case "green":
		green := color.New(color.FgGreen).SprintFunc()
		return green(str)
	case "red":
		red := color.New(color.FgRed).SprintFunc()
		return red(str)
	case "blue":
		blue := color.New(color.FgBlue).SprintFunc()
		return blue(str)
	case "magenta":
		magenta := color.New(color.FgMagenta).SprintFunc()
		return magenta(str)
	case "cyan":
		cyan := color.New(color.FgCyan).SprintFunc()
		return cyan(str)
	default:
		white := color.New(color.FgWhite).SprintFunc()
		return white(str)
	}
}

type enemyWord struct {
	x    int
	y    int
	word string
}

func startGame() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Flush()

	initConstValue()
	gameView = NewView(gameViewStartX, gameViewEndX, gameViewStartY, gameViewEndY)
	gameView.drawMainVew()

	ibox = NewInputBox(inputBoxStartX, inputBoxEndX, inputBoxStartY, inputBoxEndY)
	ibox.drawInputBox()

	var wg sync.WaitGroup
	wg.Add(1)
	x, y := 1, 1
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case <-time.After(1 * time.Second):
				gameView.clear()
				gameView.printString(x, y, "가나다라마바사", termbox.ColorWhite)
				y++
				if x >= gameView.endx {
					x = gameView.startx
				}
				if y >= gameView.endy {
					y = gameView.starty
				}
				continue
			}
		}
	}()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyCtrlQ, termbox.KeyCtrlC:
				close(done)
				break mainloop
			case termbox.KeyEnter:
				ibox.keyEnter()
				continue
			case termbox.KeySpace:
				ibox.setChar(' ')
				continue
			case termbox.KeyDelete, termbox.KeyBackspace, termbox.KeyBackspace2:
				ibox.delChar()
				continue
			default:
				if ev.Ch != 0 {
					ibox.setChar(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}

	wg.Wait()
}
