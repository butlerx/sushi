package gitissue

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/google/go-github/github"
)

// SetUp Set up .issue folder, issues, comments & config file and begin logfile.
// Appends to to gitignore to ignore the config file.
// Checks if the files exist and if they dont creates them.
func SetUp(user, oauth, key string) error {
	GitLog = logSetUp()
	err := makeFolder()
	if err != nil {
		return err
	}
	_, err = ioutil.ReadFile(*Path + ".gitissue/config.json")
	if err != nil {
		err = setUser(user, oauth, key)
		if err != nil {
			return err
		}
	}
	_, err = ioutil.ReadFile(*Path + ".gitissue/issues.json")
	if err != nil {
		temp := new([]github.Issue)
		b, err := json.Marshal(temp)
		if err == nil {
			err = ioutil.WriteFile(*Path+".gitissue/issues.json", b, 0644)
		} else {
			return err
		}
	}
	_, err = ioutil.ReadFile(*Path + ".gitissue/comments.json")
	if err != nil {
		temp := new([]Comments)
		b, err := json.Marshal(temp)
		if err == nil {
			err = ioutil.WriteFile(*Path+".gitissue/comments.json", b, 0644)
		} else {
			return err
		}
	}
	_, err = ioutil.ReadFile(*Path + ".gitignore")
	if err == nil {
		f, _ := os.OpenFile(".gitignore", os.O_APPEND, 0446)
		_, _ = f.WriteString(".gitissue/config.json")
		f.Close()
	} else {
		ignore := []byte(".gitissue/config.json")
		err = ioutil.WriteFile(".gitignore", ignore, 0644)
		return err
	}
	return nil
}

// IsSetUp Check if users config is set up and if being run in git repo.
func IsSetUp() (bool, error) {
	if checkgit() == false {
		err := errors.New("Not git repo")
		return false, err
	}
	_, err = os.Stat(*Path + ".gitissue")
	if os.IsNotExist(err) {
		return false, nil
	}
	file, err := ioutil.ReadFile(*Path + ".gitissue/config.json")
	if err != nil {
		return false, nil
	}
	temp := new(Config)
	if err = json.Unmarshal(file, temp); err != nil {
		return false, nil
	}
	if temp.Username == "" {
		return false, nil
	}
	return true, nil
}

// Check if being run in a git repo or a child in a git repo.
func checkgit() bool {
	cmd := exec.Command("git", "status")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

// logSetUp: Sets up Log file and creates logger object.
// Returns GitLog to be used to log erros.
func logSetUp() *log.Logger {
	err := makeFolder()
	if err != nil {
		return nil
	}
	_, err = ioutil.ReadFile(*Path + ".gitissue/sushi.log")
	logFile := new(os.File)
	if err != nil {
		logFile, err = os.Create(*Path + ".gitissue/sushi.log")
	} else {
		logFile, err = os.OpenFile(*Path+".gitissue/sushi.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}
	if err != nil {
		log.Fatalln("Failed to open logfile: ", err)
	}
	GitLog := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return GitLog
}

func makeFolder() error {
	_, err = os.Stat(*Path + ".gitissue")
	if os.IsNotExist(err) {
		err := os.Mkdir(*Path+".gitissue", 0755)
		if err != nil {
			log.Println("make folder: ", err)
			return err
		}
	}
	return nil
}
