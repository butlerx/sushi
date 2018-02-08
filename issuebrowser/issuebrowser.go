package issuebrowser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/butlerx/sushi/gitissue"

	"github.com/google/go-github/github"
	"github.com/jroimartin/gocui"
)

var path = "./"
var issueList []github.Issue
var comments [][]github.IssueComment
var labelList []github.Label

//typeCasting Issues so that each has a different sort method, allowing issues to be sorted by any heading
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

//used for recording user input for creating/editing issues
var tempIssueTitle string
var tempIssueBody string
var tempIssueAssignee string
var tempIssueLabels = make([]string, 0)

var entryCount = 0    //used to record which entry of a dialog the user is on
var issueState = true //true = open issues, false = closed issues

var changed bool
var previousView *gocui.View

//used to record sorting order
var sortChoice string
var orderChoice string

// PassArgs allows the calling program to pass a file path as a string
func PassArgs(s string) {
	path = s
}

//used to record filter order
var filterHeading string
var filterString string

//getRepo returns a string representation of the repository that may be required by other functions
func getRepo() string {
	dat, err := ioutil.ReadFile(*gitissue.Path + ".git/config")
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
				if strings.HasSuffix(ans, ".git") {
					ans = ans[:len(ans)-4]
				}
			} else {
				ans = ans[1:]
				if strings.HasSuffix(ans, ".git") {
					ans = ans[:len(ans)-4]
				}
			}
		}
	}
	return ans
}

//getIssues returns an array of github issues from the repository returned by getRepo()
func getIssues() []github.Issue {
	iss, err := gitissue.Issues(getRepo())
	if err != nil {
		log.Panicln(err)
	}
	return iss
}

//getComments returns an array of github issue comments from the repository returned by getRepo()
//length should be the length of the array returned by getIssues
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

//hide turns off command echoing on the terminal in order to hide user entry
func hide() {
	cmd := exec.Command("stty", "-echo")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

//unhide undoes the actions of hide
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
func GetLogin() error {
	isUp, err := gitissue.IsSetUp()
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error, sushi may only be called from inside a git repository")
		os.Exit(1)
	}
	if !isUp {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter GitHub username: ")
		user, _ := reader.ReadString('\n')
		if strings.HasSuffix(user, "\n") {
			user = user[:(len(user) - 1)]
		}
		fmt.Print("Enter GitHub Oauth token. Tokens can be generated at (https://github.com/settings/tokens) :")
		hide()
		pass, _ := reader.ReadString('\n')
		if strings.HasSuffix(pass, "\n") {
			pass = pass[:(len(pass) - 1)]
		}
		unhide()
		if err := gitissue.SetUp(user, pass, ""); err != nil {
			return err
		}
	} else {
		if err := gitissue.SetUp("", "", ""); err != nil {
			return err
		}
	}
	return nil
}

//setUp runs all of the necessary checks during startup
func setUp() error {
	if err := GetLogin(); err != nil {
		return err
	}
	if err := gitissue.Login(""); err != nil {
		return err
	}
	var err error
	issueList = getIssues()
	comments = getComments(len(issueList))
	labelList, err = gitissue.ListLabels(getRepo())
	if err != nil {
		return err
	}
	return nil
}

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
		}
	}
	return nil
}

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

