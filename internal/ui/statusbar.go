package ui

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

const statusBarViewName = "statusbar"

const statusBarStyle = "\x1b[38;5;244m"

const browseStatusBarText = "↑/↓ | k/j Move   Enter/o Open   r Rescan   <leader>e Categories   c Config   ? Help   q Quit"

// configStatusBarText covers only Tab/Shift+Tab: it's the one hint that's
// always accurate regardless of which field or button has focus. Save
// and Cancel carry their own <leader> hints directly on the buttons (see
// formpage.New) instead of living here, since '?' Help and a described
// <leader>q both stopped being reliably true across every focus state on
// this page.
const configStatusBarText = "Tab  Next field    Shift+Tab  Previous field"

// statusBarText picks the hint line for whichever view is currently
// showing, so it never advertises a key that doesn't apply there.
func (a *App) statusBarText() string {
	if a.ConfigForm.IsOpen() {
		return configStatusBarText
	}
	return browseStatusBarText
}

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
	_, err = fmt.Fprintf(v, "%s%s%s", statusBarStyle, a.statusBarText(), rowStyleReset)
	return err
}
