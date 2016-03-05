package main

import (
	"fmt"
	"os"

	"github.com/butlerx/AgileGit/issuebrowser"
	"github.com/butlerx/AgileGit/startscreen"
)

func main() {
	var path string
	switch {
	case len(os.Args) == 1:
		startscreen.Show()

	case os.Args[1] == "list":
		if len(os.Args) != 2 { //check to see if you have cmd line args
			path = os.Args[2]
			issuebrowser.PassArgs(path)
		}
		issuebrowser.Show()

	case os.Args[1] == "init":
		fmt.Println("Should init a new .issue folder with a feedback message")

	case os.Args[1] == "preq":
		fmt.Println("Should display pull request manager")

	case os.Args[1] == "create":
		fmt.Println("create new issue dialog box")
	}
}
