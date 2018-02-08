package gitissue

import (
	"context"
	"strings"

	"github.com/google/go-github/github"
)

// CreateLabel Create a label for a repo.
func CreateLabel(repo, labelName string) (github.Label, error) {
	s := strings.Split(repo, "/")
	label := new(github.Label)
	label.Name = &labelName
	ctx := context.Background()
	temp, _, err := client.Issues.CreateLabel(ctx, s[0], s[1], label)
	newLabel := *temp
	return newLabel, err
}

// AddLabel Add a label to a issue.
func AddLabel(repo, labelName string, issueNum int) ([]*github.Label, error) {
	s := strings.Split(repo, "/")
	label := []string{labelName}
	ctx := context.Background()
	labels, _, err := client.Issues.AddLabelsToIssue(ctx, s[0], s[1], issueNum, label)
	return labels, err
}

// EditLabel Change the name of a label.
func EditLabel(repo, labelName, newName string) (github.Label, error) {
	s := strings.Split(repo, "/")
	label := new(github.Label)
	label.Name = &newName
	ctx := context.Background()
	temp, _, err := client.Issues.EditLabel(ctx, s[0], s[1], labelName, label)
	editedLabel := *temp
	if err == nil {
		_, err = Issues(repo)
	}
	return editedLabel, err
}

// RemoveLabel Remove a label from an issue.
func RemoveLabel(repo, labelName string, issueNum int) error {
	s := strings.Split(repo, "/")
	ctx := context.Background()
	_, err := client.Issues.RemoveLabelForIssue(ctx, s[0], s[1], issueNum, labelName)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

// DeleteLabel Delete a label from a repo.
func DeleteLabel(repo, labelName string) error {
	s := strings.Split(repo, "/")
	ctx := context.Background()
	_, err := client.Issues.DeleteLabel(ctx, s[0], s[1], labelName)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}
