package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

func main() {
	browser := gocui.NewGui()
	if err := browser.Init(); err != nil {
		log.Panicln(err)
	}
	defer browser.Close()

	info := gocui.NewGui()
	if err := info.Init(); err != nil {
		log.Panicln(err)
	}
	defer info.Close()

	browser.SetLayout(browserLayout)
	if err := browser.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	info.SetLayout(infoLayout)
	if err := info.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := browser.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	if err := info.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func browserLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", 0, 0, maxX/2, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func infoLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("googbye", maxX/2, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Goodbye world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
