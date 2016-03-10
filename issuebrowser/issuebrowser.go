package issuebrowser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/butlerx/AgileGit/gitissue"
	"github.com/google/go-github/github"
	"github.com/jroimartin/gocui"
	"github.com/robfig/cron"
)

var path = "./"
var issueList []github.Issue
var comments [][]github.IssueComment
var tempIssueTitle string
var tempIssueBody string
var tempIssueAssignee string
var tempIssueLabels = make([]string, 0)
var entryCount = 0

// PassArgs allows the calling program to pass a file path as a string
func PassArgs(s string) {
	path = s
}

//Show is the main display function for the issue browser
func Show() {
	setUp()
	timer := cron.New()
	timer.AddFunc("0 5 * * * *", func() { issueList = getIssues() })
	timer.AddFunc("0 5 * * * *", func() { comments = getComments(len(issueList)) })
	timer.Start()
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
	timer.Stop()
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
	if helppane, err := g.SetView("helppane", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(helppane, "▲ ▼ ◀ ▶ = navigate, "+"\t"+"q = Quit, "+"\t"+"Ctrl+R = Refresh")
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("issueEd", gocui.KeyEnter, gocui.ModNone, nextEntry); err != nil {
		return err
	}

	if err := g.SetKeybinding("issueEd", gocui.KeyCtrlC, gocui.ModNone, cancel); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlN, gocui.ModNone, newIssue); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", 'l', gocui.ModNone, scrollRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", 'l', gocui.ModNone, scrollRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", 'l', gocui.ModNone, scrollRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", 'l', gocui.ModNone, scrollRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("milestonepane", 'l', gocui.ModNone, scrollRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("assigneepane", 'l', gocui.ModNone, scrollRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", 'h', gocui.ModNone, scrollLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", 'h', gocui.ModNone, scrollLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", 'h', gocui.ModNone, scrollLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", 'h', gocui.ModNone, scrollLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding("milestonepane", 'h', gocui.ModNone, scrollLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding("assigneepane", 'h', gocui.ModNone, scrollLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, scrollRight); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, scrollLeft); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", 'j', gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", 'k', gocui.ModNone, cursorUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", 'j', gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", 'k', gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", 'j', gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", 'k', gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", 'j', gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", 'k', gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("milestonepane", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("milestonepane", 'j', gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("milestonepane", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("milestonepane", 'k', gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("assigneepane", gocui.KeyArrowDown, gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("assigneepane", 'j', gocui.ModNone, scrollDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("assigneepane", gocui.KeyArrowUp, gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("assigneepane", 'k', gocui.ModNone, scrollUp); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("milestonepane", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("assigneepane", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyTab, gocui.ModNone, nextWindow); err != nil {
		return err
	}

	if err := g.SetKeybinding("issuepane", gocui.KeyTab, gocui.ModNone, nextWindow); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", gocui.KeyTab, gocui.ModNone, nextWindow); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", gocui.KeyTab, gocui.ModNone, nextWindow); err != nil {
		return err
	}

	if err := g.SetKeybinding("milestonepane", gocui.KeyTab, gocui.ModNone, nextWindow); err != nil {
		return err
	}

	if err := g.SetKeybinding("assigneepane", gocui.KeyTab, gocui.ModNone, nextWindow); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, refresh); err != nil {
		return err
	}

	return nil
}

//helper functions
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
			if strings.HasPrefix(ans, ":") {
				ans = ans[1:(len(ans) - 4)]
			} else {
				ans = ans[1:]
			}
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

func getComments(length int) [][]github.IssueComment {
	var com = make([][]github.IssueComment, length)
	var err error
	for i := 0; i < len(issueList); i++ {
		com[i], err = gitissue.ListComments(getRepo(), (*issueList[i].Number))
		if err != nil {
			log.Panic(err)
		}
	}
	return com
}

func hide() {
	cmd := exec.Command("stty", "-echo")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func unhide() {
	cmd := exec.Command("stty", "echo")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

//GetLogin sets up the config file if it does not exist
func GetLogin() {
	if _, err := os.Stat(".issue/config.json"); os.IsNotExist(err) {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter GitHub username: ")
		user, _ := reader.ReadString('\n')
		if strings.HasSuffix(user, "\n") {
			user = user[:(len(user) - 2)]
		}
		fmt.Print("Enter GitHub Oauth token: ")
		hide()
		pass, _ := reader.ReadString('\n')
		if strings.HasSuffix(pass, "\n") {
			pass = pass[:(len(pass) - 2)]
		}
		unhide()
		gitissue.SetUp(user, pass)
	}
}

func setUp() {
	GetLogin()
	gitissue.Login()
	issueList = getIssues()
	comments = getComments(len(issueList))
}

//functions called by keypress below

func refresh(g *gocui.Gui, v *gocui.View) error {
	issueList = getIssues()
	comments = getComments(len(issueList))
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	current := g.CurrentView()
	if err := g.SetCurrentView("browser"); err != nil {
		return err
	}
	browser.Clear()
	for i := 0; i < len(issueList); i++ {
		fmt.Fprint(browser, *issueList[i].Number)
		fmt.Fprintln(browser, ": "+(*issueList[i].Title))
	}
	if err := getLine(g, browser); err != nil {
		return err
	}
	if err := g.SetCurrentView(current.Name()); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func scrollDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		var l string
		var m string
		var err error
		if l, err = v.Line(cy + 1); err != nil {
			l = ""
		}
		if m, err = v.Line(cy + 2); err != nil {
			m = ""
		}
		if l != "" || m != "" {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func scrollRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		var l string
		var m string
		var err error
		if l, err = v.Word(cx, cy); err != nil {
			l = ""
		}
		if m, err = v.Word(cx+1, cy); err != nil {
			m = ""
		}
		if l != "" || m != "" {
			ox, oy := v.Origin()
			if err := v.SetCursor(cx+1, cy); err != nil {
				if err := v.SetOrigin(ox+1, oy); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func scrollLeft(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if err := v.SetCursor(cx-1, cy); err != nil && ox > 0 {
			if err := v.SetOrigin(ox-1, oy); err != nil {
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
		var l string
		var err error
		if l, err = v.Line(cy + 1); err != nil {
			l = ""
		}
		if l != "" {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
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
		for i := 0; i < (*issueList[index].Comments); i++ {
			fmt.Fprintln(commentpane, *comments[index][i].User.Login+" commented on "+(*comments[index][i].CreatedAt).Format("Mon Jan 2"))
			com := *comments[index][i].Body
			for strings.HasSuffix(com, "\n") {
				com = com[:len(com)-2]
			}
			fmt.Fprintln(commentpane, "\t\t\t\t"+com+"\n")
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

func newIssue(g *gocui.Gui, v *gocui.View) error {
	maxX, maxY := g.Size()
	if issueprompt, err := g.SetView("issueprompt", maxX/4, maxY/3, maxX-(maxX/4), (maxY/3)+(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(issueprompt, "Please enter issue title\n\n\nCtrl + c to cancel")
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
		tempIssueTitle = issueEd.Buffer()[:len(issueEd.Buffer())-2]
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Please enter issue body\n(Blank for no body)\n\nCtrl + c to cancel")
		issueEd.Clear()
		entryCount++
	case entryCount == 1:
		tempIssueBody = issueEd.Buffer()[:len(issueEd.Buffer())-2]
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Please enter issue assignee\n(Blank for no assignee)\n\nCtrl + c to cancel")
		issueEd.Clear()
		entryCount++
	case entryCount == 2:
		tempIssueAssignee = issueEd.Buffer()[:len(issueEd.Buffer())-2]
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Please enter issue labels, label titles are comma separated\n(Blank for no labels)\n\nCtrl + c to cancel")
		issueEd.Clear()
		entryCount++
	case entryCount == 3:
		tempIssueLabels = strings.Split(issueEd.Buffer()[:len(issueEd.Buffer())-2], ",")
		issueprompt.Clear()
		fmt.Fprintln(issueprompt, "Press enter to confirm entries and write out")
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
			// return err
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

func cancel(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("issueEd"); err != nil {
		return err
	}
	if err := g.DeleteView("issueprompt"); err != nil {
		return err
	}
	entryCount = 0
	if err := g.SetCurrentView("browser"); err != nil {
		return err
	}
	return nil
}

func nextWindow(g *gocui.Gui, v *gocui.View) error {
	switch {
	case v == nil || v.Name() == "assigneepane":
		return g.SetCurrentView("browser")
	case v.Name() == "browser":
		return g.SetCurrentView("issuepane")
	case v.Name() == "issuepane":
		return g.SetCurrentView("commentpane")
	case v.Name() == "commentpane":
		return g.SetCurrentView("labelpane")
	case v.Name() == "labelpane":
		return g.SetCurrentView("milestonepane")
	case v.Name() == "milestonepane":
		return g.SetCurrentView("assigneepane")
	default:
		return g.SetCurrentView("browser")
	}
}
