package issuebrowser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

//getLine parses the issue under the cursor in the browser window and displays the corresponding issue information in the relevant panels
func getLine(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error
	var index int
	var commentIndex int

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	issuepane, err := g.View("issuepane")
	if err != nil {
		return err
	}
	commentpane, err := g.View("commentpane")
	if err != nil {
		return err
	}
	labelpane, err := g.View("labelpane")
	if err != nil {
		return err
	}
	milestonepane, err := g.View("milestonepane")
	if err != nil {
		return err
	}
	assigneepane, err := g.View("assigneepane")
	if err != nil {
		return err
	}

	maxX, _ := g.Size()
	issuepane.Clear()
	commentpane.Clear()
	labelpane.Clear()
	milestonepane.Clear()
	assigneepane.Clear()

	fmt.Fprintln(commentpane, "Comments")
	fmt.Fprintln(commentpane, "")
	fmt.Fprintln(labelpane, "Labels")
	fmt.Fprintln(labelpane, "")
	fmt.Fprintln(milestonepane, "Milestone")
	fmt.Fprintln(milestonepane, "")
	fmt.Fprintln(assigneepane, "Assignee")
	fmt.Fprintln(assigneepane, "")
	if l != "" {
		//show issue body
		issNum := strings.Split(l, ":")
		for ; index < len(issueList); index++ {
			if (issNum[0]) == (strconv.Itoa(*issueList[index].Number)) {
				fmt.Fprintln(issuepane, *issueList[index].Title)
				fmt.Fprintln(issuepane, "")
				if *issueList[index].Body != "" {
					fmt.Fprintln(issuepane, *issueList[index].Body)
					fmt.Fprintln(issuepane, "")
				}
				fmt.Fprintln(issuepane, "#"+(strconv.Itoa(*issueList[index].Number))+" opened on "+((*issueList[index].CreatedAt).Format(time.UnixDate))+" by "+(*(*issueList[index].User).Login))
				break
			}
		}

		//show comments
		if *issueList[index].Comments > 0 {
			for ; commentIndex < len(comments); commentIndex++ {
				if len(comments[commentIndex]) > 0 {
					if *comments[commentIndex][0].IssueURL == *issueList[index].URL {
						break
					}
				}
			}
			for i := 0; i < (*issueList[index].Comments); i++ {
				fmt.Fprintln(commentpane, *comments[commentIndex][i].User.Login+" commented on "+(*comments[commentIndex][i].CreatedAt).Format("Mon Jan 2"))
				com := *comments[commentIndex][i].Body
				for strings.HasSuffix(com, "\n") {
					com = com[:len(com)-1]
				}
				fmt.Fprintln(commentpane, "\t\t\t\t"+com+"\n")
			}
		}

		//show labes
		labels := issueList[index].Labels
		if len(labels) == 0 {
			fmt.Fprintln(labelpane, "No Labels")
		}
		for i := 0; i < len(labels); i++ {
			fmt.Fprintln(labelpane, *labels[i].Name)
		}

		//show milestone
		if issueList[index].Milestone != nil {
			fmt.Fprintln(milestonepane, *issueList[index].Milestone.Title)
			complete := (float64(*issueList[index].Milestone.ClosedIssues) / (float64(*issueList[index].Milestone.OpenIssues) + float64(*issueList[index].Milestone.ClosedIssues)))
			barWidth := (maxX / 5) - 4
			bars := int(float64(barWidth) * complete)
			gaps := barWidth - bars
			fmt.Fprint(milestonepane, "[")
			for i := 0; i < bars; i++ {
				fmt.Fprint(milestonepane, "|")
			}
			for i := 0; i < gaps; i++ {
				fmt.Fprint(milestonepane, " ")
			}
			fmt.Fprintln(milestonepane, "]")
			complete = complete * 100
			fmt.Fprintln(milestonepane, (strconv.FormatFloat(complete, 'f', 0, 64))+"%")
		} else {
			fmt.Fprintln(milestonepane, "No Milestone")
		}

		//show assignee
		if issueList[index].Assignee != nil {
			fmt.Fprintln(assigneepane, *issueList[index].Assignee.Login)
		} else {
			fmt.Fprintln(assigneepane, "No Assignee")
		}
	} else {
		fmt.Fprintln(issuepane, "error")
		fmt.Fprintln(commentpane, "error")
		fmt.Fprintln(labelpane, "error")
		fmt.Fprintln(milestonepane, "error")
		fmt.Fprintln(assigneepane, "error")
	}
	return nil
}
