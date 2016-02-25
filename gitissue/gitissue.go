package gitissue

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Config struct {
	Username string
	Token    string
}

var (
	client = github.NewClient(nil)
	conf   *Config
)

func write(toWrite []github.Issue) error {
	b, err := json.Marshal(toWrite)
	if err == nil {
		err = ioutil.WriteFile(".issue/issues.json", b, 0644)
	}
	return err
}

func SetUp() {
	err := os.Mkdir(".issue", 0755)
	if err != nil {
		log.Println("make folder: ", err)
	}
}

func Login() error {
	file, err := ioutil.ReadFile(".issue/config.json")
	if err != nil {
		log.Println("open config: ", err)
		os.Exit(1)
	}
	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
		os.Exit(1)
	}
	conf = temp
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: conf.Token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client = github.NewClient(tc)
	user, _, err := client.Users.Get("")

	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return err
	}
	fmt.Printf("\nLogged into: %v\n", github.Stringify(user.Name))
	return nil
}

func Issues(repo string) ([]github.Issue, error) {
	var issues []github.Issue
	var err error
	s := strings.Split(repo, "/")
	if repo == "" {
		issues, _, err = client.Issues.List(false, nil)
	} else {
		issues, _, err = client.Issues.ListByRepo(s[0], s[1], nil)
	}
	if err == nil {
		write(issues)
		return issues, err
	}
	file, err := ioutil.ReadFile(".issue/issues.json")
	if err != nil {
		log.Println("open issues: ", err)
		os.Exit(1)
	}
	temp := new([]github.Issue)
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse issues: ", err)
		os.Exit(1)
	}
	issues = *temp
	return issues, err
}

func Repos() ([]github.Repository, error) {
	repos, _, err := client.Repositories.List("", nil)
	return repos, err
}

func OrgsList() ([]github.Organization, error) {
	fmt.Println(conf.Username)
	orgs, _, err := client.Organizations.List(conf.Username, nil)
	return orgs, err
}
