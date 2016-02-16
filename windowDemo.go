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

	browser.SetLayout(layout)
	if err := browser.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := browser.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("browser", 0, 0, maxX/2, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Browser")
	}
	if d, err := g.SetView("infopane", maxX/2, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(d, "Infopane")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
