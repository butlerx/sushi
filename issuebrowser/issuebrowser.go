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
	"sort"
)

var path = "./"
var issueList []github.Issue
var comments [][]github.IssueComment

type byNumber []github.Issue
type byTitle []github.Issue
type byBody []github.Issue
type byUser []github.Issue
type byAssignee []github.Issue
type byComments []github.Issue
type byClosedAt []github.Issue
type byCreatedAt []github.Issue
type byUpdatedAt []github.Issue
type byMilestone []github.Issue

func (iss byNumber) Len() int {
	return len(iss)
}
func (iss byTitle) Len() int {
	return len(iss)
}
func (iss byBody) Len() int {
	return len(iss)
}
func (iss byUser) Len() int {
	return len(iss)
}
func (iss byAssignee) Len() int {
	return len(iss)
}
func (iss byComments) Len() int {
	return len(iss)
}
func (iss byClosedAt) Len() int {
	return len(iss)
}
func (iss byCreatedAt) Len() int {
	return len(iss)
}
func (iss byUpdatedAt) Len() int {
	return len(iss)
}
func (iss byMilestone) Len() int {
	return len(iss)
}

func (iss byNumber) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byTitle) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byBody) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byUser) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byAssignee) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byComments) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byClosedAt) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byCreatedAt) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byUpdatedAt) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byMilestone) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}

func (iss byNumber) Less(i, j int) bool {
	return *iss[i].Number < *iss[j].Number
}
func (iss byTitle) Less(i, j int) bool {
	return strings.Compare(*iss[i].Title, *iss[j].Title) < 0
}
func (iss byBody) Less(i, j int) bool {
	return strings.Compare(*iss[i].Body, *iss[j].Body) < 0
}
func (iss byUser) Less(i, j int) bool {
	return strings.Compare(*iss[i].User.Login, *iss[j].User.Login) < 0
}
func (iss byAssignee) Less(i, j int) bool {
	if iss[i].Assignee != nil && iss[j].Assignee != nil {
		return strings.Compare(*iss[i].Assignee.Login, *iss[j].Assignee.Login) < 0
	} else if iss[j].Assignee != nil {
		return false
	} else {
		return true
	}
}
func (iss byComments) Less(i, j int) bool {
	return *iss[i].Comments < *iss[j].Comments
}
func (iss byClosedAt) Less(i, j int) bool {
	if iss[i].ClosedAt != nil && iss[j].ClosedAt != nil {
		return (*iss[i].ClosedAt).Before(*iss[j].ClosedAt)
	} else if iss[j].ClosedAt != nil {
		return false
	} else {
		return true
	}
}
func (iss byCreatedAt) Less(i, j int) bool {
	return (*iss[i].CreatedAt).Before(*iss[j].CreatedAt)
}
func (iss byUpdatedAt) Less(i, j int) bool {
	return (*iss[i].UpdatedAt).Before(*iss[j].UpdatedAt)
}
func (iss byMilestone) Less(i, j int) bool {
	if iss[i].Milestone != nil && iss[j].Milestone != nil {
		return strings.Compare(*iss[i].Milestone.Title, *iss[j].Milestone.Title) < 0
	} else if iss[j].Milestone != nil {
		return false
	} else {
		return true
	}
}

var tempIssueTitle string
var tempIssueBody string
var tempIssueAssignee string
var tempIssueLabels = make([]string, 0)
var entryCount = 0
var issueState = true

var previousView *gocui.View
var sortChoice string
var orderChoice string

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
	if helppane, err := g.SetView("helppane", -1, maxY-2, maxX, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(helppane, "▲ ▼ ◀ ▶ = navigate, "+"\t"+"q = Quit, "+"\t"+"Ctrl+R = Refresh")
	}
	return nil
}

