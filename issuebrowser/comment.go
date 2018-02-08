package issuebrowser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/butlerx/sushi/gitissue"
	"github.com/jroimartin/gocui"
)

//newComment opens the new comment dialog box
func newComment(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	maxX, maxY := g.Size()
	if commentPrompt, err := g.SetView("commentPrompt", maxX/4, maxY/3, maxX-(maxX/4), (maxY/3)+(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(commentPrompt, "Please enter you comment text.\nPress enter to write out.\n\nPress Ctrl+C to cancel")
	}
	if commentBody, err := g.SetView("commentBody", maxX/4, (maxY/3)+(maxY/6), maxX-(maxX/4), maxY-(maxY/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentBody.Editable = true
	}
	if err := g.SetCurrentView("commentBody"); err != nil {
		return err
	}
	return nil
}

//writeComment writes out the users comment entered into the new comment dialog box
func writeComment(g *gocui.Gui, v *gocui.View) error {
	comment := v.Buffer()
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	_, cy := browser.Cursor()
	line, err := browser.Line(cy)
	if err != nil {
		return err
	}
	issueNum := strings.Split(line, ":")
	issueIndex := 0
	for ; issueIndex < len(issueList); issueIndex++ {
		if issueNum[0] == strconv.Itoa(*issueList[issueIndex].Number) {
			break
		}
	}
	commentIndex := 0
	for ; commentIndex < len(comments); commentIndex++ {
		if len(comments[commentIndex]) > 0 {
			if *comments[commentIndex][0].IssueURL == *issueList[issueIndex].URL {
				break
			}
		}
	}
	num, err := strconv.Atoi(issueNum[0])
	if err != nil {
		return err
	}
	_, err = gitissue.Comment(getRepo(), comment, num)
	if err != nil {
		return err
	}
	if err := cancel(g, v); err != nil {
		return err
	}
	refresh(g, v)
	return nil
}

//openCommentEditor opens a list of comments for editing
func openCommentEditor(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	issueIndex := 0
	commentIndex := 0
	maxX, maxY := g.Size()
	if commentEditPrompt, err := g.SetView("commentEditPrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(commentEditPrompt, "Select the comment you wish to edit\n\nCtrl+C to cancel")
	}
	if commentBrowser, err := g.SetView("commentBrowser", maxX/4, maxY/3, maxX/2, maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentBrowser.Highlight = true
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		_, cy := browser.Cursor()
		issueLine, err := browser.Line(cy)
		if err != nil {
			return err
		}
		issueNum := strings.Split(issueLine, ":")
		var URL string
		for ; issueIndex < len(issueList); issueIndex++ {
			if issueNum[0] == strconv.Itoa(*issueList[issueIndex].Number) {
				URL = *issueList[issueIndex].URL
				break
			}
		}
		if *issueList[issueIndex].Comments > 0 {
			for ; commentIndex < len(comments); commentIndex++ {
				if len(comments[commentIndex]) > 0 {
					if URL == *comments[commentIndex][0].IssueURL {
						break
					}
				}
			}
			for i := 0; i < len(comments[commentIndex]); i++ {
				fmt.Fprintln(commentBrowser, strconv.Itoa(int(*comments[commentIndex][i].ID))+": "+*comments[commentIndex][i].User.Login+"@"+(*comments[commentIndex][i].CreatedAt).Format(time.UnixDate))
			}
		} else {
			fmt.Fprintln(commentBrowser, "This issue has no comments")
		}
		if err := g.SetCurrentView("commentBrowser"); err != nil {
			return err
		}
	}
	if commentViewer, err := g.SetView("commentViewer", maxX/2, maxY/3, maxX-(maxX/4), maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentViewer.Wrap = true
		if *issueList[issueIndex].Comments > 0 {
			fmt.Fprintln(commentViewer, *comments[commentIndex][0].Body)
		}
	}
	return nil
}

//editComment opens the individual comment chosed in openCommentEditor
func editComment(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	commentLine, err := v.Line(cy)
	if err != nil {
		return err
	}
	if commentLine == "This issue has no comments" {
		if err := cancel(g, v); err != nil {
			return err
		}
		return nil
	}
	commentEditPrompt, err := g.View("commentEditPrompt")
	if err != nil {
		return err
	}
	commentEditPrompt.Clear()
	fmt.Fprintln(commentEditPrompt, "Press enter to write out changes\n\nCtrl+C to cancel")
	if err := g.SetCurrentView("commentViewer"); err != nil {
		return err
	}
	commentViewer, err := g.View("commentViewer")
	if err != nil {
		return err
	}
	commentViewer.Editable = true
	return nil
}

//writeEditedComment writes out changes made to a comment
func writeEditedComment(g *gocui.Gui, v *gocui.View) error {
	commentBrowser, err := g.View("commentBrowser")
	if err != nil {
		return err
	}
	_, cy := commentBrowser.Cursor()
	commentLine, err := commentBrowser.Line(cy)
	if err != nil {
		return err
	}
	ID, err := strconv.Atoi((strings.Split(commentLine, ":"))[0])
	if err != nil {
		return err
	}
	if _, err = gitissue.EditComment(getRepo(), v.Buffer(), ID); err != nil {
		return err
	}
	v.Editable = false
	if err = cancel(g, v); err != nil {
		return err
	}
	if err = refresh(g, v); err != nil {
		return err
	}
	return nil
}

//openCommentDeleter opens a list of comments for deletion
func openCommentDeleter(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	issueIndex := 0
	commentIndex := 0
	maxX, maxY := g.Size()
	if commentDeletePrompt, err := g.SetView("commentDeletePrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(commentDeletePrompt, "Select the comment you wish to delete\n\nCtrl+C to cancel")
	}
	if commentDeleter, err := g.SetView("commentDeleter", maxX/4, maxY/3, maxX/2, maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentDeleter.Highlight = true
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		_, cy := browser.Cursor()
		issueLine, err := browser.Line(cy)
		if err != nil {
			return err
		}
		issueNum := strings.Split(issueLine, ":")
		var URL string
		for ; issueIndex < len(issueList); issueIndex++ {
			if issueNum[0] == strconv.Itoa(*issueList[issueIndex].Number) {
				URL = *issueList[issueIndex].URL
				break
			}
		}
		if *issueList[issueIndex].Comments > 0 {
			for ; commentIndex < len(comments); commentIndex++ {
				if len(comments[commentIndex]) > 0 {
					if URL == *comments[commentIndex][0].IssueURL {
						break
					}
				}
			}
			for i := 0; i < len(comments[commentIndex]); i++ {
				fmt.Fprintln(commentDeleter, strconv.Itoa(int(*comments[commentIndex][i].ID))+": "+*comments[commentIndex][i].User.Login+"@"+(*comments[commentIndex][i].CreatedAt).Format(time.UnixDate))
			}
		} else {
			fmt.Fprintln(commentDeleter, "This issue has no comments")
		}
		if err := g.SetCurrentView("commentDeleter"); err != nil {
			return err
		}
	}
	if commentViewer, err := g.SetView("commentViewer", maxX/2, maxY/3, maxX-(maxX/4), maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		commentViewer.Wrap = true
		if *issueList[issueIndex].Comments > 0 {
			fmt.Fprintln(commentViewer, *comments[commentIndex][0].Body)
		}
	}
	return nil
}

//deleteComment deletes a comment chosen in openCommentDeleter
func deleteComment(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	commentLine, err := v.Line(cy)
	if err != nil {
		return err
	}
	if commentLine == "This issue has no comments" {
		if err := cancel(g, v); err != nil {
			return err
		}
		return nil
	}
	ID, err := strconv.Atoi((strings.Split(commentLine, ":"))[0])
	if err != nil {
		return err
	}
	if err = gitissue.DeleteComment(getRepo(), ID); err != nil {
		return err
	}
	if err = cancel(g, v); err != nil {
		return err
	}
	if err = refresh(g, v); err != nil {
		return err
	}
	return nil
}

func getLineComment(g *gocui.Gui, v *gocui.View) error {
	issueIndex := 0
	commentIndex := 0
	_, cy := v.Cursor()
	commentLine, err := v.Line(cy)
	if err != nil {
		return err
	}
	if commentLine == "This issue has no comments" {
		return nil
	}
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	_, cy = browser.Cursor()
	issueLine, err := browser.Line(cy)
	if err != nil {
		return err
	}
	issueNum := strings.Split(issueLine, ":")
	var URL string
	for ; issueIndex < len(issueList); issueIndex++ {
		if issueNum[0] == strconv.Itoa(*issueList[issueIndex].Number) {
			URL = *issueList[issueIndex].URL
			break
		}
	}
	for ; commentIndex < len(comments); commentIndex++ {
		if len(comments[commentIndex]) > 0 {
			if URL == *comments[commentIndex][0].IssueURL {
				break
			}
		}
	}
	ID := strings.Split(commentLine, ":")
	IDnum, err := strconv.Atoi(ID[0])
	if err != nil {
		return err
	}
	commentViewer, err := g.View("commentViewer")
	if err != nil {
		return err
	}
	commentViewer.Clear()
	for i := 0; i < len(comments[commentIndex]); i++ {
		if IDnum == int(*comments[commentIndex][i].ID) {
			fmt.Fprintln(commentViewer, *comments[commentIndex][i].Body)
			break
		}
	}
	return nil
}
