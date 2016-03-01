package gitissue

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Username string
	Token    string
}

var (
	client          = github.NewClient(nil)
	conf            *Config
	folderpath, err = filepath.Abs(".")
	path            = folderpath + "/"
)

func write(toWrite []github.Issue, file string) error {
	file = path + ".issue/" + file + ".json"
	b, err := json.Marshal(toWrite)
	if err == nil {
		err = ioutil.WriteFile(file, b, 0644)
	}
	return err
}

func SetUp(user, oauth string) error {
	_, err := os.Stat(path + ".git")
	if os.IsNotExist(err) {
		return err
	}
	_, err = os.Stat(path + ".issue")
	if os.IsNotExist(err) {
		err := os.Mkdir(path+".issue", 0755)
		if err != nil {
			log.Println("make folder: ", err)
			return err
		}
	}
	_, err = ioutil.ReadFile(path + ".issue/config.json")
	if err != nil {
		temp := Config{user, oauth}
		b, err := json.Marshal(temp)
		if err == nil {
			err = ioutil.WriteFile(path+".issue/config.json", b, 0644)
		} else {
			return err
		}
	}
	return nil
}

func Login() error {
	file, err := ioutil.ReadFile(path + ".issue/config.json")
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

func IssuesFilter(repo, state, milestone, assignee, creator, sort, order string, labels []string) ([]github.Issue, error) {
	s := strings.Split(repo, "/")
	sorting := new(github.IssueListByRepoOptions)
	if len(labels) != 0 {
		sorting.Labels = labels
	}
	if state != "" {
		sorting.State = state
	}
	if milestone != "" {
		sorting.Milestone = milestone
	}
	if assignee != "" {
		sorting.Assignee = assignee
	}
	if creator != "" {
		sorting.Creator = creator
	}
	if sort != "" {
		sorting.Sort = sort
	}
	if order != "" {
		sorting.Direction = order
	}
	issues, _, err := client.Issues.ListByRepo(s[0], s[1], sorting)
	return issues, err
}

func Issues(repo string) ([]github.Issue, error) {
	issues, err := IssuesFilter(repo, "", "", "", "", "", "", nil)
	if err == nil {
		write(issues, "issues")
		return issues, err
	}
	file, err := ioutil.ReadFile(path + ".issue/issues.json")
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
