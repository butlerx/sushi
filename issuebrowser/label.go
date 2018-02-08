package issuebrowser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/butlerx/sushi/gitissue"
	"github.com/jroimartin/gocui"
)

func addLabel(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	maxX, maxY := g.Size()
	if labelPropmt, err := g.SetView("labelPrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelPropmt.Wrap = true
		fmt.Fprintln(labelPropmt, "Please choose from the list of labels below.\nPress enter to add a label to the issue.\n\nCtrl+C to exit")
	}
	if labelBrowser, err := g.SetView("labelBrowser", maxX/4, maxY/3, maxX/2, maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelBrowser.Highlight = true
		for i := 0; i < len(labelList); i++ {
			fmt.Fprintln(labelBrowser, *labelList[i].Name)
		}
		if err := g.SetCurrentView("labelBrowser"); err != nil {
			return err
		}
	}

	if labelViewer, err := g.SetView("labelViewer", maxX/2, maxY/3, maxX-(maxX/4), maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		_, cy := browser.Cursor()
		issueLine, err := browser.Line(cy)
		if err != nil {
			return err
		}
		issueNum := (strings.Split(issueLine, ":"))[0]
		for i := 0; i < len(issueList); i++ {
			if issueNum == strconv.Itoa(*issueList[i].Number) {
				for j := 0; j < len(issueList[i].Labels); j++ {
					fmt.Fprintln(labelViewer, *issueList[i].Labels[j].Name)
				}
				break
			}
		}
	}
	return nil
}

func openLabelRemover(g *gocui.Gui, v *gocui.View) error {
	previousView = v
	maxX, maxY := g.Size()
	if labelPropmt, err := g.SetView("labelPrompt", maxX/4, maxY/6, maxX-(maxX/4), maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelPropmt.Wrap = true
		fmt.Fprintln(labelPropmt, "Please choose from the list of labels below.\nPress enter to remove a label from an issue.\n\nCtrl+C to exit")
	}
	if labelRemover, err := g.SetView("labelRemover", maxX/4, maxY/3, maxX-(maxX/4), maxY-(maxY/6)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		labelRemover.Highlight = true
		browser, err := g.View("browser")
		if err != nil {
			return err
		}
		_, cy := browser.Cursor()
		issueLine, err := browser.Line(cy)
		if err != nil {
			return err
		}
		issueNum := (strings.Split(issueLine, ":")[0])
		for i := 0; i < len(issueList); i++ {
			if issueNum == strconv.Itoa(*issueList[i].Number) {
				for j := 0; j < len(issueList[i].Labels); j++ {
					fmt.Fprintln(labelRemover, *issueList[i].Labels[j].Name)
				}
				break
			}
		}
		if err := g.SetCurrentView("labelRemover"); err != nil {
			return err
		}
	}
	return nil
}

func removeLabel(g *gocui.Gui, v *gocui.View) error {
	labelRemover, err := g.View("labelRemover")
	if err != nil {
		return err
	}
	_, cy := labelRemover.Cursor()
	label, err := labelRemover.Line(cy)
	if err != nil {
		return err
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
	issueNum, err := strconv.Atoi((strings.Split(issueLine, ":"))[0])
	if err != nil {
		return err
	}
	if err := gitissue.RemoveLabel(getRepo(), label, issueNum); err != nil {
		return err
	}
	if err := cancel(g, v); err != nil {
		return err
	}
	if err := refresh(g, v); err != nil {
		return err
	}
	return nil
}

func writeLabel(g *gocui.Gui, v *gocui.View) error {
	changed = true
	labelBrowser, err := g.View("labelBrowser")
	if err != nil {
		return err
	}
	_, cy := labelBrowser.Cursor()
	label, err := labelBrowser.Line(cy)
	if err != nil {
		return err
	}
	labelViewer, err := g.View("labelViewer")
	if err != nil {
		return err
	}
	fmt.Fprintln(labelViewer, label)
	browser, err := g.View("browser")
	if err != nil {
		return err
	}
	_, cy = browser.Cursor()
	issueLine, err := browser.Line(cy)
	if err != nil {
		return err
	}
	issueNum, err := strconv.Atoi((strings.Split(issueLine, ":"))[0])
	if err != nil {
		return err
	}
	_, err = gitissue.AddLabel(getRepo(), label, issueNum)
	if err != nil {
		return err
	}
	return nil
}
