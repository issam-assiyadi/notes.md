package formmodal

import (
	"fmt"

	"github.com/awesome-gocui/gocui"
)

// Render redraws the Save/Cancel captions and the error line. It must never
// touch the Root or Ignore field views: their buffers hold live, in-progress
// user input, and Render runs every frame - clearing and rewriting them the
// way every other view in this codebase does would erase whatever the user
// is in the middle of typing.
func (m *Modal) Render(g *gocui.Gui) error {
	if !m.open {
		return nil
	}

	m.save.Render(g, m.focus == fieldSave)
	m.cancel.Render(g, m.focus == fieldCancel)

	if v, err := g.View(m.errorLineName()); err == nil {
		v.Clear()
		if m.errorText != "" {
			fmt.Fprintf(v, "\x1b[31m%s\x1b[0m", m.errorText)
		}
	}
	return nil
}
