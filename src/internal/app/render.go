package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"portoflio.com/cli/internal/markdown"
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

	contentWidth, _ := cv.Size()
	if contentWidth <= 0 {
		contentWidth = 80
	}

	md := a.Pages[a.CurrentPage].Content

	rendered, err := markdown.Render(md, contentWidth-1)
	if err != nil {
		fmt.Fprint(cv, md)
		return nil
	}

	fmt.Fprint(cv, rendered)
	return nil
}
