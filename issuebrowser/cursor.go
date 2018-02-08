package issuebrowser

import "github.com/jroimartin/gocui"

func cursordown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		var l string
		var m string
		var err error
		if l, err = v.Line(cy + 1); err != nil {
			l = ""
		}
		if m, err = v.Line(cy + 2); err != nil {
			m = ""
		}
		if l != "" || m != "" {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func cursorup(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursordownGetIssues(g *gocui.Gui, v *gocui.View) error {
	if err := cursordown(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func cursorupGetIssues(g *gocui.Gui, v *gocui.View) error {
	if err := cursorup(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func cursorupGetComments(g *gocui.Gui, v *gocui.View) error {
	if err := cursorup(g, v); err != nil {
		return err
	}
	if err := getLineComment(g, v); err != nil {
		return err
	}
	return nil
}

func cursordownGetComments(g *gocui.Gui, v *gocui.View) error {
	if err := cursordown(g, v); err != nil {
		return err
	}
	if err := getLineComment(g, v); err != nil {
		return err
	}
	return nil
}
