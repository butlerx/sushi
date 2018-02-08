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

	"github.com/butlerx/sushi/gitissue"
	"github.com/google/go-github/github"
	"github.com/jroimartin/gocui"
)

var path = "./"
var issueList []*github.Issue
var comments [][]*github.IssueComment
var labelList []*github.Label

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
func getIssues() []*github.Issue {
	iss, err := gitissue.Issues(getRepo())
	if err != nil {
		log.Panicln(err)
	}
	return iss
}

//getComments returns an array of github issue comments from the repository returned by getRepo()
//length should be the length of the array returned by getIssues
func getComments(length int) [][]*github.IssueComment {
	var com = make([][]*github.IssueComment, length)
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
				temp, err := gitissue.OpenIssue(getRepo(), issueList[i])
				if err != nil {
					return err
				}
				issueList[i] = temp
				break
			} else {
				temp, err := gitissue.CloseIssue(getRepo(), issueList[i])
				if err != nil {
					return err
				}
				issueList[i] = temp
				break
			}
		}
	}
	if err := showIssues(g); err != nil {
		return err
	}
	return nil
}
