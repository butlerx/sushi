package main

import (
	"./gitissue"
	"fmt"
)

func main() {
	gitissue.SetUp()
	gitissue.Login()
	fmt.Println(gitissue.Issues("butlerx/butlerbot"))
	//fmt.Println(gitissue.Repos())
	//fmt.Println(gitissue.OrgsList())
}
