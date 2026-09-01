package actionview

import (
	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/internal/ui/components/button"
)

var roundedFrameRunes = []rune{'─', '│', '╭', '╮', '╰', '╯'}

// actionGap is the blank row left between the content region and the
// action row above it.
const actionGap = 1

// buttonGap is the blank columns left between adjacent actions.
const buttonGap = 2

// Layout draws the page's border and title across the given rectangle,
// lays out actions anchored to y1, and returns the rectangle
// (contentX0, contentY0, contentX1, contentY1) the caller should lay its
// own content into above them. The returned rectangle is zeroed (with a
// nil error) if the given one is too small to hold both regions.
func (v *View) Layout(g *gocui.Gui, x0, y0, x1, y1 int, title string, actions []*button.Button) (contentX0, contentY0, contentX1, contentY1 int, err error) {
	if x0 >= x1 || y0 >= y1 {
		return 0, 0, 0, 0, nil
	}

	wrapper, err := g.SetView(v.wrapperName, x0, y0, x1, y1, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return 0, 0, 0, 0, err
	}
	if err == gocui.ErrUnknownView {
		wrapper.Frame = true
		wrapper.FrameRunes = roundedFrameRunes
	}
	wrapper.Title = title

	innerLeft := x0 + 1
	innerRight := x1 - 1
	innerTop := y0 + 1
	innerBottom := y1 - 1
	if innerLeft >= innerRight || innerTop >= innerBottom {
		return 0, 0, 0, 0, nil
	}

	// button.Layout(g, x0, y0, x1) sets the button's bottom edge to
	// y0+Height, so anchoring that edge to innerBottom - one row above
	// the wrapper's own bottom border, matching how every other inset
	// widget in this codebase stops short of the frame - takes
	// innerBottom-Height, not innerBottom-Height+1.
	actionTop := innerBottom - button.Height
	if actionTop <= innerTop {
		return 0, 0, 0, 0, nil
	}

	if err := layoutActions(g, innerLeft, actionTop, innerRight, actions); err != nil {
		return 0, 0, 0, 0, err
	}

	return innerLeft, innerTop, innerRight, actionTop - actionGap, nil
}

// layoutActions divides the row evenly across all actions, with buttonGap
// blank columns between adjacent ones; the last action absorbs whatever
// remainder integer division leaves.
func layoutActions(g *gocui.Gui, x0, y0, x1 int, actions []*button.Button) error {
	n := len(actions)
	if n == 0 {
		return nil
	}

	width := x1 - x0 + 1
	slot := (width - buttonGap*(n-1)) / n

	cur := x0
	for i, b := range actions {
		end := cur + slot - 1
		if i == n-1 {
			end = x1
		}
		if err := b.Layout(g, cur, y0, end); err != nil {
			return err
		}
		cur = end + 1 + buttonGap
	}
	return nil
}
