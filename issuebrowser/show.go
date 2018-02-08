package issuebrowser

import (
	"fmt"
	"log"

	"github.com/butlerx/sushi/gitissue"
	"github.com/jroimartin/gocui"
	"github.com/robfig/cron"
)

//Show is the main display function for the issue browser
func Show() {
	if err := setUp(); err != nil {
		log.Panic(err)
	}
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
	timer := cron.New()
	timer.AddFunc("0 5 * * * *", func() {
		reason, subject, update := gitissue.WatchRepo(getRepo())
		if update {
			fmt.Print("\a")
			notifications, err := window.View("notifications")
			if err != nil {
				log.Panic(err)
			}
			notifications.Clear()
			fmt.Fprintln(notifications, reason+": "+subject)
		}
	})
	timer.AddFunc("0 5 * * * *", func() {
		issueList = getIssues()
		browser, err := window.View("browser")
		if err != nil {
			log.Panic(err)
		}
		if err := sortIssues(window, browser); err != nil {
			log.Panic(err)
		}
		comments = getComments(len(issueList))
	})
	timer.Start()

	if err := window.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	timer.Stop()
}

//layout sets out the initial window layout for the program
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if windowChanger, err := g.SetView("windowChanger", maxX, maxY, maxX+1, maxY+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		windowChanger.Frame = false
	}
	if windowTabber, err := g.SetView("windowTabber", maxX, maxY, maxX+1, maxY+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		windowTabber.Frame = false
	}
	if open, err := g.SetView("open", -1, -1, maxX/6, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		open.FgColor = gocui.ColorBlack
		open.BgColor = gocui.ColorWhite
		fmt.Fprintln(open, "Open")
	}
	if closed, err := g.SetView("closed", maxX/6, -1, maxX/3, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		closed.FgColor = gocui.ColorWhite
		closed.BgColor = gocui.ColorBlack
		fmt.Fprintln(closed, "Closed")
	}
	if browser, err := g.SetView("browser", -1, 2, maxX/3, maxY); err != nil { //draw left pane
		if err != gocui.ErrUnknownView {
			return err
		}
		browser.Highlight = true
		if err := showIssues(g); err != nil {
			return err
		}
		if err := g.SetCurrentView("browser"); err != nil {
			return err
		}
	}
	if issuepane, err := g.SetView("issuepane", maxX/3, -1, maxX-(maxX/5), maxY/4); err != nil { //draw centre pane
		if err != gocui.ErrUnknownView {
			return err
		}
		issuepane.Wrap = true
	}
	if commentpane, err := g.SetView("commentpane", maxX/3, maxY/4, maxX-(maxX/5), maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentpane.Wrap = true
	}
	if labelpane, err := g.SetView("labelpane", maxX-(maxX/5), -1, maxX, maxY/3); err != nil { //draw labels pane
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(labelpane, "Labels")
	}
	if milestonepane, err := g.SetView("milestonepane", maxX-(maxX/5), maxY/3, maxX, maxY/3*2); err != nil { //draw milestone pane
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(milestonepane, "Milestone")
	}
	if assigneepane, err := g.SetView("assigneepane", maxX-(maxX/5), maxY/3*2, maxX, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(assigneepane, "Assignee")
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		if err := getLine(g, browser); err != nil {
			return err
		}
	}
	if infobar, err := g.SetView("infobar", -1, maxY-2, maxX/3, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(infobar, "F1/'?' = Help,\tq = Quit")
	}
	if _, err := g.SetView("notifications", maxX/3, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	return nil
}
