package main

import (
	"./gitissue"
	"fmt"
	"log"
)

func main() {
	err := gitissue.SetUp("a username", "a token")
	if err != nil {
		log.Println(err)
		return
	}
	gitissue.Login()
	/*issues, err := gitissue.Issues("butlerx/butlerbot")
	if err == nil {
		fmt.Println(issues)
	}
	//fmt.Println(gitissue.Repos())
	//fmt.Println(gitissue.OrgsList())*/
	labels := []string{}
	issue, err := gitissue.MakeIssue("butlerx/AgileGit", "test", "What the title says really", "butlerx", 0, labels)
	if err == nil {
		fmt.Println(issue)
	} else {
		log.Println(err)
	}
}
