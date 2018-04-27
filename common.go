package main

import (
	"math/rand"
	"time"

	"github.com/fatih/color"
	termbox "github.com/nsf/termbox-go"
)

type word struct {
	x   int
	y   int
	str string
}

var (
	gameViewStartX = 0
	gameViewStartY = 0
	gameViewEndX   = 0
	gameViewEndY   = 0

	inputBoxStartX = 0
	inputBoxEndX   = 0
	inputBoxStartY = 0
	inputBoxEndY   = 0

	enemyWords []word
	liveWords  []word
	gameView   *View
	ibox       *inputBox
)

func reset() {
	mx, my := termbox.Size()

	gameViewStartX = 0
	gameViewStartY = 0
	gameViewEndX = mx - 1
	gameViewEndY = my - (my / 4)

	inputBoxStartX = 0
	inputBoxEndX = mx - 1
	inputBoxStartY = gameViewEndY
	inputBoxEndY = my - 1

	rand.Seed(time.Now().UnixNano())
}

func render() {
	gameView.drawMainVew()
	ibox.drawInputBox()

	termbox.Flush()
}

func getColorString(cl, str string) string {
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

func setEnemyWords() {
	enemyWords = []word{
		{0, 0, "apple"},
		{0, 0, "lemon"},
		{0, 0, "okay"},
		{0, 0, "love"},
		{0, 0, "slow"},
		{0, 0, "golang"},
		{0, 0, "rainbow"},
		{0, 0, "fruite"},
		{0, 0, "bicycle"},
		{0, 0, "train"},
		{0, 0, "car"},
		{0, 0, "level"},
		{0, 0, "superman"},
	}
}

func getRandomWord() word {
	n := len(enemyWords)
	if n <= 0 {
		return word{0, 0, "---"}
	}
	i := rand.Intn(n)
	w := enemyWords[i]
	w.x = rand.Intn(gameViewEndX - len(w.str))

	// delete return word element in enemyWords
	enemyWords = append(enemyWords[:i], enemyWords[i+1:]...)

	return w
}
