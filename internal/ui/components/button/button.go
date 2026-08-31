// Package button is a small, bordered, one-line gocui button with a fixed
// accent color and a filled, reverse-video look when selected, for reuse
// anywhere a form needs an action.
package button

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

// Height is the y1-y0 delta a caller must reserve for one Button.
const Height = 2

var roundedFrameRunes = []rune{'─', '│', '╭', '╮', '╰', '╯'}

// colorSGR maps gocui's 8 basic color constants to their ANSI foreground
// SGR code, since a button's caption is drawn as raw escape codes rather
// than through gocui's own Attribute-based coloring.
var colorSGR = map[gocui.Attribute]int{
	gocui.ColorBlack:   30,
	gocui.ColorRed:     31,
	gocui.ColorGreen:   32,
	gocui.ColorYellow:  33,
	gocui.ColorBlue:    34,
	gocui.ColorMagenta: 35,
	gocui.ColorCyan:    36,
	gocui.ColorWhite:   37,
}

type Config struct {
	BaseName string
	Caption  string
	// Color is the button's accent color, used for both its border and
	// its caption. Must be one of gocui's 8 basic color constants.
	Color gocui.Attribute
}

type Button struct {
	name       string
	borderName string
	caption    string
	color      gocui.Attribute
	sgr        int
}

func New(cfg Config) *Button {
	sgr, ok := colorSGR[cfg.Color]
	if !ok {
		sgr = 39 // terminal default foreground
	}
	return &Button{
		name:       cfg.BaseName,
		borderName: cfg.BaseName + "-border",
		caption:    cfg.Caption,
		color:      cfg.Color,
		sgr:        sgr,
	}
}

// Name is the gocui view name to focus or bind keys on. It never carries a
// border of its own - see Layout.
func (b *Button) Name() string { return b.name }

// Layout draws the button as two views at the same rectangle: a
// border-only view that is never made the current view, so its accent
// color is never overridden by gocui's global focus highlight (gocui
// always recolors whichever view is g.currentView's border to the app's
// shared selection color, with no per-view opt-out), and a borderless
// content view (Name) that receives actual focus/keybindings and holds
// the caption drawn by Render.
func (b *Button) Layout(g *gocui.Gui, x0, y0, x1 int) error {
	border, err := g.SetView(b.borderName, x0, y0, x1, y0+Height, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		border.Frame = true
		border.FrameRunes = roundedFrameRunes
		border.FrameColor = b.color
	}

	v, err := g.SetView(b.name, x0, y0, x1, y0+Height, 0)
	if err != nil && err != gocui.ErrUnknownView {
		return err
	}
	if err == gocui.ErrUnknownView {
		v.Frame = false
	}
	return nil
}

// Render redraws the caption, centered across the button's full width.
// When selected, the whole row fills as reverse video in the button's
// color, rather than leaving that job to gocui's own focus highlighting,
// which only recolors the border.
func (b *Button) Render(g *gocui.Gui, selected bool) {
	v, err := g.View(b.name)
	if err != nil {
		return
	}
	v.Clear()

	width, _ := v.Size()
	line := center(b.caption, width)

	style := fmt.Sprintf("\x1b[%d;1m", b.sgr)
	if selected {
		style = fmt.Sprintf("\x1b[%d;1;7m", b.sgr)
	}
	_, _ = fmt.Fprintf(v, "%s%s\x1b[0m", style, line)
}

// Delete removes the button's views. Layout recreates them the next time
// it's called.
func (b *Button) Delete(g *gocui.Gui) {
	_ = g.DeleteView(b.name)
	_ = g.DeleteView(b.borderName)
}

func center(s string, width int) string {
	pad := width - len(s)
	if pad <= 0 {
		return s
	}
	left := pad / 2
	right := pad - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
