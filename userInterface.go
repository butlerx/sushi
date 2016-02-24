package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

func main() {
	window := gocui.NewGui()
	if err := window.Init(); err != nil {
		log.Panicln(err)
	}
	defer window.Close()

	window.SetLayout(layout)
	if err := keybindings(window); err != nil {
		log.Panicln(err)
	}
	window.SelBgColor = gocui.ColorWhite
	window.SelFgColor = gocui.ColorBlack
	window.Cursor = true

	if err := window.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	arg1 := "./"           //by default, read current dir
	if len(os.Args) != 1 { //check to see if you have cmd line args
		arg1 = os.Args[1]
	}
	files, _ := ioutil.ReadDir(arg1)

	if browser, err := g.SetView("browser", -1, -1, maxX/3, maxY); err != nil { //draw left pane
		if err != gocui.ErrUnknownView {
			return err
		}
		browser.Highlight = true

		for _, f := range files { //print file names
			fmt.Fprintln(browser, f.Name())
		}
		if err := g.SetCurrentView("browser"); err != nil {
			return err
		}
	}
	if infopane, err := g.SetView("infopane", maxX/3, -1, maxX, maxY); err != nil { //draw right pane
		if err != gocui.ErrUnknownView {
			return err
		}
		infopane.Wrap = true
		fmt.Fprintln(infopane, "")
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("browser", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyEnter, gocui.ModNone, getLine); err != nil {
		return err
	}

	return nil
}

//functions called by keypress below

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var arg1 string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	if len(os.Args) != 1 { //check to see if you have cmd line args
		arg1 = os.Args[1]
	}
	filePath := arg1 + l
	if err := g.DeleteView("infopane"); err != nil {
		return err
	}
	maxX, maxY := g.Size()
	if infopane, err := g.SetView("infopane", maxX/3, -1, maxX, maxY); err != nil { //draw right pane
		if err != gocui.ErrUnknownView {
			return err
		}
		infopane.Wrap = true
		dat, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		fmt.Fprintln(infopane, string(dat))
	}
	return nil
}
