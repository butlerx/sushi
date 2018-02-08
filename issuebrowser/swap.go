package issuebrowser

//Swap indicates the method by which the sort package should swap elements
func (iss byNumber) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byTitle) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byBody) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byUser) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byAssignee) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byComments) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byClosedAt) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byCreatedAt) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byUpdatedAt) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
func (iss byMilestone) Swap(i, j int) {
	temp := iss[i]
	iss[i] = iss[j]
	iss[j] = temp
}
