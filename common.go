package main

import (
	"math/rand"
	"strconv"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type View struct {
	startx  int
	endx    int
	starty  int
	endy    int
	fgcolor termbox.Attribute
	bgcolor termbox.Attribute
}

type inputBox struct {
	startx   int
	endx     int
	starty   int
	endy     int
	curpos   int
	inputstr string
	fgcolor  termbox.Attribute
	bgcolor  termbox.Attribute
}

type word struct {
	status int
	x      int
	y      int
	str    string
}

const (
	wordStatusEmpty = iota
	wordStatusCreated
	wordStatusDeleted
	wordStatusFreeing
)

const (
	gameStatusNone = iota
	gameStatusPlaying
	gameStatusQuit
	gameStatusGameClear
	gameStatusGameOver
)

var (
	gameViewStartX = 0
	gameViewStartY = 0
	gameViewEndX   = 0
	gameViewEndY   = 0

	inputBoxStartX = 0
	inputBoxEndX   = 0
	inputBoxStartY = 0
	inputBoxEndY   = 0

	statusViewStartX = 0
	statusViewEndX   = 0
	statusViewStartY = 0
	statusViewEndY   = 0

	tempViewStartX = 0
	tempViewEndX   = 0
	tempViewStartY = 0
	tempViewEndY   = 0

	gameClearViewStartX = 0
	gameClearViewEndX   = 0
	gameClearViewStartY = 0
	gameClearViewEndY   = 0

	gameOverViewStartX = 0
	gameOverViewEndX   = 0
	gameOverViewStartY = 0
	gameOverViewEndY   = 0

	loadedWords   []word
	fallingWords  []word
	gameView      *View
	statusView    *View
	tempView      *View
	gameClearView *View
	gameOverView  *View
	ibox          *inputBox

	msgKillCnt         = "Kill Words : "
	msgMissCnt         = "Miss Words : "
	msgCPM             = "CPM(CharacterPerMinute) : "
	msgGameScore       = "Score : "
	msgGameStatus      = "GameStatus : "
	msgLoadWordsCnt    = "LoadWords : "
	msgFallingWordsCnt = "FallingWords : "
	msgElapsedSec      = "ElapsedSec : "

	msgNewGameCmd = "NewGame : ctrl + n"
	msgQuitCmd    = "Quit : ctrl + c  or  ctrl + q"
	msgURL        = "Github : http://github.com/ysoftman/taja"
	msgGameClear  = " Game Clear "
	msgGameOver   = " Game Over "

	killCnt    = 0
	missCnt    = 0
	liveCnt    = 0
	gameScore  = 0
	prelapsec  = time.Now().Unix()
	cpm        = 0
	elapsedSec = 0

	gameStatus = gameStatusNone
)

func reset() {
	mx, my := termbox.Size()

	gameViewStartX = 0
	gameViewStartY = 0
	gameViewEndX = mx
	gameViewEndY = my - (my / 4)

	inputBoxStartX = 0
	inputBoxEndX = mx
	inputBoxStartY = gameViewEndY
	inputBoxEndY = inputBoxStartY + 3

	statusViewStartX = 0
	statusViewEndX = mx / 2
	statusViewStartY = inputBoxEndY
	statusViewEndY = my

	tempViewStartX = statusViewEndX
	tempViewEndX = mx
	tempViewStartY = inputBoxEndY
	tempViewEndY = my

	gameClearViewStartX = (gameViewEndX / 2) - len(msgGameClear) - 1
	gameClearViewEndX = gameClearViewStartX + len(msgGameClear) + 2
	gameClearViewStartY = (gameViewEndY / 2) - 3
	gameClearViewEndY = gameClearViewStartY + 3

	gameOverViewStartX = (gameViewEndX / 2) - len(msgGameOver) - 1
	gameOverViewEndX = gameOverViewStartX + len(msgGameOver) + 2
	gameOverViewStartY = (gameViewEndY / 2) - 3
	gameOverViewEndY = gameOverViewStartY + 3

	killCnt = 0
	missCnt = 0
	liveCnt = 5
	gameScore = 0
	prelapsec = time.Now().Unix()
	cpm = 0
	elapsedSec = 0

	rand.Seed(time.Now().UnixNano())

	loadedWords = []word{}
	setEnemyWords()

	fallingWords = []word{}

	if gameView != nil {
		gameView.clear()
	}
	if ibox != nil {
		ibox.clear()
	}
	if statusView != nil {
		statusView.clear()
		updateKillCnt(0)
		updateMissCnt(0)
		updateGameScore(0)
		updateCPM(0)
		updateGameStatus(0)
		updateWordStatus()
		updateCommand()
		updateElapsedSec(0)
	}
	if tempView != nil {
		tempView.clear()
	}

}

func render() {
	gameView.drawView()
	ibox.drawInputBox()
	statusView.drawView()
	tempView.drawView()
	termbox.Flush()
}

func setEnemyWords() {
	loadedWords = []word{
		{wordStatusCreated, 0, 0, "apple"},
		{wordStatusCreated, 0, 0, "lemon"},
		{wordStatusCreated, 0, 0, "okay"},
		{wordStatusCreated, 0, 0, "love"},
		{wordStatusCreated, 0, 0, "slow"},
		{wordStatusCreated, 0, 0, "golang"},
		{wordStatusCreated, 0, 0, "rainbow"},
		{wordStatusCreated, 0, 0, "fruite"},
		{wordStatusCreated, 0, 0, "bicycle"},
		{wordStatusCreated, 0, 0, "train"},
		{wordStatusCreated, 0, 0, "car"},
		{wordStatusCreated, 0, 0, "level"},
		{wordStatusCreated, 0, 0, "superman"},
	}
}

func getRandomWord() word {
	n := len(loadedWords)
	if n <= 0 {
		return word{wordStatusEmpty, 0, 0, ""}
	}
	i := rand.Intn(n)
	w := loadedWords[i]
	w.x = rand.Intn(gameViewEndX-len(w.str)-1) + 1

	// delete return word element in loadedWords
	loadedWords = append(loadedWords[:i], loadedWords[i+1:]...)

	return w
}

func deleteFallingWord(str string) bool {
	for i := range fallingWords {
		if fallingWords[i].str == str {
			fallingWords[i].status = wordStatusDeleted
			return true
		}
	}
	return false
}

func checkGameClear() bool {
	if len(fallingWords) == 0 {
		return true
	}
	for i := range fallingWords {
		if fallingWords[i].status != wordStatusDeleted {
			return false
		}
	}
	return true
}

func updateKillCnt(n int) {
	statusView.printString(1, 1, msgKillCnt+strconv.Itoa(n), termbox.ColorDefault|termbox.AttrBold)
}

func updateMissCnt(n int) {
	statusView.printString(1, 2, msgMissCnt+strconv.Itoa(n)+" / "+strconv.Itoa(liveCnt), termbox.ColorDefault|termbox.AttrBold)
}

func updateCPM(n int) {
	statusView.printString(1, 3, msgCPM+strconv.Itoa(n), termbox.ColorDefault|termbox.AttrBold)
}

func updateGameScore(n int) {
	statusView.printString(1, 4, msgGameScore+strconv.Itoa(n), termbox.ColorDefault|termbox.AttrBold)
}

func updateGameStatus(gameStatus int) {
	statusView.printString(1, 5, msgGameStatus+strconv.Itoa(gameStatus), termbox.ColorDefault|termbox.AttrBold)
}

func updateWordStatus() {
	wordsStatus := msgLoadWordsCnt + strconv.Itoa(len(loadedWords)) + ", "

	cnt := 0
	for _, v := range fallingWords {
		if v.status == wordStatusCreated {
			cnt++
		}
	}
	wordsStatus += msgFallingWordsCnt + strconv.Itoa(cnt) + "     "
	statusView.printString(1, 6, wordsStatus, termbox.ColorDefault|termbox.AttrBold)
}

func updateElapsedSec(sec int) {
	statusView.printString(1, 7, msgElapsedSec+strconv.Itoa(sec), termbox.ColorDefault|termbox.AttrBold)
}

func updateCommand() {
	statusView.printString(1, 9, "[Command]", termbox.ColorYellow|termbox.AttrBold)
	statusView.printString(1, 10, msgNewGameCmd, termbox.ColorYellow|termbox.AttrBold)
	statusView.printString(1, 11, msgQuitCmd, termbox.ColorYellow|termbox.AttrBold)
	statusView.printString(1, 12, msgURL, termbox.ColorYellow)
}

func showGameClear() {
	gameClearView.drawView()
	gameClearView.printString(1, 1, msgGameClear, termbox.ColorRed|termbox.AttrBold)
}

func showGameOver() {
	gameOverView.drawView()
	gameOverView.printString(1, 1, msgGameOver, termbox.ColorRed|termbox.AttrBold)
}
