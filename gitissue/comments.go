package gitissue

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

// Comments used for storing comments offline
// array of comments and issue number they relate to
type Comments struct {
	Issue []*github.IssueComment
	Num   int
}

// Comment Comment on a issue on github.
func Comment(repo, body string, issueNum int) (github.IssueComment, error) {
	s := strings.Split(repo, "/")
	ctx := context.Background()
	comment := new(github.IssueComment)
	comment.Body = &body
	temp, _, err := client.Issues.CreateComment(ctx, s[0], s[1], issueNum, comment)
	newComment := *temp
	return newComment, err
}

// storecomments: Write comments for issue to array and save it to file.
func storecomments(comments []*github.IssueComment, issueNum int) error {
	toWrite, err := readComments()
	if err == nil {
		toAppend := Comments{comments, issueNum}
		toWrite = append(toWrite, toAppend)
		file := *Path + ".gitissue/comments.json"
		b, err := json.Marshal(toWrite)
		if err == nil {
			err = ioutil.WriteFile(file, b, 0644)
		}
	}
	return err
}

// readComment Reads in comments from comments.json.
func readComments() ([]Comments, error) {
	file := *Path + ".gitissue/comments.json"
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

// EditComment Edit a comment already on github.
func EditComment(repo, body string, commentID int) (github.IssueComment, error) {
	s := strings.Split(repo, "/")
	comment := new(github.IssueComment)
	comment.Body = &body
	ctx := context.Background()
	temp, _, err := client.Issues.EditComment(ctx, s[0], s[1], commentID, comment)
	newComment := *temp
	return newComment, err
}

// DeleteComment Remove a comment from an issue.
func DeleteComment(repo string, commentID int) error {
	s := strings.Split(repo, "/")
	ctx := context.Background()
	_, err := client.Issues.DeleteComment(ctx, s[0], s[1], commentID)
	return err
}

// ListLabels list all possible labels in a repo.
func ListLabels(repo string) ([]*github.Label, error) {
	s := strings.Split(repo, "/")
	ctx := context.Background()
	labels, _, err := client.Issues.ListLabels(ctx, s[0], s[1], nil)
	return labels, err
}

// ListComments Lists all the comments for a given issue.
func ListComments(repo string, issueNum int) ([]*github.IssueComment, error) {
	s := strings.Split(repo, "/")
	err = nil
	ctx := context.Background()
	comments, _, err := client.Issues.ListComments(ctx, s[0], s[1], issueNum, nil)
	if err == nil {
		err = storecomments(comments, issueNum)
		return comments, err
	}
	commentStore, err := readComments()
	for i := 0; i < len(commentStore); i++ {
		if commentStore[i].Num == issueNum {
			comments = commentStore[i].Issue
		}
	}
	return comments, err
}
