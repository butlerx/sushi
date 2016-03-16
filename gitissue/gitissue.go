package gitissue

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Structure of User config file
// Username and oauth token stored
type Config struct {
	Username string
	Token    string
}

// Structure for comments to file
// used for storing comments offline
// array of comments and issue number they relate to
type Comments struct {
	Issue []github.IssueComment
	Num   int
}

var (
	client          = github.NewClient(nil)
	conf            *Config
	folderpath, err = filepath.Abs(".")
	folder          = (folderpath + "/")
	Path            = &folder
	GitLog          *log.Logger
)

// Write issues out to file.
func writeIssue(toWrite []github.Issue) error {
	file := *Path + ".issue/issues.json"
	b, err := json.Marshal(toWrite)
	if err == nil {
		err = ioutil.WriteFile(file, b, 0644)
	}
	return err
}

// Check if users config is set up and if being run in git repo.
func IsSetUp() (bool, error) {
	if checkgit() == false {
		err := errors.New("Not git repo")
		return false, err
	}
	_, err = os.Stat(*Path + ".issue")
	if os.IsNotExist(err) {
		return false, nil
	}
	file, err := ioutil.ReadFile(*Path + ".issue/config.json")
	if err != nil {
		return false, nil
	} else {
		temp := new(Config)
		if err = json.Unmarshal(file, temp); err != nil {
			return false, nil
		}
		if temp.Username == "" {
			return false, nil
		}
	}
	return true, nil
}

// Check if being run in a git repo or a child in a git repo.
func checkgit() bool {
	_, err := os.Stat(*Path + ".git")
	if os.IsNotExist(err) {
		if *Path == "/" {
			return false
		} else {
			s := strings.Split(*Path, "/")
			*Path = "/"
			for i := 1; i < len(s)-2; i++ {
				*Path = *Path + s[i] + "/"
			}
			checkgit()
		}
	}
	return true
}

// Set up .issue folder, issues, comments & config file and begin logfile.
// Appends to to gitignore to ignore the config file.
// Checks if the files exist and if they dont creates them.
func SetUp(user, oauth string) error {
	GitLog = logSetUp()
	_, err = os.Stat(*Path + ".issue")
	if os.IsNotExist(err) {
		err := os.Mkdir(*Path+".issue", 0755)
		if err != nil {
			GitLog.Println("make folder: ", err)
			return err
		}
	}
	_, err = ioutil.ReadFile(*Path + ".issue/config.json")
	if err != nil {
		temp := Config{user, oauth}
		b, err := json.Marshal(temp)
		if err == nil {
			err = ioutil.WriteFile(*Path+".issue/config.json", b, 0644)
		} else {
			return err
		}
	}
	_, err = ioutil.ReadFile(*Path + ".issue/issues.json")
	if err != nil {
		temp := new([]github.Issue)
		b, err := json.Marshal(temp)
		if err == nil {
			err = ioutil.WriteFile(*Path+".issue/issues.json", b, 0644)
		} else {
			return err
		}
	}
	_, err = ioutil.ReadFile(*Path + ".issue/comments.json")
	if err != nil {
		temp := new([]Comments)
		b, err := json.Marshal(temp)
		if err == nil {
			err = ioutil.WriteFile(*Path+".issue/comments.json", b, 0644)
		} else {
			return err
		}
	}
	_, err = ioutil.ReadFile(*Path + ".gitignore")
	if err == nil {
		f, _ := os.OpenFile(".gitignore", os.O_APPEND, 0446)
		_, _ = f.WriteString(".issue/config.json")
		f.Close()
	} else {
		ignore := []byte(".issue/config.json")
		err = ioutil.WriteFile(".gitignore", ignore, 0644)
		return err
	}
	return nil
}

