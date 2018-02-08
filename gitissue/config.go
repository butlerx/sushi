package gitissue

// Config User config file
// Username and oauth token stored
// check if user wants the token secured
type Config struct {
	Username string
	Token    string
	Secure   bool
}
