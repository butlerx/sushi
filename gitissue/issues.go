package gitissue

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

// Issues Pull all issues and write them to .issue/issues.json.
// Pulls both open and closed issues.
// Used to update issues.json.
// Method for accessing issue.json.
// Returns last pull if cant connect to server.
// Repo should be structured as "user/repo"
func Issues(repo string) ([]*github.Issue, error) {
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
	temp := new([]*github.Issue)
	if err = json.Unmarshal(file, temp); err != nil {
		GitLog.Println("parse issues: ", err)
		os.Exit(1)
	}
	issues = *temp
	return issues, err
}

// IssuesFilter Filters issues based on mileston, assignee, creatoror, labels or state.
// Pass empty strings for things that arnt to be filtered.
// Returns array of issues in order asked for.
// TODO(butlerx) filter offline issues if query of repo fails.
func IssuesFilter(repo, state, milestone, assignee, creator, sort, order string, labels []string) ([]*github.Issue, error) {
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
	ctx := context.Background()
	issues, _, err := client.Issues.ListByRepo(ctx, s[0], s[1], sorting)
	return issues, err
}

// MakeIssue Create an issue on github.
// Requires repo and title args,
// rest are optinal arg and can be passed empty.
// Make issue put milestone at 0 for no milestone.
// BUG(butlerx) Issue Creation doesnt work in offile mode.
func MakeIssue(repo, title, body, assignee string, milestone int, labels []string) (*github.Issue, error) {
	s := strings.Split(repo, "/")
	newIssue := new(github.Issue)
	state := "open"
	ctx := context.Background()
	if milestone == 0 {
		issue := new(github.IssueRequest)
		issue.Title = &title
		issue.Body = &body
		issue.Labels = &labels
		issue.Assignee = &assignee
		issue.State = &state
		newIssue, _, err = client.Issues.Create(ctx, s[0], s[1], issue)
	} else {
		var assignees []string
		issue := github.IssueRequest{&title, &body, &labels, &assignee, &state, &milestone, &assignees}
		newIssue, _, err = client.Issues.Create(ctx, s[0], s[1], &issue)
	}
	if err == nil {
		_, err = Issues(repo)
		return newIssue, err
	}
	return newIssue, err
}

// EditIssue Edit the issue object before passing it to this method.
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
	ctx := context.Background()
	updatedIssue, _, err := client.Issues.Edit(ctx, s[0], s[1], issueNum, issue)
	if err == nil {
		_, err = Issues(repo)
	} else {
		GitLog.Println("Edit issue: ", err)
	}
	return updatedIssue, err
}

// CloseIssue Marks issue as closed.
func CloseIssue(repo string, issue *github.Issue) (*github.Issue, error) {
	temp := "closed"
	issue.State = &temp
	closedIssue, err := EditIssue(repo, issue)
	return closedIssue, err
}

// OpenIssue Mark issue as open.
func OpenIssue(repo string, issue *github.Issue) (*github.Issue, error) {
	temp := "open"
	issue.State = &temp
	closedIssue, err := EditIssue(repo, issue)
	return closedIssue, err
}

// LockIssue Locks issue so it cant be changed.
func LockIssue(repo string, issueNum int) error {
	ctx := context.Background()
	s := strings.Split(repo, "/")
	_, err := client.Issues.Lock(ctx, s[0], s[1], issueNum)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

// UnlockIssue Unlocks issue to allow changes.
func UnlockIssue(repo string, issueNum int) error {
	s := strings.Split(repo, "/")
	ctx := context.Background()
	_, err := client.Issues.Unlock(ctx, s[0], s[1], issueNum)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

// writeIssue: Write issues out to file.
func writeIssue(toWrite []*github.Issue) error {
	file := *Path + ".issue/issues.json"
	b, err := json.Marshal(toWrite)
	if err == nil {
		err = ioutil.WriteFile(file, b, 0644)
	}
	return err
}
