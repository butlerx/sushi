package issuebrowser

import "github.com/jroimartin/gocui"

//keybindings assigns actions to specific keys depending on what window they are pressed in
func keybindings(g *gocui.Gui) error {
	mainWindows := []string{"browser", "issuepane", "commentpane", "labelpane", "milestonepane", "assigneepane"}
	displayWindows := []string{"issuepane", "commentpane", "labelpane", "milestonepane", "assigneepane", "helpPane", "labelBrowser", "labelRemover", "sortChoice", "selectionPane", "filterChoice"}
	controlWindows := []string{"browser", "issuepane", "commentpane", "labelpane", "milestonepane", "assigneepane", "helpPane", "labelBrowser", "labelRemover", "commentBrowser", "sortChoice", "selectionPane", "filterChoice"}
	dialogBoxes := []string{"helpPane", "labelBrowser", "labelRemover", "sortChoice", "filterChoice", "issueEd", "selectionPane", "commentBody", "commentDeleter", "commentBrowser", "commentViewer"}

	for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], gocui.KeyF1, gocui.ModNone, help); err != nil {
			return err
		}
		if err := g.SetKeybinding(mainWindows[i], '?', gocui.ModNone, help); err != nil {
			return err
		}
		if err := g.SetKeybinding(mainWindows[i], gocui.KeyF2, gocui.ModNone, toggleState); err != nil {
			return err
		}
		if err := g.SetKeybinding(mainWindows[i], 'm', gocui.ModNone, toggleState); err != nil {
			return err
		}
		if err := g.SetKeybinding(mainWindows[i], gocui.KeyF6, gocui.ModNone, showSortOrders); err != nil {
			return err
		}
		if err := g.SetKeybinding(mainWindows[i], 's', gocui.ModNone, showSortOrders); err != nil {
			return err
		}
		if err := g.SetKeybinding(mainWindows[i], gocui.KeyCtrlW, gocui.ModNone, changeWindow); err != nil {
			return err
		}
	}

	/*for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], gocui.KeyF4, gocui.ModNone, showFilterHeadings); err != nil {
			return err
		}
	}*/

	for i := 0; i < len(dialogBoxes); i++ {
		if err := g.SetKeybinding(dialogBoxes[i], gocui.KeyCtrlC, gocui.ModNone, cancel); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("selectionPane", gocui.KeyEnter, gocui.ModNone, nextEntry); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentBrowser", gocui.KeyArrowUp, gocui.ModNone, cursorupGetComments); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentBrowser", 'k', gocui.ModNone, cursorupGetComments); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentBrowser", gocui.KeyArrowDown, gocui.ModNone, cursordownGetComments); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentBrowser", 'j', gocui.ModNone, cursordownGetComments); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentBrowser", gocui.KeyEnter, gocui.ModNone, editComment); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentDeleter", gocui.KeyArrowUp, gocui.ModNone, cursorupGetComments); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentDeleter", 'k', gocui.ModNone, cursorupGetComments); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentDeleter", gocui.KeyArrowDown, gocui.ModNone, cursordownGetComments); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentDeleter", 'j', gocui.ModNone, cursordownGetComments); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentDeleter", gocui.KeyEnter, gocui.ModNone, deleteComment); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelBrowser", gocui.KeyEnter, gocui.ModNone, writeLabel); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelRemover", gocui.KeyEnter, gocui.ModNone, removeLabel); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", gocui.KeyCtrlD, gocui.ModNone, openCommentDeleter); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentViewer", gocui.KeyEnter, gocui.ModNone, writeEditedComment); err != nil {
		return err
	}

	if err := g.SetKeybinding("helpPane", gocui.KeyF1, gocui.ModNone, cancel); err != nil {
		return err
	}
	if err := g.SetKeybinding("helpPane", '?', gocui.ModNone, cancel); err != nil {
		return err
	}
	if err := g.SetKeybinding("issueEd", gocui.KeyEnter, gocui.ModNone, nextEntry); err != nil {
		return err
	}

	if err := g.SetKeybinding("sortChoice", gocui.KeyEnter, gocui.ModNone, getSortOrder); err != nil {
		return err
	}

	if err := g.SetKeybinding("filterChoice", gocui.KeyEnter, gocui.ModNone, getFilterHeading); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyCtrlN, gocui.ModNone, newIssue); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", gocui.KeyCtrlN, gocui.ModNone, newComment); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", gocui.KeyCtrlN, gocui.ModNone, addLabel); err != nil {
		return err
	}

	if err := g.SetKeybinding("labelpane", gocui.KeyCtrlD, gocui.ModNone, openLabelRemover); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentBody", gocui.KeyEnter, gocui.ModNone, writeComment); err != nil {
		return err
	}

	if err := g.SetKeybinding("commentpane", gocui.KeyCtrlE, gocui.ModNone, openCommentEditor); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyPgup, gocui.ModNone, scrollupget); err != nil {
		return err
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], gocui.KeyPgup, gocui.ModNone, scrollup); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("browser", gocui.KeyPgdn, gocui.ModNone, scrolldownget); err != nil {
		return err
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], gocui.KeyPgdn, gocui.ModNone, scrolldown); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("", gocui.KeyEnd, gocui.ModNone, scrollEnd); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyHome, gocui.ModNone, scrollHome); err != nil {
		return err
	}

	for i := 0; i < len(controlWindows); i++ {
		if err := g.SetKeybinding(controlWindows[i], '0', gocui.ModNone, scrollHome); err != nil {
			return err
		}
		if err := g.SetKeybinding(controlWindows[i], '$', gocui.ModNone, scrollEnd); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, scrollRight); err != nil {
		return err
	}

	for i := 0; i < len(controlWindows)-1; i++ {
		if err := g.SetKeybinding(controlWindows[i], 'l', gocui.ModNone, scrollRight); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, scrollLeft); err != nil {
		return err
	}

	for i := 0; i < len(controlWindows)-1; i++ {
		if err := g.SetKeybinding(controlWindows[i], 'h', gocui.ModNone, scrollLeft); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("browser", gocui.KeyArrowDown, gocui.ModNone, cursordownGetIssues); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", 'j', gocui.ModNone, cursordownGetIssues); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", gocui.KeyArrowUp, gocui.ModNone, cursorupGetIssues); err != nil {
		return err
	}

	if err := g.SetKeybinding("browser", 'k', gocui.ModNone, cursorupGetIssues); err != nil {
		return err
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], gocui.KeyArrowDown, gocui.ModNone, cursordown); err != nil {
			return err
		}
	}

	for i := 0; i < (len(displayWindows) - 1); i++ {
		if err := g.SetKeybinding(displayWindows[i], 'j', gocui.ModNone, cursordown); err != nil {
			return err
		}
	}

	for i := 0; i < len(displayWindows); i++ {
		if err := g.SetKeybinding(displayWindows[i], gocui.KeyArrowUp, gocui.ModNone, cursorup); err != nil {
			return err
		}
	}

	for i := 0; i < (len(displayWindows) - 1); i++ {
		if err := g.SetKeybinding(displayWindows[i], 'k', gocui.ModNone, cursorup); err != nil {
			return err
		}
	}

	for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], 'q', gocui.ModNone, quit); err != nil {
			return err
		}
	}

	for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], gocui.KeyTab, gocui.ModNone, toggleIssues); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, refresh); err != nil {
		return err
	}

	if err := g.SetKeybinding("windowChanger", gocui.KeyArrowUp, gocui.ModNone, windowUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", 'k', gocui.ModNone, windowUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", gocui.KeyArrowDown, gocui.ModNone, windowDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", 'j', gocui.ModNone, windowDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", gocui.KeyArrowRight, gocui.ModNone, windowRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", 'l', gocui.ModNone, windowRight); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", gocui.KeyArrowLeft, gocui.ModNone, windowLeft); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowChanger", 'h', gocui.ModNone, windowLeft); err != nil {
		return err
	}

	for i := 0; i < len(mainWindows); i++ {
		if err := g.SetKeybinding(mainWindows[i], 'g', gocui.ModNone, tabWindow); err != nil {
			return err
		}
	}

	if err := g.SetKeybinding("windowTabber", 't', gocui.ModNone, nextWindow); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowTabber", 'T', gocui.ModNone, previousWindow); err != nil {
		return err
	}
	if err := g.SetKeybinding("windowTabber", 'g', gocui.ModNone, scrollTop); err != nil {
		return err
	}
	for i := 0; i < (len(displayWindows) - 1); i++ {
		if err := g.SetKeybinding(displayWindows[i], 'G', gocui.ModNone, scrollBottom); err != nil {
			return err
		}
	}
	if err := g.SetKeybinding("browser", 'G', gocui.ModNone, scrollBottomGet); err != nil {
		return err
	}

	return nil
}
