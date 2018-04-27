// title : taja 타자 연습 게임
// author : ysoftman
// desc : 도스시절 한메 타자 산성비 게임을 생각하며...
package main

import (
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
)

var done = make(chan struct{})

func main() {
	startGame()
}

func startGame() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
	termbox.Flush()

	reset()
	setEnemyWords()

	gameView = NewView(gameViewStartX, gameViewEndX, gameViewStartY, gameViewEndY)
	gameView.drawMainVew()

	ibox = NewInputBox(inputBoxStartX, inputBoxEndX, inputBoxStartY, inputBoxEndY)
	ibox.drawInputBox()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return

			case <-time.After(1 * time.Second):
				// Add Enemy Word
				nw := getRandomWord()
				liveWords = append(liveWords, nw)

				// Refresh All Enemy Words
				for idx, _ := range liveWords {
					gameView.clearPrePos(liveWords[idx])
					liveWords[idx].y++
					if liveWords[idx].y >= gameView.endy {
						liveWords[idx].y = gameView.starty
					}
					gameView.printString(liveWords[idx].x, liveWords[idx].y, liveWords[idx].str, termbox.ColorWhite)
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
