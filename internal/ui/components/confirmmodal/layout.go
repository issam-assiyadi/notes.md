package confirmmodal

import (
	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/internal/ui/components/button"
)

var roundedFrameRunes = []rune{'─', '│', '╭', '╮', '╰', '╯'}

// messageHeight is the y1-y0 delta the message line needs to show one row
// of text - see textinput.Height's doc comment for why this is 2, not 1.
const messageHeight = 2

const (
	width      = 56
	minWidth   = 40
	buttonGap  = 2
	sectionGap = 1
)

// bounds returns a fixed-size box centered in the terminal, clamped to a
// usable minimum.
func bounds(maxX, maxY int) (x0, y0, x1, y1 int) {
	w := width
	if w > maxX-2 {
		w = maxX - 2
	}
	if w < minWidth {
		w = min(minWidth, maxX)
	}

	h := messageHeight + button.Height + sectionGap*2 + 2
	if h > maxY-2 {
		h = maxY - 2
	}

	x0 = (maxX - w) / 2
	y0 = (maxY - h) / 2
	x1 = x0 + w - 1
	y1 = y0 + h - 1
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
	if err := m.layoutMessage(g, innerLeft, y, innerRight); err != nil {
		return err
	}
	y += messageHeight + sectionGap

	return m.layoutButtons(g, innerLeft, y, innerRight)
}

func (m *Modal) layoutMessage(g *gocui.Gui, x0, y0, x1 int) error {
	v, err := g.SetView(m.messageViewName(), x0, y0, x1, y0+messageHeight, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		v.Frame = false
		v.Wrap = true
	}
	return nil
}

func (m *Modal) layoutButtons(g *gocui.Gui, x0, y0, x1 int) error {
	mid := x0 + (x1-x0)/2
	if err := m.confirm.Layout(g, x0, y0, mid-buttonGap/2); err != nil {
		return err
	}
	return m.dismiss.Layout(g, mid+buttonGap/2, y0, x1)
}
