// title : taja 타자 연습 게임
// author : ysoftman
// desc : 도스시절 한메 타자 산성비 게임을 생각하며...
package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	// "github.com/jroimartin/gocui" // 한글(utf8) 출력에 문제가 있음
	"github.com/ysoftman/gocui"
)

var done = make(chan struct{})

func main() {
	StartGoCui()
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

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("words view", 0, 0, maxX-2, maxY-11); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "taja game"
		v.SetCursor(((maxY / 2) - len(v.Title)/2), 0)

	}
	if v, err := g.SetView("input area", 0, maxY-10, maxX-2, maxY-8); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Highlight = true
		v.Frame = true
		v.SetCursor(0, 0)
		g.SetCurrentView("input area")
	}
	if v, err := g.SetView("status area", 0, maxY-7, maxX-2, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "status"
		v.SetCursor(0, 0)
		g.FgColor = gocui.ColorGreen
	}

	return nil
}

func inputAction(g *gocui.Gui, v *gocui.View) error {
	g.Update(func(g *gocui.Gui) error {
		inputView, _ := g.View("input view")
		word := strings.TrimSpace(inputView.Buffer())
		inputView.Clear()
		inputView.SetCursor(0, 0)

		wordsView, _ := g.View("words view")
		wordsView.Clear()
		fmt.Fprint(wordsView, GetColorString("green", word))
		return nil
	})
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	close(done)
	return gocui.ErrQuit
}

func StartGoCui() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()
	g.SetManagerFunc(layout)
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, inputAction); err != nil {
		log.Panicln(err)
	}

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
				g.Update(func(g *gocui.Gui) error {
					wordsView, _ := g.View("words view")
					wordsView.Clear()
					wordsView.SetCursor(x, y)
					debugstr := fmt.Sprintf("(%d,%d)", x, y)
					fmt.Fprintln(wordsView, GetColorString("", "apple"+debugstr))
					y++
					maxx, maxy := g.Size()
					if x >= maxx {
						x = 1
					}
					if y >= maxy-15 {
						y = 1
					}

					return nil
				})
			}
		}
	}()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	wg.Wait()
}
