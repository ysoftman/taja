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
	msgNewGameCmd      = "NewGame : ctrl + n"
	msgQuitCmd         = "Quit : ctrl + c  or  ctrl + q"
	msgGameClear       = " Game Clear "
	msgGameOver        = " Game Over "
	paddingStr         = "     "

	killCnt      = 0
	missCnt      = 0
	liveCnt      = 0
	gameScore    = 0
	prelapsec    = time.Now().Unix()
	cpmValue     = 0
	matchWordLen = 0
	matchLapSec  = 0
	elapsedSec   = 0

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
	cpmValue = 0
	matchWordLen = 0
	matchLapSec = 0
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
		updateStatus()
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

func updateStatus() {
	line := 0
	line++
	statusView.printString(1, line, "[Status]", termbox.ColorDefault|termbox.AttrBold)
	line++
	statusView.printString(1, line, msgKillCnt+strconv.Itoa(killCnt), termbox.ColorDefault|termbox.AttrBold)
	line++
	statusView.printString(1, line, msgMissCnt+strconv.Itoa(missCnt)+" / "+strconv.Itoa(liveCnt), termbox.ColorDefault|termbox.AttrBold)
	cpmStatus := msgCPM + "matchWord(" + strconv.Itoa(matchWordLen) + ")*60 / lapsec(" + strconv.Itoa(matchLapSec) + ") = " + strconv.Itoa(cpmValue) + paddingStr
	line++
	statusView.printString(1, line, cpmStatus, termbox.ColorDefault|termbox.AttrBold)
	line++
	statusView.printString(1, line, msgGameScore+strconv.Itoa(gameScore), termbox.ColorDefault|termbox.AttrBold)
	line++
	statusView.printString(1, line, msgGameStatus+strconv.Itoa(gameStatus), termbox.ColorDefault|termbox.AttrBold)
	cnt := 0
	for _, v := range fallingWords {
		if v.status == wordStatusCreated {
			cnt++
		}
	}
	wordsStatus := msgLoadWordsCnt + strconv.Itoa(len(loadedWords)) + ", " + msgFallingWordsCnt + strconv.Itoa(cnt) + paddingStr
	line++
	statusView.printString(1, line, wordsStatus, termbox.ColorDefault|termbox.AttrBold)
	line++
	statusView.printString(1, line, msgElapsedSec+strconv.Itoa(elapsedSec), termbox.ColorDefault|termbox.AttrBold)
	line++
	line++
	statusView.printString(1, line, "[Command]", termbox.ColorYellow|termbox.AttrBold)
	line++
	statusView.printString(1, line, msgNewGameCmd, termbox.ColorYellow|termbox.AttrBold)
	line++
	statusView.printString(1, line, msgQuitCmd, termbox.ColorYellow|termbox.AttrBold)
	line++
	line++
	statusView.printString(1, line, time.Now().String(), termbox.ColorBlack|termbox.AttrBold)
}

func showGameClear() {
	gameClearView.drawView()
	gameClearView.printString(1, 1, msgGameClear, termbox.ColorRed|termbox.AttrBold)
}

func showGameOver() {
	gameOverView.drawView()
	gameOverView.printString(1, 1, msgGameOver, termbox.ColorRed|termbox.AttrBold)
}
