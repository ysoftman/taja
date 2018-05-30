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
	gameStatusLevelUp
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

	levelupViewStartX = 0
	levelupViewEndX   = 0
	levelupViewStartY = 0
	levelupViewEndY   = 0

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
	levelupView   *View
	gameClearView *View
	gameOverView  *View
	ibox          *inputBox

	msgKillCnt         = "Kill Words : "
	msgMissCnt         = "Miss Words : "
	msgCPM             = "CPM(CharacterPerMinute) : "
	msgGameScore       = "Score : "
	msgGameStatus      = "GameStatus : "
	msgLevel           = "Level : "
	msgLoadWordsCnt    = "LoadWords : "
	msgFallingWordsCnt = "FallingWords : "
	msgElapsedSec      = "ElapsedSec : "
	msgNewGameCmd      = "NewGame : ctrl + n"
	msgQuitCmd         = "Quit : ctrl + c  or  ctrl + q"
	msgLevelUp         = " Level Up : "
	msgGameClear       = " Game Clear "
	msgGameOver        = " Game Over "
	paddingStr         = "     "

	killCnt      = 0
	missCnt      = 0
	liveCnt      = 0
	gameScore    = 0
	gameLevel    = 1
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

	levelupViewStartX = (gameViewEndX / 2) - len(msgLevelUp) - 1
	levelupViewEndX = gameClearViewStartX + len(msgLevelUp) + 4
	levelupViewStartY = (gameViewEndY / 2) - 3
	levelupViewEndY = gameClearViewStartY + 3

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
	liveCnt = 20
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
	gameView.drawFrame()
	ibox.drawInputBox()
	statusView.drawFrame()
	tempView.drawFrame()
	termbox.Flush()
}

