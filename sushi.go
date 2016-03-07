package main

import (
	"fmt"
	"os"

	"github.com/butlerx/AgileGit/issuebrowser"
)

func main() {
	var path string
	switch {
	case len(os.Args) == 1:
		fmt.Println("Please add an arguement")

	case os.Args[1] == "ilist":
		if len(os.Args) != 2 { //check to see if you have cmd line args
			path = os.Args[2]
			issuebrowser.PassArgs(path)
		}
		issuebrowser.Show()

	case os.Args[1] == "init":
		fmt.Println("Should init a new .issue folder with a feedback message")

	case os.Args[1] == "pullreq":
		fmt.Println("Should display pull request manager")

	case os.Args[1] == "cissue":
		fmt.Println("Create new issue dialog box")

	case os.Args[1] == "config":
		fmt.Println("Write username and password to config file")
	}
}
