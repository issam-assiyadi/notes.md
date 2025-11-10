package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// N.B: we will call this function in every re-size and when we call g.Update()

func (a *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// TODO: check if we have enough space to render our app
	// If we don't we will display a msg, (like: Not enough space to render the app)
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
	sv.Clear()

	// draw sidebar content 
	for i, p := range a.Pages {
		prefix := "  "
		if i == a.CurrentPage {
			prefix = "> "
		}

		fmt.Fprintf(sv, "%s%s\n", prefix, p.Title)
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
	cv.Clear()

	// draw the content
	if len(a.Pages) == 0 {
		fmt.Fprintln(cv, "No pages")
		return nil
	}

	// TODO: render the md content
	fmt.Fprintf(cv, "%s", a.Pages[a.CurrentPage].Content)

	return nil
}
