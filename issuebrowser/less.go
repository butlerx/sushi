package issuebrowser

import "strings"

//Less returns true if element i is less than element j
//It is implemented by the sort package
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
