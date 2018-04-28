// title : 타자 연습 게임
// author : ysoftman
// desc : 도스시절 한메 타자 산성비 게임을 생각하며...
package main

import (
	"strconv"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
)

var done = make(chan struct{})

func main() {
	startGame()
}

func updateKillCnt(n int) {
	statusView.printString(1, 1, msgKillCnt+strconv.Itoa(n), termbox.ColorDefault)
}

func updateMissCnt(n int) {
	statusView.printString(1, 2, msgMissCnt+strconv.Itoa(n), termbox.ColorDefault)
}

func updateCPM(n int) {
	statusView.printString(1, 3, msgCPM+strconv.Itoa(n), termbox.ColorDefault)
}

func updateGameScore(n int) {
	statusView.printString(1, 4, msgGameScore+strconv.Itoa(n), termbox.ColorDefault)
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

	gameView = NewView(gameViewStartX, gameViewEndX, gameViewStartY, gameViewEndY, termbox.ColorGreen|termbox.AttrBold, termbox.ColorDefault)
	gameView.drawView()

	ibox = NewInputBox(inputBoxStartX, inputBoxEndX, inputBoxStartY, inputBoxEndY, termbox.ColorYellow|termbox.AttrBold, termbox.ColorDefault)
	ibox.drawInputBox()

	statusView = NewView(statusViewStartX, statusViewEndX, statusViewStartY, statusViewEndY, termbox.ColorMagenta|termbox.AttrBold, termbox.ColorDefault)
	statusView.drawView()

	cmdView = NewView(cmdViewStartX, cmdViewEndX, cmdViewStartY, cmdViewEndY, termbox.ColorCyan|termbox.AttrBold, termbox.ColorDefault)
	cmdView.drawView()

	debugView = NewView(debugViewStartX, debugViewEndX, debugViewStartY, debugViewEndY, termbox.ColorRed|termbox.AttrBold, termbox.ColorDefault)
	debugView.drawView()

	cmdView.printString(1, 1, msgQuitCmd, termbox.ColorDefault)
	cmdView.printString(1, 2, msgURL, termbox.ColorDefault)

	updateKillCnt(0)
	updateMissCnt(0)
	updateGameScore(0)
	updateCPM(0)

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
							updateMissCnt(missCnt)
						}
					}
					if fallingWords[idx].status == wordStatusCreated {
						gameView.printString(fallingWords[idx].x, fallingWords[idx].y, fallingWords[idx].str, termbox.ColorWhite)
					}
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
				// debugView.debug(ibox.inputstr)
				if deleteFallingWord(ibox.inputstr) {
					// debugView.debug(ibox.inputstr)
					// debugView.debug(strconv.Itoa(len(fallingWords)))
					if checkGameClear() {
						gameView.printString(gameView.endx/2-len(msgGameClear), gameView.endy/2, msgGameClear, termbox.ColorWhite)
						close(done)
						break mainloop
					}
					killCnt++
					updateKillCnt(killCnt)
					gameScore += len(ibox.inputstr)
					updateGameScore(gameScore)

					cpm = int(len(ibox.inputstr) * 60 / int(time.Now().Unix()-prelapsec))
					prelapsec = time.Now().Unix()
					updateCPM(cpm)
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
