package gitissue

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Username string
	Token    string
}

var (
	client = github.NewClient(nil)
	conf   *Config
)

func Login() error {

	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Println("open config: ", err)
		os.Exit(1)
	}

	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		log.Println("parse config: ", err)
		os.Exit(1)
	}
	fmt.Println(temp)
	conf = temp
	fmt.Println(conf)

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

func Issues() ([]github.Issue, error) {
	issues, _, err := client.Issues.List(false, nil)
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
