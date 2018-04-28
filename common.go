package main

import (
	"math/rand"
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

	cmdViewStartX = 0
	cmdViewEndX   = 0
	cmdViewStartY = 0
	cmdViewEndY   = 0

	debugViewStartX = 0
	debugViewEndX   = 0
	debugViewStartY = 0
	debugViewEndY   = 0

	loadedWords  []word
	fallingWords []word
	gameView     *View
	statusView   *View
	cmdView      *View
	debugView    *View
	ibox         *inputBox

	msgKillCnt   = "Kill Words : "
	msgMissCnt   = "Miss Words : "
	msgCPM       = "CPM(CharacterPerMinute) : "
	msgGameScore = "Score : "

	msgQuitCmd   = "Quit : ctrl + c  or  ctrl + q"
	msgURL       = "http://github.com/ysoftman/taja"
	msgGameClear = "Game Clear"

	killCnt   = 0
	missCnt   = 0
	gameScore = 0
	prelapsec = time.Now().Unix()
	cpm       = 0
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

	cmdViewStartX = statusViewEndX
	cmdViewEndX = mx
	cmdViewStartY = inputBoxEndY
	cmdViewEndY = my - 3

	debugViewStartX = statusViewEndX
	debugViewEndX = mx
	debugViewStartY = cmdViewEndY
	debugViewEndY = my

	rand.Seed(time.Now().UnixNano())
}

func render() {
	gameView.drawView()
	ibox.drawInputBox()
	statusView.drawView()
	cmdView.drawView()
	debugView.drawView()
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
