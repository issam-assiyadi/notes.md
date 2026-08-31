package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

const statusBarViewName = "statusbar"

const statusBarStyle = "\x1b[38;5;244m"

const statusBarText = "↑/↓ | k/j Move   Enter/Space Open   r Rescan   c Categories   ? Help   q Quit"

func (a *App) layoutStatusBar(g *gocui.Gui, maxX, maxY int) error {
	v, err := g.SetView(statusBarViewName, -1, maxY-2, maxX, maxY, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		v.Frame = false
	}
	return nil
}

func (a *App) renderStatusBar(g *gocui.Gui) error {
	v, err := g.View(statusBarViewName)
	if err != nil {
		return nil
	}
	v.Clear()
	_, err = fmt.Fprintf(v, "%s%s%s", statusBarStyle, statusBarText, rowStyleReset)
	return err
}
