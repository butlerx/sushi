package issuebrowser

import (
	"sort"

	"github.com/jroimartin/gocui"
)

//sortIssues sorts the list of issues depending on the sortChoice variable and then refreshes the display
func sortIssues(g *gocui.Gui, v *gocui.View) error {
	switch {
	case sortChoice == "Number":
		if orderChoice == "Ascending" {
			sort.Sort(byNumber(issueList))
		} else {
			sort.Sort(sort.Reverse(byNumber(issueList)))
		}
	case sortChoice == "Title":
		if orderChoice == "Ascending" {
			sort.Sort(byTitle(issueList))
		} else {
			sort.Sort(sort.Reverse(byTitle(issueList)))
		}
	case sortChoice == "Body":
		if orderChoice == "Ascending" {
			sort.Sort(byBody(issueList))
		} else {
			sort.Sort(sort.Reverse(byBody(issueList)))
		}
	case sortChoice == "User":
		if orderChoice == "Ascending" {
			sort.Sort(byUser(issueList))
		} else {
			sort.Sort(sort.Reverse(byUser(issueList)))
		}
	case sortChoice == "Assignee":
		if orderChoice == "Ascending" {
			sort.Sort(byAssignee(issueList))
		} else {
			sort.Sort(sort.Reverse(byAssignee(issueList)))
		}
	case sortChoice == "Comments":
		if orderChoice == "Ascending" {
			sort.Sort(byComments(issueList))
		} else {
			sort.Sort(sort.Reverse(byComments(issueList)))
		}
	case sortChoice == "Date Closed":
		if orderChoice == "Ascending" {
			sort.Sort(byClosedAt(issueList))
		} else {
			sort.Sort(sort.Reverse(byClosedAt(issueList)))
		}
	case sortChoice == "Date Created":
		if orderChoice == "Ascending" {
			sort.Sort(byCreatedAt(issueList))
		} else {
			sort.Sort(sort.Reverse(byCreatedAt(issueList)))
		}
	case sortChoice == "Date Updated":
		if orderChoice == "Ascending" {
			sort.Sort(byUpdatedAt(issueList))
		} else {
			sort.Sort(sort.Reverse(byUpdatedAt(issueList)))
		}
	case sortChoice == "Milestone Title":
		if orderChoice == "Ascending" {
			sort.Sort(byMilestone(issueList))
		} else {
			sort.Sort(sort.Reverse(byMilestone(issueList)))
		}
	case sortChoice == "":
		sort.Sort(byNumber(issueList))
	}
	if err := showIssues(g); err != nil {
		return err
	}
	if err := cancel(g, v); err != nil {
		return err
	}
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	if err := getLine(g, browser); err != nil {
		return err
	}
	return nil
}
