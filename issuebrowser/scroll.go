package issuebrowser

import "github.com/jroimartin/gocui"

func scrollTop(g *gocui.Gui, v *gocui.View) error {
	if previousView.Name() == "browser" {
		err := scrollTopGet(g, v)
		if err != nil {
			return err
		}
	} else {
		if err := g.SetCurrentView(previousView.Name()); err != nil {
			return err
		}
		if previousView != nil {
			if err := previousView.SetOrigin(0, 0); err != nil {
				return err
			}
			if err := previousView.SetCursor(0, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

func scrollBottom(g *gocui.Gui, v *gocui.View) error {
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
		for l != "" || m != "" {
			if err := v.SetCursor(cx, cy+1); err != nil {
				ox, oy := v.Origin()
				if err := v.SetOrigin(ox, oy+1); err != nil {
					return err
				}
			}
			cy++
			if l, err = v.Line(cy + 1); err != nil {
				l = ""
			}
			if m, err = v.Line(cy + 2); err != nil {
				m = ""
			}
		}
	}
	return nil
}

func scrollTopGet(g *gocui.Gui, v *gocui.View) error {
	if err := g.SetCurrentView(previousView.Name()); err != nil {
		return err
	}
	if previousView != nil {
		if err := previousView.SetOrigin(0, 0); err != nil {
			return err
		}
		if err := previousView.SetCursor(0, 0); err != nil {
			return err
		}
	}
	if err := getLine(g, previousView); err != nil {
		return err
	}
	return nil
}

func scrollBottomGet(g *gocui.Gui, v *gocui.View) error {
	if err := scrollBottom(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func scrollup(g *gocui.Gui, v *gocui.View) error {
	_, maxY := v.Size()
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		for i := 0; i < maxY; i++ {
			if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
				if err := v.SetOrigin(ox, oy-1); err != nil {
					return err
				}
			}
			cy--
			oy--
		}
	}
	return nil
}

func scrolldown(g *gocui.Gui, v *gocui.View) error {
	_, maxY := v.Size()
	if v != nil {
		cx, cy := v.Cursor()
		var l string
		var m string
		var err error
		for i := 0; i < maxY; i++ {
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
			cy++
		}
	}
	return nil
}

func scrollEnd(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		var l string
		var m string
		var err error
		line, err := v.Line(cy)
		if err != nil {
			return err
		}
		for i := 0; i < len(line); i++ {
			if l, err = v.Word(cx, cy); err != nil {
				l = ""
			}
			if m, err = v.Word(cx+1, cy); err != nil {
				m = ""
			}
			if l != "" || m != "" {
				ox, oy := v.Origin()
				if err := v.SetCursor(cx+1, cy); err != nil {
					if err := v.SetOrigin(ox+1, oy); err != nil {
						return err
					}
				}
				cx++
			}
		}
	}
	return nil
}

func scrollHome(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		line, err := v.Line(cy)
		if err != nil {
			return err
		}
		for i := 0; i < len(line); i++ {
			if err := v.SetCursor(cx-1, cy); err != nil && ox > 0 {
				if err := v.SetOrigin(ox-1, oy); err != nil {
					return err
				}
			}
			cx--
			ox--
		}
	}
	return nil
}

func scrollRight(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		var l string
		var m string
		var err error
		if l, err = v.Word(cx, cy); err != nil {
			l = ""
		}
		if m, err = v.Word(cx+1, cy); err != nil {
			m = ""
		}
		if l != "" || m != "" {
			ox, oy := v.Origin()
			if err := v.SetCursor(cx+1, cy); err != nil {
				if err := v.SetOrigin(ox+1, oy); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func scrollLeft(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		ox, oy := v.Origin()
		if err := v.SetCursor(cx-1, cy); err != nil && ox > 0 {
			if err := v.SetOrigin(ox-1, oy); err != nil {
				return err
			}
		}
	}
	return nil
}

func scrollupget(g *gocui.Gui, v *gocui.View) error {
	if err := scrollup(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}

func scrolldownget(g *gocui.Gui, v *gocui.View) error {
	if err := scrolldown(g, v); err != nil {
		return err
	}
	if err := getLine(g, v); err != nil {
		return err
	}
	return nil
}
