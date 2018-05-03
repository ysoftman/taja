// title : 타자 연습 게임
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

	gameView = NewView(gameViewStartX, gameViewEndX, gameViewStartY, gameViewEndY, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault)
	gameView.drawView()

	ibox = NewInputBox(inputBoxStartX, inputBoxEndX, inputBoxStartY, inputBoxEndY, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault)
	ibox.drawInputBox()

	statusView = NewView(statusViewStartX, statusViewEndX, statusViewStartY, statusViewEndY, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault)
	statusView.drawView()

	tempView = NewView(tempViewStartX, tempViewEndX, tempViewStartY, tempViewEndY, termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault)
	tempView.drawView()

	gameClearView = NewView(gameClearViewStartX, gameClearViewEndX, gameClearViewStartY, gameClearViewEndY, termbox.ColorBlack|termbox.AttrBold, termbox.ColorDefault)

	gameOverView = NewView(gameOverViewStartX, gameOverViewEndX, gameOverViewStartY, gameOverViewEndY, termbox.ColorBlack|termbox.AttrBold, termbox.ColorDefault)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return

			case <-time.After(1 * time.Second):
				updateStatus()
				if gameStatus == gameStatusNone {
					reset()
					render()
					continue
				}
				if gameStatus == gameStatusPlaying {
					// Add Enemy Word
					nw := getRandomWord()
					if nw.status == wordStatusCreated {
						fallingWords = append(fallingWords, nw)
					}
					// Refresh All Enemy Words
					for idx := range fallingWords {

						gameView.clearPrePos(fallingWords[idx])
						fallingWords[idx].y++
						if fallingWords[idx].y >= gameView.endy {
							fallingWords[idx].y = gameView.starty + 1
							if fallingWords[idx].status == wordStatusCreated {
								missCnt++
								if missCnt == liveCnt {
									gameStatus = gameStatusGameOver
									break

								}
							}
						}
						if fallingWords[idx].status == wordStatusCreated {
							gameView.printString(fallingWords[idx].x, fallingWords[idx].y, fallingWords[idx].str, termbox.ColorWhite)
						}
					}
					elapsedSec++
					continue
				}
				if gameStatus == gameStatusGameClear {
					showGameClear()
					continue
				} else if gameStatus == gameStatusGameOver {
					showGameOver()
					continue
				}
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
			case termbox.KeyCtrlN:
				gameStatus = gameStatusPlaying
				reset()
				gameView.clear()
				continue
			case termbox.KeyEnter:
				ibox.keyEnter()
				if deleteFallingWord(ibox.inputstr) {
					if checkGameClear() {
						gameStatus = gameStatusGameClear
						ibox.inputstr = ""
						continue
					}
					killCnt++

					gameScore += len(ibox.inputstr)
					matchLapSec = int(time.Now().Unix() - prelapsec)
					if matchLapSec == 0 {
						matchLapSec = 1
					}
					matchWordLen = len(ibox.inputstr)
					cpmValue = (matchWordLen * 60) / matchLapSec
					prelapsec = time.Now().Unix()
				}
				ibox.inputstr = ""
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
