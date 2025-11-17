package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
)


func (a *App) Render (g *gocui.Gui) error {

	// render the sidebar content
	sv, err := g.View(a.SidebarView)

	if err != nil {
		return err
	}
	sv.Clear()

	for i, p := range a.Pages {
		prefix := "  "
		if i == a.CurrentPage {
			prefix = "> "
		}

		fmt.Fprintf(sv, "%s%s\n", prefix, p.Title)
	}


	// render the content
	cv, err := g.View(a.ContentView)
	
	if err != nil {
		return err 
	}
	cv.Clear()

	if len(a.Pages) == 0 {
		fmt.Fprintln(cv, "No pages")
		return nil
	}

	// TODO: render the md content
	fmt.Fprintf(cv, "%s", a.Pages[a.CurrentPage].Content)

	return nil
}
