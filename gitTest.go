package main

import (
	"./gitissue"
	"fmt"
)

func main() {
	//gitissue.SetUp("a username", "a token")
	gitissue.Login()
	issues, err := gitissue.Issues("butlerx/butlerbot")
	if err == nil {
		fmt.Println(issues[0])
	}
	//fmt.Println(gitissue.Repos())
	//fmt.Println(gitissue.OrgsList())
}