func keybindings(g *gocui.Gui) error {
	mainWindows := []string{"browser", "issuepane", "commentpane", "labelpane", "milestonepane", "assigneepane"}
	displayWindows := []string{"issuepane", "commentpane", "labelpane", "milestonepane", "assigneepane", "sortChoice"}
	controlWindows := []string{"browser", "issuepane", "commentpane", "labelpane", "milestonepane", "assigneepane", "sortChoice"}
	if err := g.SetKeybinding("browser", gocui.KeyF6, gocui.ModNone, showSortOrders); err != nil {
		return err
	}

	for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], gocui.KeyCtrlW, gocui.ModNone, changeWindow); err != nil {
			return err
		}
	}
	if err := g.SetKeybinding("issueEd", gocui.KeyEnter, gocui.ModNone, nextEntry); err != nil {
		return err
	}

	if err := g.SetKeybinding("issueEd", gocui.KeyCtrlC, gocui.ModNone, cancel); err != nil {
		return err
	}

	if err := g.SetKeybinding("sortChoice", gocui.KeyEnter, gocui.ModNone, getSortOrder); err != nil {
		return err
	}

	if err := g.SetKeybinding("sortChoice", gocui.KeyCtrlC, gocui.ModNone, cancel); err != nil {
		return err
	}
	for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], gocui.KeyCtrlN, gocui.ModNone, newIssue); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("browser", gocui.KeyPgup, gocui.ModNone, scrollupget); err != nil {
		return err
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], gocui.KeyPgup, gocui.ModNone, scrollup); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("browser", gocui.KeyPgdn, gocui.ModNone, scrolldownget); err != nil {
		return err
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], gocui.KeyPgdn, gocui.ModNone, scrolldown); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, scrollEnd); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyHome, gocui.ModNone, scrollHome); err != nil {
		return err
	}

	for i := 0; i < len(controlWindows); i++ {
		if err := g.SetKeybinding(controlWindows[i], '0', gocui.ModNone, scrollHome); err != nil {
			return err
		}
		if err := g.SetKeybinding(controlWindows[i], '$', gocui.ModNone, scrollEnd); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, scrollRight); err != nil {
		return err
	}

	for i := 0; i < len(controlWindows); i++ {
		if err := g.SetKeybinding(controlWindows[i], 'l', gocui.ModNone, scrollRight); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, scrollLeft); err != nil {
		return err
	}

	for i := 0; i < len(controlWindows); i++ {
		if err := g.SetKeybinding(controlWindows[i], 'h', gocui.ModNone, scrollLeft); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("browser", gocui.KeyArrowDown, gocui.ModNone, cursordownGet); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", 'j', gocui.ModNone, cursordownGet); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyArrowUp, gocui.ModNone, cursorupGet); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", 'k', gocui.ModNone, cursorupGet); err != nil {
		return err
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], gocui.KeyArrowDown, gocui.ModNone, cursordown); err != nil {
			return err
		}
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], 'j', gocui.ModNone, cursordown); err != nil {
			return err
		}
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], gocui.KeyArrowUp, gocui.ModNone, cursorup); err != nil {
			return err
		}
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], 'k', gocui.ModNone, cursorup); err != nil {
			return err
		}
	}

	for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], 'q', gocui.ModNone, quit); err != nil {
			return err
		}
	}

	for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], gocui.KeyTab, gocui.ModNone, toggleIssues); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, refresh); err != nil {
		return err
	}

	if err := g.SetKeybinding("windowChanger", gocui.KeyArrowUp, gocui.ModNone, windowUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", 'k', gocui.ModNone, windowUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", gocui.KeyArrowDown, gocui.ModNone, windowDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", 'j', gocui.ModNone, windowDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", gocui.KeyArrowRight, gocui.ModNone, windowRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", 'l', gocui.ModNone, windowRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", gocui.KeyArrowLeft, gocui.ModNone, windowLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", 'h', gocui.ModNone, windowLeft); err != nil {
		return err
	}

	for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], 'g', gocui.ModNone, tabWindow); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("windowTabber", 't', gocui.ModNone, nextWindow); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowTabber", 'T', gocui.ModNone, previousWindow); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowTabber", 'g', gocui.ModNone, scrollTop); err != nil {
		return err
	}
	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], 'G', gocui.ModNone, scrollBottom); err != nil {
			return err
		}
	}
	if err := g.SetKeybinding("browser", 'G', gocui.ModNone, scrollBottomGet); err != nil {
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

func showIssues(g *gocui.Gui) error {
	browser, err := g.View("browser")
	if err != nil {
		return err
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
	return nil
}

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

//functions called by keypress below
func showSortOrders(g *gocui.Gui, v *gocui.View) error {
	previousView = g.CurrentView()
	maxX, maxY := g.Size()
	if sortPrompt, err := g.SetView("sortPrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		sortPrompt.Wrap = true
		fmt.Fprintln(sortPrompt, "Please select a heading to sort by below")
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
	if err := showIssues(g); err != nil {
		return err
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

func scrollTop(g *gocui.Gui, v *gocui.View) error {
	if previousView.Name() == "browser" {
		err := scrollTopGet(g, v)
		if err != nil {
			return err
		}
	} else {
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
		if previousView != nil {
			if err := previousView.SetOrigin(0, 0); err != nil {
				return err
			}
			if err := previousView.SetCursor(0, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

func scrollBottom(g *gocui.Gui, v *gocui.View) error {
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
		for l != "" || m != "" {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
			cy++
			if l, err = v.Line(cy + 1); err != nil {
				l = ""
			}
			if m, err = v.Line(cy + 2); err != nil {
				m = ""
			}
		}
	}
	return nil
}

func scrollTopGet(g *gocui.Gui, v *gocui.View) error {
	if err := g.SetCurrentView(previousView.Name()); err != nil {
		return err
	}
	if previousView != nil {
		if err := previousView.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := previousView.SetCursor(0, 0); err != nil {
			return err
		}
	}
	if err := getLine(g, previousView); err != nil {
		return err
	}
	return nil
}

func scrollBottomGet(g *gocui.Gui, v *gocui.View) error {
	if err := scrollBottom(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func scrollup(g *gocui.Gui, v *gocui.View) error {
	_, maxY := v.Size()
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		for i := 0; i < maxY; i++ {
			if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
				if err := v.SetOrigin(ox, oy-1); err != nil {
					return err
				}
			}
			cy--
			oy--
		}
	}
	return nil
}

func scrolldown(g *gocui.Gui, v *gocui.View) error {
	_, maxY := v.Size()
	if v != nil {
		cx, cy := v.Cursor()
		var l string
		var m string
		var err error
		for i := 0; i < maxY; i++ {
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
			cy++
		}
	}
	return nil
}

func cursordown(g *gocui.Gui, v *gocui.View) error {
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

func scrollEnd(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		var l string
		var m string
		var err error
		line, err := v.Line(cy)
		if err != nil {
			return err
		}
		for i := 0; i < len(line); i++ {
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
				cx++
			}
		}
	}
	return nil
}

func scrollHome(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		line, err := v.Line(cy)
		if err != nil {
			return err
		}
		for i := 0; i < len(line); i++ {
			if err := v.SetCursor(cx-1, cy); err != nil && ox > 0 {
				if err := v.SetOrigin(ox-1, oy); err != nil {
					return err
				}
			}
			cx--
			ox--
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

func cursorup(g *gocui.Gui, v *gocui.View) error {
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

func scrollupget(g *gocui.Gui, v *gocui.View) error {
	if err := scrollup(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func scrolldownget(g *gocui.Gui, v *gocui.View) error {
	if err := scrolldown(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func cursordownGet(g *gocui.Gui, v *gocui.View) error {
	if err := cursordown(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func cursorupGet(g *gocui.Gui, v *gocui.View) error {
	if err := cursorup(g, v); err != nil {
		return err
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
					com = com[:len(com)-2]
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

func cancel(g *gocui.Gui, v *gocui.View) error {
	if (g.CurrentView()).Name() == "issueEd" {
		if err := g.DeleteView("issueEd"); err != nil {
			return err
		}
		if err := g.DeleteView("issueprompt"); err != nil {
			return err
		}
		entryCount = 0
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "sortChoice" {
		if err := g.DeleteView("sortChoice"); err != nil {
			return err
		}
		if err := g.DeleteView("sortPrompt"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	}
	return nil
}

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
