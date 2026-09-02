package formpage

import "github.com/awesome-gocui/gocui"

func (p *Page) setFocus(g *gocui.Gui, focus field) error {
	p.focus = focus
	// Set together with the view change, in the same handler call, so a
	// redraw can never land between "current view changed" and "cursor
	// visibility updated" and show one without the other for a frame.
	g.Cursor = focus == fieldRoot || focus == fieldIgnore
	if _, err := g.SetCurrentView(p.stops()[p.focus]); err != nil {
		return err
	}
	switch focus {
	case fieldRoot:
		return p.root.Focus(g)
	case fieldIgnore:
		return p.ignore.Focus(g)
	}
	return nil
}

func (p *Page) focusNext(g *gocui.Gui, v *gocui.View) error {
	return p.setFocus(g, (p.focus+1)%fieldCount)
}

func (p *Page) focusPrev(g *gocui.Gui, v *gocui.View) error {
	return p.setFocus(g, (p.focus-1+fieldCount)%fieldCount)
}

// BindKeys registers Tab/Shift+Tab focus-cycling on all four stops, Esc
// cancel on all four stops, and Ctrl+S submit on all four stops (a
// save-from-anywhere shortcut). Enter's meaning depends on which stop it's
// pressed on: submit on the Root/Ignore fields and the Save button, cancel
// on the Cancel button - matching what each stop visually represents.
// These are all Key-type bindings rather than rune bindings, so gocui
// dispatches them before ever reaching the focused field's Editor - see
// matchView in gocui's keybinding.go: it only blocks rune keybindings on an
// Editable view.
func (p *Page) BindKeys(g *gocui.Gui, onSubmit, onCancel func(*gocui.Gui, *gocui.View) error) error {
	for _, name := range p.stops() {
		enter := onSubmit
		if name == p.cancel.Name() {
			enter = onCancel
		}

		bindings := []struct {
			key interface{}
			fn  func(*gocui.Gui, *gocui.View) error
		}{
			{gocui.KeyTab, p.focusNext},
			{gocui.KeyBacktab, p.focusPrev},
			{gocui.KeyEnter, enter},
			{gocui.KeyCtrlS, onSubmit},
			{gocui.KeyEsc, onCancel},
		}

		for _, kb := range bindings {
			if err := g.SetKeybinding(name, kb.key, gocui.ModNone, kb.fn); err != nil {
				return err
			}
		}
	}
	return nil
}
