package issuebrowser

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

//help displays a list of keybinds
func help(g *gocui.Gui, v *gocui.View) error {
	previousView = g.CurrentView()
	maxX, maxY := g.Size()
	if helpPane, err := g.SetView("helpPane", maxX/4, maxY/6, maxX-(maxX/4), maxY-(maxY/3)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(helpPane, "F1, '?' = Help")
		fmt.Fprintln(helpPane, "F2, 'm' = Toggle open/closed")
		//fmt.Fprintln(helpPane, "F4 = Filter")
		fmt.Fprintln(helpPane, "F6, 's' = Sort by heading")
		fmt.Fprintln(helpPane, "")
		fmt.Fprintln(helpPane, "▼ , 'j' = Down")
		fmt.Fprintln(helpPane, "▲ , 'k' = Up")
		fmt.Fprintln(helpPane, "◀ , 'h' = Left")
		fmt.Fprintln(helpPane, "▶ , 'l'	= Right")
		fmt.Fprintln(helpPane, "")
		fmt.Fprintln(helpPane, "Home, '0'	= Home")
		fmt.Fprintln(helpPane, "End, '$' = End")
		fmt.Fprintln(helpPane, "PgUp = Page Up")
		fmt.Fprintln(helpPane, "PgDn = Page Down")
		fmt.Fprintln(helpPane, "'gg' = Window Top")
		fmt.Fprintln(helpPane, "'G' = Window Bottom")
		fmt.Fprintln(helpPane, "")
		fmt.Fprintln(helpPane, "Ctrl + N = New Issue/Comment/Label")
		fmt.Fprintln(helpPane, "Ctrl + C = Cancel out of dialog box")
		fmt.Fprintln(helpPane, "Ctrl + R = Refresh issue list from remote repo")
		fmt.Fprintln(helpPane, "Ctrl + E = Edit Comment/Label")
		fmt.Fprintln(helpPane, "Ctrl + D = Delete Comment/Label")
		fmt.Fprintln(helpPane, "")
		fmt.Fprintln(helpPane, "Ctrl + W = Enter window changing mode, nav keys to pick a window")
		fmt.Fprintln(helpPane, "'gt' = Next window")
		fmt.Fprintln(helpPane, "'gT' = Previous window")

		if err := g.SetCurrentView("helpPane"); err != nil {
			return err
		}
	}
	return nil
}