//refresh queries the remote repository for an up to date version of the issues and comments and updates all local variables according
func refresh(g *gocui.Gui, v *gocui.View) error {
	var err error
	issueList = getIssues()
	comments = getComments(len(issueList))
	labelList, err = gitissue.ListLabels(getRepo())
	if err != nil {
		return err
	}
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

//quit exits the application
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

func cursordownGetIssues(g *gocui.Gui, v *gocui.View) error {
	if err := cursordown(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func cursorupGetIssues(g *gocui.Gui, v *gocui.View) error {
	if err := cursorup(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

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

//toggleState changes and open issue to closed and vice versa
func toggleState(g *gocui.Gui, v *gocui.View) error {
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	_, cy := browser.Cursor()
	current, err := browser.Line(cy)
	if err != nil {
		return err
	}
	line := strings.Split(current, ":")
	for i := 0; i < len(issueList); i++ {
		if line[0] == strconv.Itoa(*issueList[i].Number) {
			if *issueList[i].State != "open" {
				temp, err := gitissue.OpenIssue(getRepo(), &issueList[i])
				if err != nil {
					return err
				}
				issueList[i] = *temp
				break
			} else {
				temp, err := gitissue.CloseIssue(getRepo(), &issueList[i])
				if err != nil {
					return err
				}
				issueList[i] = *temp
				break
			}
		}
	}
	if err := showIssues(g); err != nil {
		return err
	}
	return nil
}

//newComment opens the new comment dialog box
func newComment(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	maxX, maxY := g.Size()
	if commentPrompt, err := g.SetView("commentPrompt", maxX/4, maxY/3, maxX-(maxX/4), (maxY/3)+(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(commentPrompt, "Please enter you comment text.\nPress enter to write out.\n\nPress Ctrl+C to cancel")
	}
	if commentBody, err := g.SetView("commentBody", maxX/4, (maxY/3)+(maxY/6), maxX-(maxX/4), maxY-(maxY/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentBody.Editable = true
	}
	if err := g.SetCurrentView("commentBody"); err != nil {
		return err
	}
	return nil
}

//writeComment writes out the users comment entered into the new comment dialog box
func writeComment(g *gocui.Gui, v *gocui.View) error {
	comment := v.Buffer()
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	_, cy := browser.Cursor()
	line, err := browser.Line(cy)
	if err != nil {
		return err
	}
	issueNum := strings.Split(line, ":")
	issueIndex := 0
	for ; issueIndex < len(issueList); issueIndex++ {
		if issueNum[0] == strconv.Itoa(*issueList[issueIndex].Number) {
			break
		}
	}
	commentIndex := 0
	for ; commentIndex < len(comments); commentIndex++ {
		if len(comments[commentIndex]) > 0 {
			if *comments[commentIndex][0].IssueURL == *issueList[issueIndex].URL {
				break
			}
		}
	}
	num, err := strconv.Atoi(issueNum[0])
	if err != nil {
		return err
	}
	_, err = gitissue.Comment(getRepo(), comment, num)
	if err != nil {
		return err
	}
	if err := cancel(g, v); err != nil {
		return err
	}
	refresh(g, v)
	return nil
}

//openCommentEditor opens a list of comments for editing
func openCommentEditor(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	issueIndex := 0
	commentIndex := 0
	maxX, maxY := g.Size()
	if commentEditPrompt, err := g.SetView("commentEditPrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(commentEditPrompt, "Select the comment you wish to edit\n\nCtrl+C to cancel")
	}
	if commentBrowser, err := g.SetView("commentBrowser", maxX/4, maxY/3, maxX/2, maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentBrowser.Highlight = true
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		_, cy := browser.Cursor()
		issueLine, err := browser.Line(cy)
		if err != nil {
			return err
		}
		issueNum := strings.Split(issueLine, ":")
		var URL string
		for ; issueIndex < len(issueList); issueIndex++ {
			if issueNum[0] == strconv.Itoa(*issueList[issueIndex].Number) {
				URL = *issueList[issueIndex].URL
				break
			}
		}
		if *issueList[issueIndex].Comments > 0 {
			for ; commentIndex < len(comments); commentIndex++ {
				if len(comments[commentIndex]) > 0 {
					if URL == *comments[commentIndex][0].IssueURL {
						break
					}
				}
			}
			for i := 0; i < len(comments[commentIndex]); i++ {
				fmt.Fprintln(commentBrowser, strconv.Itoa(*comments[commentIndex][i].ID)+": "+*comments[commentIndex][i].User.Login+"@"+(*comments[commentIndex][i].CreatedAt).Format(time.UnixDate))
			}
		} else {
			fmt.Fprintln(commentBrowser, "This issue has no comments")
		}
		if err := g.SetCurrentView("commentBrowser"); err != nil {
			return err
		}
	}
	if commentViewer, err := g.SetView("commentViewer", maxX/2, maxY/3, maxX-(maxX/4), maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentViewer.Wrap = true
		if *issueList[issueIndex].Comments > 0 {
			fmt.Fprintln(commentViewer, *comments[commentIndex][0].Body)
		}
	}
	return nil
}

//editComment opens the individual comment chosed in openCommentEditor
func editComment(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	commentLine, err := v.Line(cy)
	if err != nil {
		return err
	}
	if commentLine == "This issue has no comments" {
		if err := cancel(g, v); err != nil {
			return err
		}
		return nil
	}
	commentEditPrompt, err := g.View("commentEditPrompt")
	if err != nil {
		return err
	}
	commentEditPrompt.Clear()
	fmt.Fprintln(commentEditPrompt, "Press enter to write out changes\n\nCtrl+C to cancel")
	if err := g.SetCurrentView("commentViewer"); err != nil {
		return err
	}
	commentViewer, err := g.View("commentViewer")
	if err != nil {
		return err
	}
	commentViewer.Editable = true
	return nil
}

//writeEditedComment writes out changes made to a comment
func writeEditedComment(g *gocui.Gui, v *gocui.View) error {
	commentBrowser, err := g.View("commentBrowser")
	if err != nil {
		return err
	}
	_, cy := commentBrowser.Cursor()
	commentLine, err := commentBrowser.Line(cy)
	if err != nil {
		return err
	}
	ID, err := strconv.Atoi((strings.Split(commentLine, ":"))[0])
	if err != nil {
		return err
	}
	if _, err = gitissue.EditComment(getRepo(), v.Buffer(), ID); err != nil {
		return err
	}
	v.Editable = false
	if err = cancel(g, v); err != nil {
		return err
	}
	if err = refresh(g, v); err != nil {
		return err
	}
	return nil
}

//openCommentDeleter opens a list of comments for deletion
func openCommentDeleter(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	issueIndex := 0
	commentIndex := 0
	maxX, maxY := g.Size()
	if commentDeletePrompt, err := g.SetView("commentDeletePrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(commentDeletePrompt, "Select the comment you wish to delete\n\nCtrl+C to cancel")
	}
	if commentDeleter, err := g.SetView("commentDeleter", maxX/4, maxY/3, maxX/2, maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentDeleter.Highlight = true
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		_, cy := browser.Cursor()
		issueLine, err := browser.Line(cy)
		if err != nil {
			return err
		}
		issueNum := strings.Split(issueLine, ":")
		var URL string
		for ; issueIndex < len(issueList); issueIndex++ {
			if issueNum[0] == strconv.Itoa(*issueList[issueIndex].Number) {
				URL = *issueList[issueIndex].URL
				break
			}
		}
		if *issueList[issueIndex].Comments > 0 {
			for ; commentIndex < len(comments); commentIndex++ {
				if len(comments[commentIndex]) > 0 {
					if URL == *comments[commentIndex][0].IssueURL {
						break
					}
				}
			}
			for i := 0; i < len(comments[commentIndex]); i++ {
				fmt.Fprintln(commentDeleter, strconv.Itoa(*comments[commentIndex][i].ID)+": "+*comments[commentIndex][i].User.Login+"@"+(*comments[commentIndex][i].CreatedAt).Format(time.UnixDate))
			}
		} else {
			fmt.Fprintln(commentDeleter, "This issue has no comments")
		}
		if err := g.SetCurrentView("commentDeleter"); err != nil {
			return err
		}
	}
	if commentViewer, err := g.SetView("commentViewer", maxX/2, maxY/3, maxX-(maxX/4), maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentViewer.Wrap = true
		if *issueList[issueIndex].Comments > 0 {
			fmt.Fprintln(commentViewer, *comments[commentIndex][0].Body)
		}
	}
	return nil
}

//deleteComment deletes a comment chosen in openCommentDeleter
func deleteComment(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	commentLine, err := v.Line(cy)
	if err != nil {
		return err
	}
	if commentLine == "This issue has no comments" {
		if err := cancel(g, v); err != nil {
			return err
		}
		return nil
	}
	ID, err := strconv.Atoi((strings.Split(commentLine, ":"))[0])
	if err != nil {
		return err
	}
	if err = gitissue.DeleteComment(getRepo(), ID); err != nil {
		return err
	}
	if err = cancel(g, v); err != nil {
		return err
	}
	if err = refresh(g, v); err != nil {
		return err
	}
	return nil
}

func cursorupGetComments(g *gocui.Gui, v *gocui.View) error {
	if err := cursorup(g, v); err != nil {
		return err
	}
	if err := getLineComment(g, v); err != nil {
		return err
	}
	return nil
}

func cursordownGetComments(g *gocui.Gui, v *gocui.View) error {
	if err := cursordown(g, v); err != nil {
		return err
	}
	if err := getLineComment(g, v); err != nil {
		return err
	}
	return nil
}

func getLineComment(g *gocui.Gui, v *gocui.View) error {
	issueIndex := 0
	commentIndex := 0
	_, cy := v.Cursor()
	commentLine, err := v.Line(cy)
	if err != nil {
		return err
	}
	if commentLine == "This issue has no comments" {
		return nil
	}
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	_, cy = browser.Cursor()
	issueLine, err := browser.Line(cy)
	if err != nil {
		return err
	}
	issueNum := strings.Split(issueLine, ":")
	var URL string
	for ; issueIndex < len(issueList); issueIndex++ {
		if issueNum[0] == strconv.Itoa(*issueList[issueIndex].Number) {
			URL = *issueList[issueIndex].URL
			break
		}
	}
	for ; commentIndex < len(comments); commentIndex++ {
		if len(comments[commentIndex]) > 0 {
			if URL == *comments[commentIndex][0].IssueURL {
				break
			}
		}
	}
	ID := strings.Split(commentLine, ":")
	IDnum, err := strconv.Atoi(ID[0])
	if err != nil {
		return err
	}
	commentViewer, err := g.View("commentViewer")
	if err != nil {
		return err
	}
	commentViewer.Clear()
	for i := 0; i < len(comments[commentIndex]); i++ {
		if IDnum == *comments[commentIndex][i].ID {
			fmt.Fprintln(commentViewer, *comments[commentIndex][i].Body)
			break
		}
	}
	return nil
}

func addLabel(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	maxX, maxY := g.Size()
	if labelPropmt, err := g.SetView("labelPrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelPropmt.Wrap = true
		fmt.Fprintln(labelPropmt, "Please choose from the list of labels below.\nPress enter to add a label to the issue.\n\nCtrl+C to exit")
	}
	if labelBrowser, err := g.SetView("labelBrowser", maxX/4, maxY/3, maxX/2, maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelBrowser.Highlight = true
		for i := 0; i < len(labelList); i++ {
			fmt.Fprintln(labelBrowser, *labelList[i].Name)
		}
		if err := g.SetCurrentView("labelBrowser"); err != nil {
			return err
		}
	}

	if labelViewer, err := g.SetView("labelViewer", maxX/2, maxY/3, maxX-(maxX/4), maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		_, cy := browser.Cursor()
		issueLine, err := browser.Line(cy)
		if err != nil {
			return err
		}
		issueNum := (strings.Split(issueLine, ":"))[0]
		for i := 0; i < len(issueList); i++ {
			if issueNum == strconv.Itoa(*issueList[i].Number) {
				for j := 0; j < len(issueList[i].Labels); j++ {
					fmt.Fprintln(labelViewer, *issueList[i].Labels[j].Name)
				}
				break
			}
		}
	}
	return nil
}

func openLabelRemover(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	maxX, maxY := g.Size()
	if labelPropmt, err := g.SetView("labelPrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelPropmt.Wrap = true
		fmt.Fprintln(labelPropmt, "Please choose from the list of labels below.\nPress enter to remove a label from an issue.\n\nCtrl+C to exit")
	}
	if labelRemover, err := g.SetView("labelRemover", maxX/4, maxY/3, maxX-(maxX/4), maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelRemover.Highlight = true
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		_, cy := browser.Cursor()
		issueLine, err := browser.Line(cy)
		if err != nil {
			return err
		}
		issueNum := (strings.Split(issueLine, ":")[0])
		for i := 0; i < len(issueList); i++ {
			if issueNum == strconv.Itoa(*issueList[i].Number) {
				for j := 0; j < len(issueList[i].Labels); j++ {
					fmt.Fprintln(labelRemover, *issueList[i].Labels[j].Name)
				}
				break
			}
		}
		if err := g.SetCurrentView("labelRemover"); err != nil {
			return err
		}
	}
	return nil
}

func removeLabel(g *gocui.Gui, v *gocui.View) error {
	labelRemover, err := g.View("labelRemover")
	if err != nil {
		return err
	}
	_, cy := labelRemover.Cursor()
	label, err := labelRemover.Line(cy)
	if err != nil {
		return err
	}
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	_, cy = browser.Cursor()
	issueLine, err := browser.Line(cy)
	if err != nil {
		return err
	}
	issueNum, err := strconv.Atoi((strings.Split(issueLine, ":"))[0])
	if err != nil {
		return err
	}
	if err := gitissue.RemoveLabel(getRepo(), label, issueNum); err != nil {
		return err
	}
	if err := cancel(g, v); err != nil {
		return err
	}
	if err := refresh(g, v); err != nil {
		return err
	}
	return nil
}

func writeLabel(g *gocui.Gui, v *gocui.View) error {
	changed = true
	labelBrowser, err := g.View("labelBrowser")
	if err != nil {
		return err
	}
	_, cy := labelBrowser.Cursor()
	label, err := labelBrowser.Line(cy)
	if err != nil {
		return err
	}
	labelViewer, err := g.View("labelViewer")
	if err != nil {
		return err
	}
	fmt.Fprintln(labelViewer, label)
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	_, cy = browser.Cursor()
	issueLine, err := browser.Line(cy)
	if err != nil {
		return err
	}
	issueNum, err := strconv.Atoi((strings.Split(issueLine, ":"))[0])
	if err != nil {
		return err
	}
	_, err = gitissue.AddLabel(getRepo(), label, issueNum)
	if err != nil {
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

//cancel is used to close dialog boxes
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
	} else if (g.CurrentView()).Name() == "filterChoice" {
		if err := g.DeleteView("filterChoice"); err != nil {
			return err
		}
		if err := g.DeleteView("filterPrompt"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "commentBody" {
		if err := g.DeleteView("commentBody"); err != nil {
			return err
		}
		if err := g.DeleteView("commentPrompt"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "commentBrowser" || (g.CurrentView()).Name() == "commentViewer" {
		if err := g.DeleteView("commentEditPrompt"); err != nil {
			return err
		}
		if err := g.DeleteView("commentBrowser"); err != nil {
			return err
		}
		if err := g.DeleteView("commentViewer"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "commentDeleter" {
		if err := g.DeleteView("commentDeletePrompt"); err != nil {
			return err
		}
		if err := g.DeleteView("commentDeleter"); err != nil {
			return err
		}
		if err := g.DeleteView("commentViewer"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "labelBrowser" {
		if err := g.DeleteView("labelPrompt"); err != nil {
			return err
		}
		if err := g.DeleteView("labelBrowser"); err != nil {
			return err
		}
		if err := g.DeleteView("labelViewer"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
		if changed == true {
			changed = false
			if err := refresh(g, v); err != nil {
				return err
			}
		}
	} else if (g.CurrentView()).Name() == "labelRemover" {
		if err := g.DeleteView("labelRemover"); err != nil {
			return err
		}
		if err := g.DeleteView("labelPrompt"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "helpPane" {
		if err := g.DeleteView("helpPane"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	}
	return nil
}

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