// Sets up Log file and creates logger object.
// Returns GitLog to be used to log erros.
func logSetUp() *log.Logger {
	_, err = ioutil.ReadFile(*Path + ".issue/sushi.log")
	logFile := new(os.File)
	if err != nil {
		logFile, err = os.Create(*Path + ".issue/sushi.log")
	} else {
		logFile, err = os.OpenFile(*Path+".issue/sushi.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}
	if err != nil {
		log.Fatalln("Failed to open logfile: ", err)
	}
	GitLog := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return GitLog
}

// Logs in to github using oauth.
// Returns error if login fails.
func Login() error {
	file, err := ioutil.ReadFile(*Path + ".issue/config.json")
	if err != nil {
		GitLog.Println("open config: ", err)
		os.Exit(1)
	}
	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		GitLog.Println("parse config: ", err)
		os.Exit(1)
	}
	conf = temp
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = github.NewClient(tc)
	user, _, err := client.Users.Get("")

	if err != nil {
		GitLog.Printf("\nerror: %v\n", err)
		return err
	}
	log.Printf("\nLogged into: %v\n", github.Stringify(user.Login))
	GitLog.Print("Logged into: ", github.Stringify(user.Login))
	return nil
}

// Filters issues based on mileston, assignee, creatoror, labels or state.
// Pass empty strings for things that arnt to be filtered.
// Returns array of issues in order asked for.
// TODO(butlerx) filter offline issues if query of repo fails.
func IssuesFilter(repo, state, milestone, assignee, creator, sort, order string, labels []string) ([]github.Issue, error) {
	s := strings.Split(repo, "/")
	sorting := new(github.IssueListByRepoOptions)
	if len(labels) != 0 {
		sorting.Labels = labels
	}
	if state != "" {
		sorting.State = state
	}
	if milestone != "" {
		sorting.Milestone = milestone
	}
	if assignee != "" {
		sorting.Assignee = assignee
	}
	if creator != "" {
		sorting.Creator = creator
	}
	if sort != "" {
		sorting.Sort = sort
	}
	if order != "" {
		sorting.Direction = order
	}
	issues, _, err := client.Issues.ListByRepo(s[0], s[1], sorting)
	return issues, err
}

// Pull all issues and write them to .issue/issues.json.
// Pulls both open and closed issues.
// Used to update issues.json.
// Method for accessing issue.json.
// Returns last pull if cant connect to server.
// Repo should be structured as "user/repo"
func Issues(repo string) ([]github.Issue, error) {
	issues, err := IssuesFilter(repo, "all", "", "", "", "", "", nil)
	if err == nil {
		writeIssue(issues)
		return issues, err
	}
	file, err := ioutil.ReadFile(*Path + ".issue/issues.json")
	if err != nil {
		GitLog.Println("open issues: ", err)
		os.Exit(1)
	}
	temp := new([]github.Issue)
	if err = json.Unmarshal(file, temp); err != nil {
		GitLog.Println("parse issues: ", err)
		os.Exit(1)
	}
	issues = *temp
	return issues, err
}

// Create an issue on github.
// Requires repo and title args,
// rest are optinal arg and can be passed empty.
// Make issue put milestone at 0 for no milestone.
// BUG(butlerx) doesnt work in offile mode.
func MakeIssue(repo, title, body, assignee string, milestone int, labels []string) (*github.Issue, error) {
	s := strings.Split(repo, "/")
	newIssue := new(github.Issue)
	state := "open"
	if milestone == 0 {
		issue := new(github.IssueRequest)
		issue.Title = &title
		issue.Body = &body
		issue.Labels = &labels
		issue.Assignee = &assignee
		issue.State = &state
		newIssue, _, err = client.Issues.Create(s[0], s[1], issue)
	} else {
		issue := github.IssueRequest{&title, &body, &labels, &assignee, &state, &milestone}
		newIssue, _, err = client.Issues.Create(s[0], s[1], &issue)
	}
	if err == nil {
		_, err = Issues(repo)
		return newIssue, err
	} else {
		return newIssue, err
	}
}

// Edit a github issue.
// Edit the issue object before passing it to this method.
func EditIssue(repo string, oldIssue *github.Issue) (*github.Issue, error) {
	s := strings.Split(repo, "/")
	issueNum := *oldIssue.Number
	issue := new(github.IssueRequest)
	if len(oldIssue.Labels) != 0 {
		labels := []string{oldIssue.Labels[0].String()}
		for i := 1; i < len(oldIssue.Labels); i++ {
			var label string
			label = oldIssue.Labels[i].String()
			labels = append(labels, label)
		}
		issue.Labels = &labels
	}
	issue.Title = oldIssue.Title
	issue.Body = oldIssue.Body
	if oldIssue.Assignee != nil {
		issue.Assignee = oldIssue.Assignee.Login
	}
	issue.State = oldIssue.State
	updatedIssue, _, err := client.Issues.Edit(s[0], s[1], issueNum, issue)
	if err == nil {
		_, err = Issues(repo)
	} else {
		GitLog.Println("Edit issue: ", err)
	}
	return updatedIssue, err
}

// Marks issue as closed.
func CloseIssue(repo string, issue *github.Issue) (*github.Issue, error) {
	temp := "closed"
	issue.State = &temp
	closedIssue, err := EditIssue(repo, issue)
	return closedIssue, err
}

// Mark issue as open.
func OpenIssue(repo string, issue *github.Issue) (*github.Issue, error) {
	temp := "open"
	issue.State = &temp
	closedIssue, err := EditIssue(repo, issue)
	return closedIssue, err
}

// Locks issue so it cant be changed.
func LockIssue(repo string, issueNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.Lock(s[0], s[1], issueNum)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

// Unlocks issue to allow changes.
func UnlockIssue(repo string, issueNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.Unlock(s[0], s[1], issueNum)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

// Lists all the comments for a given issue.
func ListComments(repo string, issueNum int) ([]github.IssueComment, error) {
	s := strings.Split(repo, "/")
	err = nil
	comments, _, err := client.Issues.ListComments(s[0], s[1], issueNum, nil)
	if err == nil {
		err = storecomments(comments, issueNum)
		return comments, err
	} else {
		commentStore, err := readComments()
		for i := 0; i < len(commentStore); i++ {
			if commentStore[i].Num == issueNum {
				comments = commentStore[i].Issue
			}
		}
		return comments, err
	}
}

// Write comments for issue to array and save it to file.
func storecomments(comments []github.IssueComment, issueNum int) error {
	toWrite, err := readComments()
	if err == nil {
		toAppend := Comments{comments, issueNum}
		toWrite = append(toWrite, toAppend)
		file := *Path + ".issue/comments.json"
		b, err := json.Marshal(toWrite)
		if err == nil {
			err = ioutil.WriteFile(file, b, 0644)
		}
	}
	return err
}

// Reads in comments from comments.json.
func readComments() ([]Comments, error) {
	file := *Path + ".issue/comments.json"
	read, err := ioutil.ReadFile(file)
	if err != nil {
		GitLog.Println("open comments: ", err)
		os.Exit(1)
	}
	temp := new([]Comments)
	if err = json.Unmarshal(read, temp); err != nil {
		GitLog.Println("parse comments: ", err)
		os.Exit(1)
	}
	comments := *temp
	return comments, err
}

// Comment on a issue on github.
func Comment(repo, body string, issueNum int) (github.IssueComment, error) {
	s := strings.Split(repo, "/")
	comment := new(github.IssueComment)
	comment.Body = &body
	temp, _, err := client.Issues.CreateComment(s[0], s[1], issueNum, comment)
	newComment := *temp
	return newComment, err
}

// Edit a comment already on github.
func editComment(repo, body string, commentId int) (github.IssueComment, error) {
	s := strings.Split(repo, "/")
	comment := new(github.IssueComment)
	comment.Body = &body
	temp, _, err := client.Issues.EditComment(s[0], s[1], commentId, comment)
	newComment := *temp
	return newComment, err
}

// Remove a comment from an issue.
func DeleteComment(repo string, commentId int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.DeleteComment(s[0], s[1], commentId)
	return err
}

// list all possible labels in a repo.
func ListLabels(repo string) ([]github.Label, error) {
	s := strings.Split(repo, "/")
	labels, _, err := client.Issues.ListLabels(s[0], s[1], nil)
	return labels, err
}

// Create a label for a repo.
func CreateLabel(repo, labelName string) (github.Label, error) {
	s := strings.Split(repo, "/")
	label := new(github.Label)
	label.Name = &labelName
	temp, _, err := client.Issues.CreateLabel(s[0], s[1], label)
	newLabel := *temp
	return newLabel, err
}

// Add a label to a repo.
func AddLabel(repo, labelName string, issueNum int) ([]github.Label, error) {
	s := strings.Split(repo, "/")
	label := []string{labelName}
	labels, _, err := client.Issues.AddLabelsToIssue(s[0], s[1], issueNum, label)
	return labels, err
}

// Change the name of a label.
func EditLabel(repo, labelName, newName string) (github.Label, error) {
	s := strings.Split(repo, "/")
	label := new(github.Label)
	label.Name = &newName
	temp, _, err := client.Issues.EditLabel(s[0], s[1], labelName, label)
	editedLabel := *temp
	if err == nil {
		_, err = Issues(repo)
	}
	return editedLabel, err
}

// Remove a label from an issue.
func RemoveLabel(repo, labelName string, issueNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.RemoveLabelForIssue(s[0], s[1], issueNum, labelName)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

// Delete a label from a repo.
func DeleteLabel(repo, labelName string) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.DeleteLabel(s[0], s[1], labelName)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

// Create a Milestone for a repo.
func CreateMilestone(repo, milestone string) (github.Milestone, error) {
	s := strings.Split(repo, "/")
	temp := new(github.Milestone)
	temp.Title = &milestone
	temp, _, err := client.Issues.CreateMilestone(s[0], s[1], temp)
	ms := *temp
	return ms, err
}

// Add Milestone to an issue.
// BUG(butlerx) currently not supported as milestones in the api are a mix of strings and ints.
func AddMilestone() {}

// List all Milestones in a repo.
func ListMilestones(repo string) ([]github.Milestone, error) {
	s := strings.Split(repo, "/")
	milestones, _, err := client.Issues.ListMilestones(s[0], s[1], nil)
	return milestones, err
}

// Change the title of a milestone in a repo.
func EditMilestone(repo, newTitle string, mileNum int) (github.Milestone, error) {
	s := strings.Split(repo, "/")
	temp, _, err := client.Issues.GetMilestone(s[0], s[1], mileNum)
	if err != nil {
		milestone := *temp
		return milestone, err
	}
	temp.Title = &newTitle
	newmilestone, _, err := client.Issues.EditMilestone(s[0], s[1], mileNum, temp)
	if err == nil {
		_, err = Issues(repo)
	}
	milestone := *newmilestone
	return milestone, err
}

// Remove Milestone to an issue.
// BUG(butlerx) currently not supported as milestones in the api are a mix of strings and ints.
func RemoveMilestone() {}

// Delet a Milestone from a repo.
func DeleteMilestone(repo string, mileNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.DeleteMilestone(s[0], s[1], mileNum)
	return err
}

// monitor for change in repo

// list all all of a users repos.
// currently unused.
func repos() ([]github.Repository, error) {
	repos, _, err := client.Repositories.List("", nil)
	return repos, err
}

// List all orgs a users a part of.
// currently unused.
func orgsList() ([]github.Organization, error) {
	GitLog.Println(conf.Username)
	orgs, _, err := client.Organizations.List(conf.Username, nil)
	return orgs, err
}
