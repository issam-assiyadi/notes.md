package formmodal

import (
	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/internal/ui/components/button"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/textinput"
)

var roundedFrameRunes = []rune{'─', '│', '╭', '╮', '╰', '╯'}

// lineHeight is the y1-y0 delta the error line (the one piece of this
// modal's content not built from a reusable component) needs to show one
// row of text - see textinput.Height's doc comment for why this is 2, not
// 1.
const lineHeight = 2

const (
	formWidth  = 60
	minWidth   = 40
	minHeight  = 14
	buttonGap  = 2
	sectionGap = 1
)

// bounds returns a fixed-size box centered in the terminal, clamped to a
// usable minimum.
func bounds(maxX, maxY int) (x0, y0, x1, y1 int) {
	width := formWidth
	if width > maxX-2 {
		width = maxX - 2
	}
	if width < minWidth {
		width = min(minWidth, maxX)
	}

	// Root field, Ignore field, an error line, and a button row, each
	// worth their own height, with a sectionGap row between them, plus
	// the wrapper's own top/bottom border.
	height := textinput.Height + textinput.Height + lineHeight + button.Height + sectionGap*3 + 2
	if height > maxY-2 {
		height = maxY - 2
	}
	if height < minHeight {
		height = min(minHeight, maxY)
	}

	x0 = (maxX - width) / 2
	y0 = (maxY - height) / 2
	x1 = x0 + width - 1
	y1 = y0 + height - 1
	if x1 >= maxX {
		x1 = maxX - 1
	}
	if y1 >= maxY {
		y1 = maxY - 1
	}
	return x0, y0, x1, y1
}

func (m *Modal) Layout(g *gocui.Gui, maxX, maxY int) error {
	if !m.open {
		return nil
	}
	if maxX <= 0 || maxY <= 0 {
		return nil
	}

	x0, y0, x1, y1 := bounds(maxX, maxY)
	if x0 >= x1 || y0 >= y1 {
		return nil
	}

	wrapper, err := g.SetView(m.wrapperName(), x0, y0, x1, y1, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		wrapper.Frame = true
		wrapper.FrameRunes = roundedFrameRunes
	}
	wrapper.Title = m.title

	innerLeft := x0 + 1
	innerRight := x1 - 1
	if innerLeft >= innerRight {
		return nil
	}

	y := y0 + 1
	if err := m.root.Layout(g, innerLeft, y, innerRight, m.rootValue); err != nil {
		return err
	}
	y += textinput.Height + sectionGap

	if err := m.ignore.Layout(g, innerLeft, y, innerRight, m.ignoreValue); err != nil {
		return err
	}
	y += textinput.Height + sectionGap

	if err := m.layoutErrorLine(g, innerLeft, y, innerRight); err != nil {
		return err
	}
	y += lineHeight + sectionGap

	return m.layoutButtons(g, innerLeft, y, innerRight)
}

func (m *Modal) layoutErrorLine(g *gocui.Gui, x0, y0, x1 int) error {
	v, err := g.SetView(m.errorLineName(), x0, y0, x1, y0+lineHeight, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		v.Frame = false
		v.Wrap = false
	}
	return nil
}

func (m *Modal) layoutButtons(g *gocui.Gui, x0, y0, x1 int) error {
	mid := x0 + (x1-x0)/2

	if err := m.save.Layout(g, x0, y0, mid-buttonGap/2); err != nil {
		return err
	}
	return m.cancel.Layout(g, mid+buttonGap/2, y0, x1)
}
