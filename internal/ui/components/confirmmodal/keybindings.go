package confirmmodal

import "github.com/awesome-gocui/gocui"

func (m *Modal) setFocus(g *gocui.Gui, focus field) error {
	m.focus = focus
	_, err := g.SetCurrentView(m.stops()[m.focus])
	return err
}

func (m *Modal) focusNext(g *gocui.Gui, v *gocui.View) error {
	return m.setFocus(g, (m.focus+1)%fieldCount)
}

func (m *Modal) focusPrev(g *gocui.Gui, v *gocui.View) error {
	return m.setFocus(g, (m.focus-1+fieldCount)%fieldCount)
}

// BindKeys registers Tab/Shift+Tab focus-cycling and Esc dismiss on both
// stops. Enter's meaning depends on which stop it's pressed on: onConfirm
// on the confirm button, onDismiss on the dismiss button.
func (m *Modal) BindKeys(g *gocui.Gui, onConfirm, onDismiss func(*gocui.Gui, *gocui.View) error) error {
	for _, name := range m.stops() {
		enter := onConfirm
		if name == m.dismiss.Name() {
			enter = onDismiss
		}

		bindings := []struct {
			key interface{}
			fn  func(*gocui.Gui, *gocui.View) error
		}{
			{gocui.KeyTab, m.focusNext},
			{gocui.KeyBacktab, m.focusPrev},
			{gocui.KeyEnter, enter},
			{gocui.KeyEsc, onDismiss},
		}

		for _, kb := range bindings {
			if err := g.SetKeybinding(name, kb.key, gocui.ModNone, kb.fn); err != nil {
				return err
			}
		}
	}
	return nil
}
