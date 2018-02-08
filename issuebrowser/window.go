package issuebrowser

import "github.com/jroimartin/gocui"

//changeWindow is used to change the window in active focus
func changeWindow(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	if err := g.SetCurrentView("windowChanger"); err != nil {
		return err
	}
	return nil
}

func tabWindow(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	if err := g.SetCurrentView("windowTabber"); err != nil {
		return err
	}
	return nil
}

func windowUp(g *gocui.Gui, v *gocui.View) error {
	switch {
	case previousView == nil || previousView.Name() == "assigneepane":
		return g.SetCurrentView("milestonepane")
	case previousView.Name() == "browser":
		return g.SetCurrentView("browser")
	case previousView.Name() == "issuepane":
		return g.SetCurrentView("issuepane")
	case previousView.Name() == "commentpane":
		return g.SetCurrentView("issuepane")
	case previousView.Name() == "labelpane":
		return g.SetCurrentView("labelpane")
	case previousView.Name() == "milestonepane":
		return g.SetCurrentView("labelpane")
	default:
		return g.SetCurrentView("browser")
	}
}

func windowDown(g *gocui.Gui, v *gocui.View) error {
	switch {
	case previousView == nil || previousView.Name() == "assigneepane":
		return g.SetCurrentView("assigneepane")
	case previousView.Name() == "browser":
		return g.SetCurrentView("browser")
	case previousView.Name() == "issuepane":
		return g.SetCurrentView("commentpane")
	case previousView.Name() == "commentpane":
		return g.SetCurrentView("commentpane")
	case previousView.Name() == "labelpane":
		return g.SetCurrentView("milestonepane")
	case previousView.Name() == "milestonepane":
		return g.SetCurrentView("assigneepane")
	default:
		return g.SetCurrentView("browser")
	}
}

func windowRight(g *gocui.Gui, v *gocui.View) error {
	switch {
	case previousView == nil || previousView.Name() == "assigneepane":
		return g.SetCurrentView("assigneepane")
	case previousView.Name() == "browser":
		return g.SetCurrentView("issuepane")
	case previousView.Name() == "issuepane":
		return g.SetCurrentView("labelpane")
	case previousView.Name() == "commentpane":
		return g.SetCurrentView("milestonepane")
	case previousView.Name() == "labelpane":
		return g.SetCurrentView("labelpane")
	case previousView.Name() == "milestonepane":
		return g.SetCurrentView("milestonepane")
	default:
		return g.SetCurrentView("browser")
	}
}

func windowLeft(g *gocui.Gui, v *gocui.View) error {
	switch {
	case previousView == nil || previousView.Name() == "assigneepane":
		return g.SetCurrentView("commentpane")
	case previousView.Name() == "browser":
		return g.SetCurrentView("browser")
	case previousView.Name() == "issuepane":
		return g.SetCurrentView("browser")
	case previousView.Name() == "commentpane":
		return g.SetCurrentView("browser")
	case previousView.Name() == "labelpane":
		return g.SetCurrentView("issuepane")
	case previousView.Name() == "milestonepane":
		return g.SetCurrentView("commentpane")
	default:
		return g.SetCurrentView("browser")
	}
}

//nextWindow moves to the next window in sequence
func nextWindow(g *gocui.Gui, v *gocui.View) error {
	switch {
	case previousView == nil || previousView.Name() == "assigneepane":
		return g.SetCurrentView("browser")
	case previousView.Name() == "browser":
		return g.SetCurrentView("issuepane")
	case previousView.Name() == "issuepane":
		return g.SetCurrentView("commentpane")
	case previousView.Name() == "commentpane":
		return g.SetCurrentView("labelpane")
	case previousView.Name() == "labelpane":
		return g.SetCurrentView("milestonepane")
	case previousView.Name() == "milestonepane":
		return g.SetCurrentView("assigneepane")
	default:
		return g.SetCurrentView("browser")
	}
}

//previousWindow moves to the previous window in sequence
func previousWindow(g *gocui.Gui, v *gocui.View) error {
	switch {
	case previousView == nil || previousView.Name() == "assigneepane":
		return g.SetCurrentView("milestonepane")
	case previousView.Name() == "browser":
		return g.SetCurrentView("assigneepane")
	case previousView.Name() == "issuepane":
		return g.SetCurrentView("browser")
	case previousView.Name() == "commentpane":
		return g.SetCurrentView("issuepane")
	case previousView.Name() == "labelpane":
		return g.SetCurrentView("commentpane")
	case previousView.Name() == "milestonepane":
		return g.SetCurrentView("labelpane")
	default:
		return g.SetCurrentView("browser")
	}
}
