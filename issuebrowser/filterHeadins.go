package issuebrowser

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

func showFilterHeadings(g *gocui.Gui, v *gocui.View) error {
	previousView = g.CurrentView()
	maxX, maxY := g.Size()
	if filterPrompt, err := g.SetView("filterPrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		filterPrompt.Wrap = true
		fmt.Fprintln(filterPrompt, "Please select a heading to filter by below")
		fmt.Fprintln(filterPrompt, "")
		fmt.Fprintln(filterPrompt, "Ctrl + C to cancel")
	}
	if filterChoice, err := g.SetView("filterChoice", maxX/4, maxY/3, maxX-(maxX/4), maxY-(maxY/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		filterChoice.Highlight = true
		fmt.Fprintln(filterChoice, "Clear Filter")
		fmt.Fprintln(filterChoice, "Number")
		fmt.Fprintln(filterChoice, "Title")
		fmt.Fprintln(filterChoice, "Body")
		fmt.Fprintln(filterChoice, "User")
		fmt.Fprintln(filterChoice, "Assignee")
		fmt.Fprintln(filterChoice, "Comments")
		fmt.Fprintln(filterChoice, "Milestone Title")
		if err := g.SetCurrentView("filterChoice"); err != nil {
			return err
		}
	}
	return nil
}

func getFilterHeading(g *gocui.Gui, v *gocui.View) error {
	filterChoice, err := g.View("filterChoice")
	if err != nil {
		return err
	}
	_, cy := filterChoice.Cursor()
	selection, err := filterChoice.Line(cy)
	if err != nil {
		return err
	}
	if selection == "Clear Filter" {
		filterHeading = ""
		filterString = ""
		if err := cancel(g, v); err != nil {
			return err
		}
	} else if filterHeading == "" {
		filterHeading = selection
		if strings.HasSuffix(filterHeading, "\n") {
			filterHeading = filterHeading[:len(filterHeading)-1]
		}
		filterChoice.Clear()
		if err := filterChoice.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := filterChoice.SetCursor(0, 0); err != nil {
			return err
		}
		filterChoice.Editable = true
	} else {
		filterString = selection
		if strings.HasSuffix(filterString, "\n") {
			filterString = filterString[:len(filterString)-1]
		}
		filterString = strings.Trim(filterString, " ")
		if err := cancel(g, v); err != nil {
			return err
		}
	}
	if err := showIssues(g); err != nil {
		return err
	}
	return nil
}
