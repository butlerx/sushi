package gitissue

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Username string
	Token    string
}

type Comments struct {
	Issue []github.IssueComment
	Num   int
}

var (
	client          = github.NewClient(nil)
	conf            *Config
	folderpath, err = filepath.Abs(".")
	path            = folderpath + "/"
	GitLog          *log.Logger
)

func write(toWrite []github.Issue, file string) error {
	file = path + ".issue/" + file + ".json"
	b, err := json.Marshal(toWrite)
	if err == nil {
		err = ioutil.WriteFile(file, b, 0644)
	}
	return err
}

func IsSetUp() (bool, error) {
	_, err = os.Stat(path + ".issue")
	if os.IsNotExist(err) {
		return false, err
	}
	file, err := ioutil.ReadFile(path + ".issue/config.json")
	if err != nil {
		return false, err
	} else {
		temp := new(Config)
		if err = json.Unmarshal(file, temp); err != nil {
			return false, err
		}
		if temp.Username == "" {
			return false, err
		}
	}
	return true, nil
}

func SetUp(user, oauth string) error {
	GitLog = logSetUp()
	_, err := os.Stat(path + ".git")
	if os.IsNotExist(err) {
		return err
	}
	_, err = os.Stat(path + ".issue")
	if os.IsNotExist(err) {
		err := os.Mkdir(path+".issue", 0755)
		if err != nil {
			GitLog.Println("make folder: ", err)
			return err
		}
	}
	_, err = ioutil.ReadFile(path + ".issue/config.json")
	if err != nil {
		temp := Config{user, oauth}
		b, err := json.Marshal(temp)
		if err == nil {
			err = ioutil.WriteFile(path+".issue/config.json", b, 0644)
		} else {
			return err
		}
	}
	_, err = ioutil.ReadFile(path + ".issue/issues.json")
	if err != nil {
		temp := new([]github.Issue)
		b, err := json.Marshal(temp)
		if err == nil {
			err = ioutil.WriteFile(path+".issue/issues.json", b, 0644)
		} else {
			return err
		}
	}
	_, err = ioutil.ReadFile(path + ".issue/comments.json")
	if err != nil {
		temp := new([]Comments)
		b, err := json.Marshal(temp)
		if err == nil {
			err = ioutil.WriteFile(path+".issue/comments.json", b, 0644)
		} else {
			return err
		}
	}
	_, err = ioutil.ReadFile(path + ".gitignore")
	if err == nil {
		f, _ := os.OpenFile(".gitignore", os.O_APPEND, 0446)
		_, _ = f.WriteString(".issue/config.json")
		f.Close()
	} else {
		ignore := []byte(".issue/config.json")
		err = ioutil.WriteFile(".gitignore", ignore, 0644)
	}
	return err
}

func logSetUp() *log.Logger {
	logFile, err := os.OpenFile("sushi.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open logfile: ", err)
	}
	GitLog := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return GitLog
}

func Login() error {
	file, err := ioutil.ReadFile(path + ".issue/config.json")
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
	GitLog.Printf("\nLogged into: %v\n", github.Stringify(user.Login))
	return nil
}

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

