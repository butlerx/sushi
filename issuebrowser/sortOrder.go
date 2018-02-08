package issuebrowser

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

//showSortOrders displays a list of potential sorting options
func showSortOrders(g *gocui.Gui, v *gocui.View) error {
	previousView = g.CurrentView()
	maxX, maxY := g.Size()
	if sortPrompt, err := g.SetView("sortPrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sortPrompt.Wrap = true
		fmt.Fprintln(sortPrompt, "Please select a heading to sort by below")
		fmt.Fprintln(sortPrompt, "")
		fmt.Fprintln(sortPrompt, "Ctrl + C to cancel")
	}
	if sortChoice, err := g.SetView("sortChoice", maxX/4, maxY/3, maxX-(maxX/4), maxY-(maxY/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sortChoice.Highlight = true
		fmt.Fprintln(sortChoice, "Number")
		fmt.Fprintln(sortChoice, "Title")
		fmt.Fprintln(sortChoice, "Body")
		fmt.Fprintln(sortChoice, "User")
		fmt.Fprintln(sortChoice, "Assignee")
		fmt.Fprintln(sortChoice, "Comments")
		fmt.Fprintln(sortChoice, "Date Closed")
		fmt.Fprintln(sortChoice, "Date Created")
		fmt.Fprintln(sortChoice, "Date Updated")
		fmt.Fprintln(sortChoice, "Milestone Title")
		if err := g.SetCurrentView("sortChoice"); err != nil {
			return err
		}
	}
	return nil
}

//getSortOrder sets the sortOrder variable to the users choice
func getSortOrder(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	selection, err := v.Line(cy)
	if err != nil {
		return err
	}
	switch {
	case selection == "Ascending":
		orderChoice = selection
		v.Clear()
		if err := sortIssues(g, v); err != nil {
			return err
		}
	case selection == "Descending":
		orderChoice = selection
		v.Clear()
		if err := sortIssues(g, v); err != nil {
			return err
		}
	default:
		sortChoice = selection
		v.Clear()
		fmt.Fprintln(v, "Ascending")
		fmt.Fprintln(v, "Descending")
		if err := v.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := v.SetCursor(0, 0); err != nil {
			return err
		}
	}
	return nil
}
