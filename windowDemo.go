package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jroimartin/gocui"
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
	if v, err := g.SetView("browser", -1, -1, maxX/3, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		arg1 := "./"
		if len(os.Args) != 1 {
			arg1 = os.Args[1]
		}
		files, _ := ioutil.ReadDir(arg1)
		for _, f := range files {
			fmt.Fprintln(v, f.Name())
		}
	}
	if d, err := g.SetView("infopane", maxX/3, -1, maxX, maxY); err != nil {
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
