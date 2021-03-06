package issuebrowser

import (
	"fmt"
	"strings"

	"github.com/butlerx/sushi/gitissue"
	"github.com/jroimartin/gocui"
)

//nextEntry is used to cycle through each option in the newIssue creation process
func nextEntry(g *gocui.Gui, v *gocui.View) error {
	issueEd, err := g.View("issueEd")
	if err != nil {
		return err
	}
	issueprompt, err := g.View("issueprompt")
	if err != nil {
		return err
	}
	ox, oy := issueEd.Origin()
	if err := issueEd.SetCursor(ox, oy); err != nil {
		return err
	}
	switch {
	case entryCount == 0:
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
		if issueEd.Buffer() != "" {
			tempIssueBody = issueEd.Buffer()[:len(issueEd.Buffer())-2]
		}
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Please enter issue assignee\n(Leave blank for no assignee)\n\nCtrl + c to cancel")
		issueEd.Clear()
		entryCount++
	case entryCount == 2:
		if issueEd.Buffer() != "" {
			tempIssueAssignee = issueEd.Buffer()[:len(issueEd.Buffer())-2]
		}
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Please enter issue labels, label titles are comma separated\n(Leave blank for no labels)\n\nCtrl + c to cancel")
		issueEd.Clear()
		entryCount++
	case entryCount == 3:
		if issueEd.Buffer() != "" {
			tempIssueLabels = strings.Split(issueEd.Buffer()[:len(issueEd.Buffer())-2], ",")
		}
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Press enter to confirm entries and write out")
		fmt.Fprintln(issueprompt, "")
		fmt.Fprintln(issueprompt, "Press Ctrl + c to cancel")
		issueEd.Clear()
		issueEd.Editable = false
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
	case entryCount == 4:
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
