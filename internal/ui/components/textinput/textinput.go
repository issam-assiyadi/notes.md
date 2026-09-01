// Package textinput is a single-line, bordered, editable gocui text field
// with rounded corners, for reuse anywhere a form needs a free-text value.
package textinput

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

// Height is the y1-y0 delta a caller must reserve for one Input: gocui
// insets a view's content by one cell from x0/y0/x1/y1, so a single line
// of text needs two rows of span.
const Height = 2

var roundedFrameRunes = []rune{'─', '│', '╭', '╮', '╰', '╯'}

type Config struct {
	BaseName string
	Title    string
	// ConsumeChord, if set, is asked before a plain printable rune is
	// inserted normally. Returning true means the rune completed
	// something else (typically a <leader> chord) and must not be
	// written into the field; the field otherwise behaves exactly like
	// gocui's own DefaultEditor.
	ConsumeChord func(g *gocui.Gui, v *gocui.View, ch rune) bool
}

type Input struct {
	name         string
	title        string
	consumeChord func(g *gocui.Gui, v *gocui.View, ch rune) bool
}

func New(cfg Config) *Input {
	return &Input{name: cfg.BaseName, title: cfg.Title, consumeChord: cfg.ConsumeChord}
}

// Name is the gocui view name to focus or bind keys on.
func (i *Input) Name() string { return i.name }

// Layout creates or repositions the field at (x0,y0)-(x1,y0+Height).
// initial pre-fills the field's text the first time it's created; it has
// no effect on later calls once the view already exists, so a caller may
// pass it on every frame without it clobbering in-progress typing.
func (i *Input) Layout(g *gocui.Gui, x0, y0, x1 int, initial string) error {
	v, err := g.SetView(i.name, x0, y0, x1, y0+Height, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		v.Frame = true
		v.FrameRunes = roundedFrameRunes
		v.Title = i.title
		v.Editable = true
		v.Wrap = false
		_, _ = fmt.Fprint(v, initial)
		_ = v.SetCursor(len([]rune(initial)), 0)
		scrollCursorIntoView(v)
		if i.consumeChord != nil {
			consume := i.consumeChord
			v.Editor = gocui.EditorFunc(func(view *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
				if ch != 0 && mod == 0 && consume(g, view, ch) {
					return
				}
				gocui.DefaultEditor.Edit(view, key, ch, mod)
			})
		}
	}
	return nil
}

// Value reads the field's current text.
func (i *Input) Value(g *gocui.Gui) string {
	v, err := g.View(i.name)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(v.Buffer())
}

// Focus scrolls the field's horizontal origin so its cursor is visible.
// gocui hides the terminal cursor outright when it falls outside a view's
// origin/size window, and nothing recomputes that just because the view
// becomes current again - call this whenever this field becomes focused,
// or a cursor left past the visible width (e.g. at the end of a long
// pre-filled value) can stay invisible after focus returns to it.
func (i *Input) Focus(g *gocui.Gui) error {
	v, err := g.View(i.name)
	if err != nil {
		return err
	}
	scrollCursorIntoView(v)
	return nil
}

// Delete removes the field's view. Layout recreates it, with a fresh
// pre-fill, the next time it's called.
func (i *Input) Delete(g *gocui.Gui) {
	_ = g.DeleteView(i.name)
}

func scrollCursorIntoView(v *gocui.View) {
	cx, _ := v.Cursor()
	maxX, _ := v.Size()
	ox := cx - maxX + 1
	if ox < 0 {
		ox = 0
	}
	_ = v.SetOrigin(ox, 0)
}
