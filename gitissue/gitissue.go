package gitissue

import (
	"github.com/butlerx/sushi/encrypt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	client          = github.NewClient(nil)
	Conf            *Config
	folderpath, err = filepath.Abs(".")
	folder          = (folderpath + "/")
	Path            = &folder
	GitLog          *log.Logger
)

// Set Auth token and password in config.
// Pass an empty string if the user doesnt want to secure there oauth
func setUser(user, oauth, userkey string) error {
	temp := new(Config)
	if userkey == "" {
		temp = &Config{user, oauth, false}
	} else {
		key := []byte(userkey)
		oauth := encrypt.Encrypt(key, oauth)
		temp = &Config{user, oauth, true}
	}
	b, err := json.Marshal(temp)
	if err == nil {
		err = ioutil.WriteFile(*Path+".issue/config.json", b, 0644)
		return err
	}
	return nil
}

// ChangeKey Change user encyption key
func ChangeKey(oldKey, newKey string) error {
	file, err := ioutil.ReadFile(*Path + ".issue/config.json")
	if err != nil {
		GitLog.Println("open config: ", err)
		os.Exit(1)
	}
	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		GitLog.Println("parse config: ", err)
		os.Exit(1)
	}
	Conf = temp
	token := Conf.Token
	if Conf.Secure {
		key := []byte(oldKey)
		token = encrypt.Decrypt(key, token)
		if err != nil {
			return err
		}
	}
	err = ChangeLogin(Conf.Username, token, newKey)
	return err
}

// ChangeLogin Change user name and Auth token and relogin.
func ChangeLogin(user, oauth, key string) error {
	err := setUser(user, oauth, key)
	if err != nil {
		return err
	}
	err = Login(key)
	return err
}

// Login Logs in to github using oauth.
// Returns error if login fails.
func Login(userkey string) error {
	file, err := ioutil.ReadFile(*Path + ".issue/config.json")
	if err != nil {
		GitLog.Println("open config: ", err)
		os.Exit(1)
	}
	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		GitLog.Println("parse config: ", err)
		os.Exit(1)
	}
	Conf = temp
	token := Conf.Token
	if Conf.Secure {
		key := []byte(userkey)
		token = encrypt.Decrypt(key, token)
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = github.NewClient(tc)
	user, _, err := client.Users.Get("")

	if err != nil {
		GitLog.Printf("\nerror: %v\n", err)
		return err
	}
	log.Printf("\nLogged into: %v\n", github.Stringify(user.Login))
	GitLog.Print("Logged into: ", github.Stringify(user.Login))
	return nil
}

// PossibleAssignees Get a list of all possible assingees
func PossibleAssignees(repo string) ([]github.User, error) {
	s := strings.Split(repo, "/")
	assignees, _, err := client.Issues.ListAssignees(s[0], s[1], nil)
	return assignees, err
}

// WatchRepo Monitor for change in repo.
// Rings terminal bell and
// returns true and the reason if something happened in repo.
// returns false if nothing changed
func WatchRepo(repo string) (string, string, bool) {
	s := strings.Split(repo, "/")
	subscription, _, err := client.Activity.GetRepositorySubscription(s[0], s[1])
	if err == nil && subscription != nil {
		notification, _, err := client.Activity.GetThread(*subscription.ThreadURL)
		if err == nil && *notification.Unread == true {
			_, err = client.Activity.MarkThreadRead(*subscription.ThreadURL)
			return *notification.Reason, *notification.Subject.Title, true, err
		}
	}
	return "", "", false, nil
}

// list all all of a users repos.
// currently unused.
func repos() ([]github.Repository, error) {
	repos, _, err := client.Repositories.List("", nil)
	return repos, err
}

// List all orgs a users a part of.
// currently unused.
func orgsList() ([]github.Organization, error) {
	GitLog.Println(Conf.Username)
	orgs, _, err := client.Organizations.List(Conf.Username, nil)
	return orgs, err
}
