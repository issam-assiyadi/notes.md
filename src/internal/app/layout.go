package app

import (
	"github.com/jroimartin/gocui"
)

func (a *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if maxX <= 0 || maxY <= 0 {
		return nil
	}

	sidebarWidth := max(maxX/4, 20)

	sv, err := g.SetView(a.SidebarView, 0, 0, sidebarWidth-1, maxY-1)

	// err == gocui.ErrUnknownView => view doesn't exist
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		sv.Title = "Pages"
		sv.Highlight = true
		sv.SelBgColor = gocui.ColorGreen
		sv.Wrap = true
	}

	cv, err := g.SetView(a.ContentView, sidebarWidth, 0, maxX-1, maxY-1)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		cv.Title = "Content"
		cv.Wrap = true
		cv.Autoscroll = false
	}

	return nil
}
