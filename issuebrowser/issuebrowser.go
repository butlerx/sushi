package issuebrowser

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/butlerx/AgileGit/gitissue"
	"github.com/google/go-github/github"
	"github.com/jroimartin/gocui"
)

var path = "./"
var issueList = getIssues()

func getRepo() string {
	dat, err := ioutil.ReadFile(".git/config")
	if err != nil {
		panic(err)
	}
	list := strings.Split(string(dat), "\n")
	ans := ""
	for i := 0; i < len(list); i++ {
		if strings.Contains(list[i], "github.com") {
			sublist := strings.Split(list[i], "github.com")
			ans = sublist[len(sublist)-1]
		}
	}
	return ans
}

func getIssues() []github.Issue {
	iss, err := gitissue.Issues(getRepo())
	if err != nil {
		log.Panicln(err)
	}
	return iss
}

// PassArgs allows the calling program to pass a file path as a string
func PassArgs(s string) {
	path = s
}

//Show is the main display function for the issue browser
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
	if browser, err := g.SetView("browser", -1, -1, maxX/3, maxY); err != nil { //draw left pane
		if err != gocui.ErrUnknownView {
			return err
		}
		browser.Highlight = true
		for i := 0; i < len(issueList); i++ {
			fmt.Fprint(browser, *issueList[i].Number)
			fmt.Fprintln(browser, ": "+(*issueList[i].Title))
		}
		if err := g.SetCurrentView("browser"); err != nil {
			return err
		}
	}
	if issuepane, err := g.SetView("issuepane", maxX/3, -1, maxX-(maxX/5), maxY); err != nil { //draw centre pane
		if err != gocui.ErrUnknownView {
			return err
		}
		issuepane.Wrap = true
	}
	if labelpane, err := g.SetView("labelpane", maxX-(maxX/5), -1, maxX, maxY/3); err != nil { //draw labels pane
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(labelpane, "Labels")
	}
	if milestonepane, err := g.SetView("milestonepane", maxX-(maxX/5), maxY/3, maxX, maxY); err != nil { //draw milestone pane
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(milestonepane, "Milestone")
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

	if err := g.SetKeybinding("issuepane", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyArrowRight, gocui.ModNone, moveRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", gocui.KeyArrowLeft, gocui.ModNone, moveLeft); err != nil {
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
	var index int

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	if err := g.DeleteView("issuepane"); err != nil {
		return err
	}
	if err := g.DeleteView("labelpane"); err != nil {
		return err
	}
	if err := g.DeleteView("milestonepane"); err != nil {
		return err
	}

	maxX, maxY := g.Size()
	if issuepane, err := g.SetView("issuepane", maxX/3, -1, maxX-(maxX/5), maxY); err != nil { //draw centre pane
		if err != gocui.ErrUnknownView {
			return err
		}
		issuepane.Wrap = true
		if l != "" {
			issNum := strings.Split(l, ":")
			for ; index < len(issueList); index++ {
				if (issNum[0]) == (strconv.Itoa(*issueList[index].Number)) {
					fmt.Fprintln(issuepane, *issueList[index].Title)
					fmt.Fprintln(issuepane, "")
					fmt.Fprintln(issuepane, "#"+(strconv.Itoa(*issueList[index].Number))+" opened on "+((*issueList[index].CreatedAt).Format(time.UnixDate))+" by "+(*(*issueList[index].User).Login))
					break
				}
			}
		} else {
			fmt.Fprintln(issuepane, "error")
		}
	}

	if labelpane, err := g.SetView("labelpane", maxX-(maxX/5), -1, maxX, maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelpane.Wrap = true
		labels := issueList[index].Labels
		if l != "" {
			fmt.Fprintln(labelpane, "Labels")
			fmt.Fprintln(labelpane, "")
			for i := 0; i < len(labels); i++ {
				fmt.Fprintln(labelpane, *labels[i].Name)
			}
		} else {
			fmt.Fprintln(labelpane, "error")
		}
	}
	if milestonepane, err := g.SetView("milestonepane", maxX-(maxX/5), maxY/3, maxX, maxY); err != nil { //draw milestone pane
		if err != gocui.ErrUnknownView {
			return err
		}
		if l != "" {
			fmt.Fprintln(milestonepane, "Milestone")
			fmt.Fprintln(milestonepane, "")
			if issueList[index].Milestone != nil {
				complete := ((*issueList[index].Milestone.ClosedIssues)/(*issueList[index].Milestone.OpenIssues) + (*issueList[index].Milestone.ClosedIssues)) * 100
				fmt.Fprintln(milestonepane, (strconv.Itoa(complete))+"%")
			} else {
				fmt.Fprintln(milestonepane, "No Milestone")
			}
		} else {
			fmt.Fprintln(milestonepane, "error")
		}
	}
	return nil
}

func moveRight(g *gocui.Gui, v *gocui.View) error {
	if err := g.SetCurrentView("issuepane"); err != nil {
		return err
	}
	return nil
}

func moveLeft(g *gocui.Gui, v *gocui.View) error {
	if err := g.SetCurrentView("browser"); err != nil {
		return err
	}
	return nil
}
