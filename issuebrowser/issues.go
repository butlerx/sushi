package issuebrowser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

//printIssues is used to print a specific issue Title, given it's index in the stored array
func printIssues(g *gocui.Gui, v *gocui.View, i int) {
	fmt.Fprint(v, *issueList[i].Number)
	fmt.Fprintln(v, ": "+(*issueList[i].Title))
}

//showIssues prints the list of issues to the browser window
func showIssues(g *gocui.Gui) error {
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	browser.Clear()
	if len(issueList) == 0 {
		fmt.Fprintln(browser, "Error reading issues or config file, please ensure your login and Oauth token are correct and that you have a stable internet connection")
	}
	for i := 0; i < len(issueList); i++ {
		if issueState {
			if *issueList[i].State == "open" {
				switch {
				case filterHeading == "Number":
					if strings.Contains(strconv.Itoa(*issueList[i].Number), filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "Title":
					if strings.Contains(*issueList[i].Title, filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "Body":
					if strings.Contains(*issueList[i].Body, filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "User":
					if strings.Contains(*issueList[i].User.Login, filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "Assignee":
					if strings.Contains(*issueList[i].Assignee.Login, filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "Comments":
					commentIndex := 0
					contains := false
					if *issueList[i].Comments > 0 {
						for ; commentIndex < len(comments); commentIndex++ {
							if len(comments[commentIndex]) > 0 {
								if *comments[commentIndex][0].IssueURL == *issueList[i].URL {
									break
								}
							}
						}
						for i := 0; i < len(comments[commentIndex]); i++ {
							if strings.Contains(*comments[commentIndex][i].Body, filterString) {
								contains = true
							}
						}
						if contains {
							printIssues(g, browser, i)
						}
					}
				case filterHeading == "Milestone Title":
					if strings.Contains(*issueList[i].Milestone.Title, filterString) {
						printIssues(g, browser, i)
					}
				default:
					printIssues(g, browser, i)
				}
			}
		} else {
			if *issueList[i].State == "closed" {
				switch {
				case filterHeading == "Number":
					if strings.Contains(strconv.Itoa(*issueList[i].Number), filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "Title":
					if strings.Contains(*issueList[i].Title, filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "Body":
					if strings.Contains(*issueList[i].Body, filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "User":
					if strings.Contains(*issueList[i].User.Login, filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "Assignee":
					if strings.Contains(*issueList[i].Assignee.Login, filterString) {
						printIssues(g, browser, i)
					}
				case filterHeading == "Comments":
					commentIndex := 0
					contains := false
					if *issueList[i].Comments > 0 {
						for ; commentIndex < len(comments); commentIndex++ {
							if len(comments[commentIndex]) > 0 && *comments[commentIndex][0].IssueURL == *issueList[i].URL {
								break
							}
						}
						for i := 0; i < len(comments[commentIndex]); i++ {
							if strings.Contains(*comments[commentIndex][i].Body, filterString) {
								contains = true
							}
						}
						if contains {
							printIssues(g, browser, i)
						}
					}
				case filterHeading == "Milestone Title":
					if strings.Contains(*issueList[i].Milestone.Title, filterString) {
						printIssues(g, browser, i)
					}
				default:
					printIssues(g, browser, i)
				}
			}
		}
	}
	return nil
}

//toggleIssues is used to swap between displaying open and closed issues
func toggleIssues(g *gocui.Gui, v *gocui.View) error {
	open, err := g.View("open")
	if err != nil {
		return err
	}
	closed, err := g.View("closed")
	if err != nil {
		return err
	}
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	if issueState {
		open.FgColor = gocui.ColorWhite
		open.BgColor = gocui.ColorBlack
		closed.FgColor = gocui.ColorBlack
		closed.BgColor = gocui.ColorWhite
		issueState = !issueState
	} else {
		open.FgColor = gocui.ColorBlack
		open.BgColor = gocui.ColorWhite
		closed.FgColor = gocui.ColorWhite
		closed.BgColor = gocui.ColorBlack
		issueState = !issueState
	}
	browser.Clear()
	for i := 0; i < len(issueList); i++ {
		if issueState {
			if *issueList[i].State == "open" {
				fmt.Fprint(browser, *issueList[i].Number)
				fmt.Fprintln(browser, ": "+(*issueList[i].Title))
			}
		} else {
			if *issueList[i].State == "closed" {
				fmt.Fprint(browser, *issueList[i].Number)
				fmt.Fprintln(browser, ": "+(*issueList[i].Title))
			}
		}
	}
	if err := g.SetCurrentView("browser"); err != nil {
		return err
	}
	if err := browser.SetCursor(0, 0); err != nil {
		return err
	}
	if err := getLine(g, browser); err != nil {
		return err
	}
	return nil
}

//newIssue opens the dialog box for generating a new issue
func newIssue(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	maxX, maxY := g.Size()
	if issueprompt, err := g.SetView("issueprompt", maxX/4, maxY/3, maxX-(maxX/4), (maxY/3)+(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(issueprompt, "Please enter issue title\n(Title is mandatory)\n\nCtrl + c to cancel")
	}
	if issueEd, err := g.SetView("issueEd", maxX/4, (maxY/3)+(maxY/6), maxX-(maxX/4), maxY-(maxY/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		issueEd.Editable = true
	}
	if err := g.SetCurrentView("issueEd"); err != nil {
		return err
	}
	return nil
}
