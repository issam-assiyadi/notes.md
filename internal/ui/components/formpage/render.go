package formpage

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

// Render redraws the Save/Cancel captions and the error line. It must never
// touch the Root or Ignore field views: their buffers hold live, in-progress
// user input, and Render runs every frame - clearing and rewriting them the
// way every other view in this codebase does would erase whatever the user
// is in the middle of typing.
func (p *Page) Render(g *gocui.Gui) error {
	if !p.open {
		return nil
	}

	p.save.Render(g, p.focus == fieldSave)
	p.cancel.Render(g, p.focus == fieldCancel)

	if v, err := g.View(p.errorLineName()); err == nil {
		v.Clear()
		if p.errorText != "" {
			_, _ = fmt.Fprintf(v, "\x1b[31m%s\x1b[0m", p.errorText)
		}
	}
	return nil
}
