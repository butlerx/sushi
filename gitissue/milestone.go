package gitissue

import (
	"strings"

	"github.com/google/go-github/github"
)

// CreateMilestone Create a Milestone for a repo.
func CreateMilestone(repo, milestone string) (github.Milestone, error) {
	s := strings.Split(repo, "/")
	temp := new(github.Milestone)
	temp.Title = &milestone
	temp, _, err := client.Issues.CreateMilestone(s[0], s[1], temp)
	ms := *temp
	return ms, err
}

// AddMilestone Add Milestone to an issue.
// BUG(butlerx) Currently adding a milestone is not supported as milestones in the api are a mix of strings and ints.
// Bug is noted in library docs.
func AddMilestone() error {
	return error.new("Not yet implemented")
}

// ListMilestones List all Milestones in a repo.
func ListMilestones(repo string) ([]github.Milestone, error) {
	s := strings.Split(repo, "/")
	milestones, _, err := client.Issues.ListMilestones(s[0], s[1], nil)
	return milestones, err
}

// EditMilestone Change the title of a milestone in a repo.
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

// RemoveMilestone Remove Milestone to an issue.
// BUG(butlerx) currently Removing milestones is not supported as milestones in the api are a mix of strings and ints.
// Bug is noted in library docs.
func RemoveMilestone() error {
	return error.new("Not yet implemented")
}

// DeleteMilestone Delet a Milestone from a repo.
func DeleteMilestone(repo string, mileNum int) error {
	s := strings.Split(repo, "/")
	_, err := client.Issues.DeleteMilestone(s[0], s[1], mileNum)
	return err
}
