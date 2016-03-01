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
	issues, err := gitissue.Issues("butlerx/butlerbot")
	if err == nil {
		fmt.Println(issues[0])
	}
	//fmt.Println(gitissue.Repos())
	//fmt.Println(gitissue.OrgsList())
}
