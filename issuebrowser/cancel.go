package issuebrowser

import "github.com/jroimartin/gocui"

//cancel is used to close dialog boxes
func cancel(g *gocui.Gui, v *gocui.View) error {
	if (g.CurrentView()).Name() == "issueEd" {
		if err := g.DeleteView("issueEd"); err != nil {
			return err
		}
		if err := g.DeleteView("issueprompt"); err != nil {
			return err
		}
		entryCount = 0
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "sortChoice" {
		if err := g.DeleteView("sortChoice"); err != nil {
			return err
		}
		if err := g.DeleteView("sortPrompt"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "filterChoice" {
		if err := g.DeleteView("filterChoice"); err != nil {
			return err
		}
		if err := g.DeleteView("filterPrompt"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "commentBody" {
		if err := g.DeleteView("commentBody"); err != nil {
			return err
		}
		if err := g.DeleteView("commentPrompt"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "commentBrowser" || (g.CurrentView()).Name() == "commentViewer" {
		if err := g.DeleteView("commentEditPrompt"); err != nil {
			return err
		}
		if err := g.DeleteView("commentBrowser"); err != nil {
			return err
		}
		if err := g.DeleteView("commentViewer"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "commentDeleter" {
		if err := g.DeleteView("commentDeletePrompt"); err != nil {
			return err
		}
		if err := g.DeleteView("commentDeleter"); err != nil {
			return err
		}
		if err := g.DeleteView("commentViewer"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "labelBrowser" {
		if err := g.DeleteView("labelPrompt"); err != nil {
			return err
		}
		if err := g.DeleteView("labelBrowser"); err != nil {
			return err
		}
		if err := g.DeleteView("labelViewer"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
		if changed == true {
			changed = false
			if err := refresh(g, v); err != nil {
				return err
			}
		}
	} else if (g.CurrentView()).Name() == "labelRemover" {
		if err := g.DeleteView("labelRemover"); err != nil {
			return err
		}
		if err := g.DeleteView("labelPrompt"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	} else if (g.CurrentView()).Name() == "helpPane" {
		if err := g.DeleteView("helpPane"); err != nil {
			return err
		}
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
	}
	return nil
}