func setEnemyWords() {
	switch gameLevel {
	case 1:
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
			{wordStatusCreated, 0, 0, "xelloss"},
		}
	case 2:
		loadedWords = []word{
			{wordStatusCreated, 0, 0, "red"},
			{wordStatusCreated, 0, 0, "orange"},
			{wordStatusCreated, 0, 0, "yellow"},
			{wordStatusCreated, 0, 0, "green"},
			{wordStatusCreated, 0, 0, "blue"},
			{wordStatusCreated, 0, 0, "indigo"},
			{wordStatusCreated, 0, 0, "violet"},
			{wordStatusCreated, 0, 0, "white"},
			{wordStatusCreated, 0, 0, "black"},
			{wordStatusCreated, 0, 0, "color"},
			{wordStatusCreated, 0, 0, "pentagon"},
			{wordStatusCreated, 0, 0, "hexagon"},
			{wordStatusCreated, 0, 0, "heptagon"},
			{wordStatusCreated, 0, 0, "octagon"},
			{wordStatusCreated, 0, 0, "nonagon"},
			{wordStatusCreated, 0, 0, "decagon"},
			{wordStatusCreated, 0, 0, "hendecagon"},
			{wordStatusCreated, 0, 0, "dodecagon"},
			{wordStatusCreated, 0, 0, "circle"},
			{wordStatusCreated, 0, 0, "square"},
			{wordStatusCreated, 0, 0, "rectangle"},
			{wordStatusCreated, 0, 0, "oval"},
			{wordStatusCreated, 0, 0, "heart"},
			{wordStatusCreated, 0, 0, "abcdefghijklmnopqrstuvwxyz"},
		}
	case 3:
		loadedWords = []word{
			{wordStatusCreated, 0, 0, "There is no spoon."},
			{wordStatusCreated, 0, 0, "Stay hungry stay foolish."},
			{wordStatusCreated, 0, 0, "What time is it now?"},
			{wordStatusCreated, 0, 0, "I will be back."},
			{wordStatusCreated, 0, 0, "I'm your father"},
			{wordStatusCreated, 0, 0, "Make it count."},
			{wordStatusCreated, 0, 0, "I want to live with the flow"},
			{wordStatusCreated, 0, 0, "If the sun were to rise in the west, my love would be unchaged forever"},
			{wordStatusCreated, 0, 0, "Keep Calm and Carry On"},
			{wordStatusCreated, 0, 0, "Let it go, the cold never bothered me anyway"},
		}
	case 4:
		loadedWords = []word{
			{wordStatusCreated, 0, 0, "우주용사"},
			{wordStatusCreated, 0, 0, "미리내"},
			{wordStatusCreated, 0, 0, "꽃다발"},
			{wordStatusCreated, 0, 0, "허수아비"},
			{wordStatusCreated, 0, 0, "연못"},
			{wordStatusCreated, 0, 0, "사자"},
			{wordStatusCreated, 0, 0, "개구리"},
			{wordStatusCreated, 0, 0, "걸음걸이"},
			{wordStatusCreated, 0, 0, "마차"},
			{wordStatusCreated, 0, 0, "신데렐라"},
			{wordStatusCreated, 0, 0, "오리온"},
			{wordStatusCreated, 0, 0, "초코파이"},
			{wordStatusCreated, 0, 0, "아롱다롱"},
			{wordStatusCreated, 0, 0, "후레쉬맨"},
			{wordStatusCreated, 0, 0, "독수리오형제"},
		}
	case 5:
		loadedWords = []word{
			{wordStatusCreated, 0, 0, "가는 말이 고와야 오는 말이 곱다."},
			{wordStatusCreated, 0, 0, "천리길도 한걸음 부터."},
			{wordStatusCreated, 0, 0, "내가 너를 모르는데 넌들 나를 알겠느냐?"},
			{wordStatusCreated, 0, 0, "오 그대여, 가지 마세요. 나는 지금 울잖아요."},
			{wordStatusCreated, 0, 0, "사랑이 어떻게 변하니?"},
			{wordStatusCreated, 0, 0, "간장 공장 공장장은 강 공장장이고 된장 공장 공장장은 공 공장장이다."},
			{wordStatusCreated, 0, 0, "이건 입에서 나는 소리가 아니야."},
			{wordStatusCreated, 0, 0, "가나다라마바사아자차카타파하"},
			{wordStatusCreated, 0, 0, "똠방각하"},
			{wordStatusCreated, 0, 0, "동해물과 백두산이 마르고 닳도록"},
			{wordStatusCreated, 0, 0, "지못미, 지켜주지 못해 미안해"},
			{wordStatusCreated, 0, 0, "마상, 마음의 상처"},
			{wordStatusCreated, 0, 0, "아도겐 아따따뚜겐"},
			{wordStatusCreated, 0, 0, "오늘도 행복하세요."},
			{wordStatusCreated, 0, 0, "성공한 사람의 인생은 성공한 후에 포장되어 평범한 사람을 망친다."},
			{wordStatusCreated, 0, 0, "마른 하늘에 날벼락."},
			{wordStatusCreated, 0, 0, "나한테 왜 그랬어요? 넌 나에게 모욕감을 줬어"},
			{wordStatusCreated, 0, 0, "느그 아부지 뭐하시노?"},
			{wordStatusCreated, 0, 0, "밥은 먹고 다니냐?"},
			{wordStatusCreated, 0, 0, "너나 잘 하세요."},
		}
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
	line++
	statusView.printString(1, line, msgLevel+strconv.Itoa(gameLevel), termbox.ColorDefault|termbox.AttrBold)
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
	// line++
	// statusView.printString(1, line, time.Now().String(), termbox.ColorWhite|termbox.AttrBold)
	line++
	line++

	cmdColor := termbox.ColorYellow | termbox.AttrBold
	if time.Now().Second()%5 == 0 {
		cmdColor = termbox.ColorBlue | termbox.AttrBold
	} else if time.Now().Second()%5 == 1 {
		cmdColor = termbox.ColorCyan | termbox.AttrBold
	} else if time.Now().Second()%5 == 2 {
		cmdColor = termbox.ColorGreen | termbox.AttrBold
	} else if time.Now().Second()%5 == 3 {
		cmdColor = termbox.ColorMagenta | termbox.AttrBold
	} else if time.Now().Second()%5 == 4 {
		cmdColor = termbox.ColorRed | termbox.AttrBold
	}

	statusView.printString(1, line, "[Command]", cmdColor)
	line++
	statusView.printString(1, line, msgNewGameCmd, cmdColor)
	line++
	statusView.printString(1, line, msgQuitCmd, cmdColor)
}

func showLevelUp() {
	levelupView.drawFrame()
	levelupView.printString(1, 1, msgLevelUp+strconv.Itoa(gameLevel), termbox.ColorRed|termbox.AttrBold)
}

func showGameClear() {
	gameClearView.drawFrame()
	gameClearView.printString(1, 1, msgGameClear, termbox.ColorRed|termbox.AttrBold)
}

func showGameOver() {
	gameOverView.drawFrame()
	gameOverView.printString(1, 1, msgGameOver, termbox.ColorRed|termbox.AttrBold)
}