// pull all issues and write them to .issue/issues.json
func Issues(repo string) ([]github.Issue, error) {
	issues, err := IssuesFilter(repo, "all", "", "", "", "", "", nil)
	if err == nil {
		write(issues, "issues")
		return issues, err
	}
	file, err := ioutil.ReadFile(path + ".issue/issues.json")
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

func MakeIssue(repo, title, body, assignee string, milestone int, labels []string) (*github.Issue, error) { // make issue put milestone at 0 for no milestone
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

func EditIssue(repo string, oldIssue github.Issue) (*github.Issue, error) { // make issue put milestone at 0 for no milestone
	s := strings.Split(repo, "/")
	issueNum := *oldIssue.Number
	issue := new(github.IssueRequest)
	var labels []string
	for i := 0; i < len(oldIssue.Labels); i++ {
		var label string
		label = oldIssue.Labels[i].String()
		labels = append(labels, label)
	}
	issue.Labels = &labels
	issue.Title = oldIssue.Title
	issue.Body = oldIssue.Body
	issue.Assignee = oldIssue.Assignee.Login
	issue.State = oldIssue.State
	updatedIssue, _, err := client.Issues.Edit(s[0], s[1], issueNum, issue)
	if err == nil {
		_, err = Issues(repo)
	}
	return updatedIssue, err
}

func CloseIssue(repo string, issue github.Issue) (*github.Issue, error) {
	temp := "closed"
	issue.State = &temp
	closedIssue, err := EditIssue(repo, issue)
	return closedIssue, err
}

func OpenIssue(repo string, issue github.Issue) (*github.Issue, error) {
	temp := "open"
	issue.State = &temp
	closedIssue, err := EditIssue(repo, issue)
	return closedIssue, err
}

func LockIssue(repo string, issueNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.Lock(s[0], s[1], issueNum)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

func UnlockIssue(repo string, issueNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.Unlock(s[0], s[1], issueNum)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

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

func storecomments(comments []github.IssueComment, issueNum int) error {
	toWrite, err := readComments()
	if err == nil {
		toAppend := Comments{comments, issueNum}
		toWrite = append(toWrite, toAppend)
		file := path + ".issue/comments.json"
		b, err := json.Marshal(toWrite)
		if err == nil {
			err = ioutil.WriteFile(file, b, 0644)
		}
	}
	return err
}

func readComments() ([]Comments, error) {
	file := path + ".issue/comments.json"
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

func Comment(repo, body string, issueNum int) (github.IssueComment, error) {
	s := strings.Split(repo, "/")
	comment := new(github.IssueComment)
	comment.Body = &body
	temp, _, err := client.Issues.CreateComment(s[0], s[1], issueNum, comment)
	newComment := *temp
	return newComment, err
}

func editComment(repo, body string, commentId int) (github.IssueComment, error) {
	s := strings.Split(repo, "/")
	comment := new(github.IssueComment)
	comment.Body = &body
	temp, _, err := client.Issues.EditComment(s[0], s[1], commentId, comment)
	newComment := *temp
	return newComment, err
}

func DeleteComment(repo string, commentId int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.DeleteComment(s[0], s[1], commentId)
	return err
}

func ListLabels(repo string) ([]github.Label, error) {
	s := strings.Split(repo, "/")
	labels, _, err := client.Issues.ListLabels(s[0], s[1], nil)
	return labels, err
}

func CreateLabel(repo, labelName string) (github.Label, error) {
	s := strings.Split(repo, "/")
	label := new(github.Label)
	label.Name = &labelName
	temp, _, err := client.Issues.CreateLabel(s[0], s[1], label)
	newLabel := *temp
	return newLabel, err
}

func AddLabel(repo, labelName string, issueNum int) ([]github.Label, error) {
	s := strings.Split(repo, "/")
	label := []string{labelName}
	labels, _, err := client.Issues.AddLabelsToIssue(s[0], s[1], issueNum, label)
	return labels, err
}

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

func RemoveLabel(repo, labelName string, issueNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.RemoveLabelForIssue(s[0], s[1], issueNum, labelName)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

func DeleteLabel(repo, labelName string) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.DeleteLabel(s[0], s[1], labelName)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

func CreateMilestone(repo, milestone string) (github.Milestone, error) {
	s := strings.Split(repo, "/")
	temp := new(github.Milestone)
	temp.Title = &milestone
	temp, _, err := client.Issues.CreateMilestone(s[0], s[1], temp)
	ms := *temp
	return ms, err
}

func AddMilestone() {}

func ListMilestones(repo string) ([]github.Milestone, error) {
	s := strings.Split(repo, "/")
	milestones, _, err := client.Issues.ListMilestones(s[0], s[1], nil)
	return milestones, err
}

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

func RemoveMilestone() {}
func DeleteMilestone(repo string, mileNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.DeleteMilestone(s[0], s[1], mileNum)
	return err
}

func Repos() ([]github.Repository, error) {
	repos, _, err := client.Repositories.List("", nil)
	return repos, err
}

func OrgsList() ([]github.Organization, error) {
	GitLog.Println(conf.Username)
	orgs, _, err := client.Organizations.List(conf.Username, nil)
	return orgs, err
}
