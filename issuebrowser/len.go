package issuebrowser

//Len returns the length of the array to be sorted
func (iss byNumber) Len() int {
	return len(iss)
}
func (iss byTitle) Len() int {
	return len(iss)
}
func (iss byBody) Len() int {
	return len(iss)
}
func (iss byUser) Len() int {
	return len(iss)
}
func (iss byAssignee) Len() int {
	return len(iss)
}
func (iss byComments) Len() int {
	return len(iss)
}
func (iss byClosedAt) Len() int {
	return len(iss)
}
func (iss byCreatedAt) Len() int {
	return len(iss)
}
func (iss byUpdatedAt) Len() int {
	return len(iss)
}
func (iss byMilestone) Len() int {
	return len(iss)
}
