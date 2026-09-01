package formpage

import (
	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/internal/ui/components/button"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/textinput"
)

// lineHeight is the y1-y0 delta the error line (the one piece of this
// page's content not built from a reusable component) needs to show one
// row of text - see textinput.Height's doc comment for why this is 2, not
// 1.
const lineHeight = 2

const sectionGap = 1

// Layout draws the page across the given rectangle: actionview owns the
// border, title, and the Save/Cancel row pinned to its bottom, and the
// fields stack downward from the top of whatever content rectangle it
// hands back.
func (p *Page) Layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if !p.open {
		return nil
	}

	contentX0, contentY0, contentX1, contentY1, err := p.av.Layout(g, x0, y0, x1, y1, p.title, []*button.Button{p.save, p.cancel})
	if err != nil {
		return err
	}
	if contentX0 >= contentX1 || contentY0 >= contentY1 {
		return nil
	}

	y := contentY0
	if err := p.root.Layout(g, contentX0, y, contentX1, p.rootValue); err != nil {
		return err
	}
	y += textinput.Height + sectionGap

	if err := p.ignore.Layout(g, contentX0, y, contentX1, p.ignoreValue); err != nil {
		return err
	}
	y += textinput.Height + sectionGap

	return p.layoutErrorLine(g, contentX0, y, contentX1)
}

func (p *Page) layoutErrorLine(g *gocui.Gui, x0, y0, x1 int) error {
	v, err := g.SetView(p.errorLineName(), x0, y0, x1, y0+lineHeight, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		v.Frame = false
		v.Wrap = false
	}
	return nil
}
