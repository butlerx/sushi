package startscreen

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var path = "./"

// PassArgs allows the calling program to pass a file path as a string
func PassArgs(s string) {
	path = s
}

//Show is the main display function for the start screen
func Show() {
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

	if browser, err := g.SetView("browser", -1, -1, maxX/3, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		browser.Highlight = true
		fmt.Fprintln(browser, "First Time Startup")
		fmt.Fprintln(browser, "Issue Browser")
		fmt.Fprintln(browser, "Pull Request Manager")
		if err := g.SetCurrentView("browser"); err != nil {
			return err
		}
	}

	if infopane, err := g.SetView("infopane", maxX/3, -1, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		infopane.Wrap = true
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		if err := getLine(g, browser); err != nil {
			return err
		}
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

	if err := g.SetKeybinding("infopane", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("infopane", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	return nil
}

//functions called by keypress below

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
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

func scrollUp(g *gocui.Gui, v *gocui.View) error {
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
	if err := getLine(g, v); err != nil {
		return err
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
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	if err := g.DeleteView("infopane"); err != nil {
		return err
	}

	maxX, maxY := g.Size()
	if infopane, err := g.SetView("infopane", maxX/3, -1, maxX, maxY); err != nil { //draw right pane
		if err != gocui.ErrUnknownView {
			return err
		}
		infopane.Wrap = true
		if l == "First Time Startup" {
			fmt.Fprintln(infopane, "Initializes sushi for the first time in a repo, creating a .issue folder and storing github issues locally within the repository.")
		} else if l == "Issue Browser" {
			fmt.Fprintln(infopane, "A browser for viewing, managing and editing github issues.")
		} else if l == "Pull Request Manager" {
			fmt.Fprintln(infopane, "A browser for viewing and managing pull requests.")
		} else {
			fmt.Fprintln(infopane, "Error")
		}
	}
	return nil
}
