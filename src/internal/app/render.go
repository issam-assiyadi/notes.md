package app

import (
	"fmt"

	"github.com/jroimartin/gocui"
	"portoflio.com/cli/internal/markdown"
)

func (a *App) Render(g *gocui.Gui) error {
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

	return a.Content.Render(g, func(v *gocui.View, contentWidth int) error {
		if len(a.Pages) == 0 {
			fmt.Fprintln(v, "No pages")
			return nil
		}

		if contentWidth <= 0 {
			contentWidth = 80
		}

		md := a.Pages[a.CurrentPage].Content
		// FIXME: I should investigate more in that, it seems like we do this in a wrong way.
		rendered, err := markdown.Render(md, contentWidth)
		if err != nil {
			fmt.Fprint(v, md)
			return nil
		}

		fmt.Fprint(v, rendered)
		return nil
	})
}
