package issuebrowser

import (
	"fmt"

	"github.com/butlerx/sushi/gitissue"
	"github.com/jroimartin/gocui"
)

//nextEntry is used to cycle through each option in the newIssue creation process
func nextEntry(g *gocui.Gui, v *gocui.View) error {
	onceThrough := false
	issueprompt, err := g.View("issueprompt")
	if err != nil {
		return err
	}
	switch {
	case entryCount == 0:
		issueEd, err := g.View("issueEd")
		if err != nil {
			return err
		}
		ox, oy := issueEd.Origin()
		if err := issueEd.SetCursor(ox, oy); err != nil {
			return err
		}
		if len(issueEd.Buffer()) >= 2 {
			tempIssueTitle = issueEd.Buffer()[:len(issueEd.Buffer())-2]
		}
		issueprompt.Clear()
		if tempIssueTitle == "" {
			fmt.Fprintln(issueprompt, "Please enter issue title\n(Title is mandatory)\n\nCtrl + c to cancel")
		} else {
			fmt.Fprintln(issueprompt, "Please enter issue body\n(Leave blank for no body)\n\nCtrl + c to cancel")
			issueEd.Clear()
			entryCount++
		}
	case entryCount == 1:
		issueEd, err := g.View("issueEd")
		if err != nil {
			return err
		}
		ox, oy := issueEd.Origin()
		if err := issueEd.SetCursor(ox, oy); err != nil {
			return err
		}
		if issueEd.Buffer() != "" {
			tempIssueBody = issueEd.Buffer()[:len(issueEd.Buffer())-2]
		}
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Please enter issue assignee\n(Leave blank for no assignee)\n\nCtrl + c to cancel")
		issueEd.Clear()
		entryCount++
	case entryCount == 2:
		if !onceThrough {
			issueEd, err := g.View("issueEd")
			if err != nil {
				return err
			}
			ox, oy := issueEd.Origin()
			if err := issueEd.SetCursor(ox, oy); err != nil {
				return err
			}
			if issueEd.Buffer() != "" {
				tempIssueAssignee = issueEd.Buffer()[:len(issueEd.Buffer())-2]
			}
			issueEd.Clear()
			if err := g.DeleteView("issueEd"); err != nil {
				return err
			}
		}
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Please choose a label(s) to assign to this issue.\n(Choose a blank line for no labels)\n\nCtrl + c to cancel")
		axX, maxY := g.Size()
		selectionPane, err := g.SetView("selectionPane", maxX/4, maxY/3, maxX/2, maxY-(maxY/6))
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			selectionPane.Highlight = true
			fmt.Fprintln(selectionPane, "\t")
			for i := 0; i < len(labelList); i++ {
				fmt.Fprintln(selectionPane, *labelList[i].Name)
			}
		}
		if _, err := g.SetView("selectionDisplay", maxX/2, maxY/3, maxX-(maxX/4), maxY-(maxY/6)); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
		}
		if err := g.SetCurrentView("selectionPane"); err != nil {
			return err
		}
		entryCount++
	case entryCount == 3:
		onceThrough = true
		selectionPane, err := g.View("selectionPane")
		if err != nil {
			return err
		}
		_, cy := selectionPane.Cursor()
		selection, err := selectionPane.Line(cy)
		if err != nil {
			return err
		}
		if err := g.SetCurrentView("selectionPane"); err != nil {
			return err
		}
		selectionDisplay, err := g.View("selectionDisplay")
		if err != nil {
			return err
		}
		if selection != "\t" {
			tempIssueLabels = append(tempIssueLabels, selection)
			for i := 0; i < len(tempIssueLabels); i++ {
				fmt.Fprintln(selectionDisplay, tempIssueLabels[i])
			}
			issueprompt.Clear()
			fmt.Fprintln(issueprompt, "Add another label?")
			selectionPane.Clear()
			fmt.Fprintln(selectionPane, "Yes")
			fmt.Fprintln(selectionPane, "No")
			selectionPane.SetOrigin(0, 0)
			selectionPane.SetCursor(0, 0)
			another, err := selectionPane.Line(0)
			if err != nil {
				return err
			}
			if another == "Yes" {
				entryCount--
			} else {
				entryCount++
			}
		} else {
			entryCount++
		}
	case entryCount == 4:
		selectionPane, err := g.View("selectionPane")
		if err != nil {
			return err
		}
		maxX, maxY := g.Size()
		if err := g.DeleteView("selectionDisplay"); err != nil {
			return err
		}
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Press enter to confirm entries and write out")
		fmt.Fprintln(issueprompt, "")
		fmt.Fprintln(issueprompt, "Press Ctrl + c to cancel")
		selectionPane.Clear()
		if err := g.DeleteView("selectionPane"); err != nil {
			return err
		}
		issueEd, err := g.SetView("issueEd", maxX/4, maxY/3, maxX-(maxX/4), maxY-(maxY/6))
		if err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			issueEd.Editable = false
		}
		if err := g.SetCurrentView("issueEd"); err != nil {
			return err
		}
		fmt.Fprintln(issueEd, "Title: "+tempIssueTitle)
		fmt.Fprintln(issueEd, "Body: "+tempIssueBody)
		fmt.Fprintln(issueEd, "Assignee: "+tempIssueAssignee)
		fmt.Fprint(issueEd, "Lables: ")
		for i := 0; i < len(tempIssueLabels); i++ {
			if i == 0 {
				fmt.Fprintln(issueEd, tempIssueLabels[i])
			} else {
				fmt.Fprintln(issueEd, "        "+tempIssueLabels[i])
			}
		}
		entryCount++
	case entryCount == 5:
		_, err := gitissue.MakeIssue(getRepo(), tempIssueTitle, tempIssueBody, tempIssueAssignee, 0, tempIssueLabels)
		if err != nil {
			return err
		}
		err = cancel(g, v)
		if err != nil {
			return err
		}
		tempIssueTitle = ""
		tempIssueBody = ""
		tempIssueAssignee = ""
		tempIssueLabels = make([]string, 0)
		refresh(g, v)
	default:
		fmt.Fprintln(issueprompt, "Error reading header")
	}
	return nil
}
