package gitissue

import (
	"strings"

	"github.com/google/go-github/github"
)

// CreateLabel Create a label for a repo.
func CreateLabel(repo, labelName string) (github.Label, error) {
	s := strings.Split(repo, "/")
	label := new(github.Label)
	label.Name = &labelName
	temp, _, err := client.Issues.CreateLabel(s[0], s[1], label)
	newLabel := *temp
	return newLabel, err
}

// AddLabel Add a label to a issue.
func AddLabel(repo, labelName string, issueNum int) ([]github.Label, error) {
	s := strings.Split(repo, "/")
	label := []string{labelName}
	labels, _, err := client.Issues.AddLabelsToIssue(s[0], s[1], issueNum, label)
	return labels, err
}

// EditLabel Change the name of a label.
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

// RemoveLabel Remove a label from an issue.
func RemoveLabel(repo, labelName string, issueNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.RemoveLabelForIssue(s[0], s[1], issueNum, labelName)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}

// DeleteLabel Delete a label from a repo.
func DeleteLabel(repo, labelName string) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.DeleteLabel(s[0], s[1], labelName)
	if err == nil {
		_, err = Issues(repo)
	}
	return err
}
