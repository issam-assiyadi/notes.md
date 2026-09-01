// Package actionview lays out a bordered page as a content region on top
// and a row of actions pinned to the bottom, so the actions' screen
// position never depends on how much content sits above them - adding a
// field above can't push a Save button further down the page. It's
// geometry-only: it never touches a gocui buffer or a caller's own
// widgets' content, so it can't collide with a caller's own per-frame
// Render.
package actionview

import "github.com/awesome-gocui/gocui"

type Config struct {
	BaseName string
}

type View struct {
	wrapperName string
}

func New(cfg Config) *View {
	return &View{wrapperName: cfg.BaseName + "-wrapper"}
}

func (v *View) WrapperName() string { return v.wrapperName }

// Delete removes the wrapper view. Layout recreates it the next time
// it's called. The caller remains responsible for deleting its own
// content widgets and actions.
func (v *View) Delete(g *gocui.Gui) {
	if v == nil || g == nil {
		return
	}
	_ = g.DeleteView(v.wrapperName)
}
