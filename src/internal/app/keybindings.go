package app

import (
	"github.com/jroimartin/gocui"
	"portoflio.com/cli/internal/app/actions"
)

func (a *App) BindKeys (g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }); err != nil {
		return err
	}

	// FIXME: this is wrong we should link this to the content view, becauser we will have
	// conflic with moving pages.
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func (g *gocui.Gui, v *gocui.View) error {		
		// you can see this as a re-try to get the view if we get a null one (since we have a gui instance).
		if v == nil {
			var err error
			v, err = g.View(a.ContentView)
			if err != nil {
				return err
			}
		}		
		return actions.ScrollDown(v)
	}); err != nil {
		return err
	}

	// FIXME: same here
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func (g *gocui.Gui, v *gocui.View) error {		
		if v == nil {
			var err error
			v, err = g.View(a.ContentView)
			if err != nil {
				return err
			}
		}	
		return actions.ScrollUp(v)
	}); err != nil {
		return err
	}

	return nil
}