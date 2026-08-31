package ui

import (
	"log"
	"time"

	"github.com/awesome-gocui/gocui"
)

type appKeyBinding struct {
	key  interface{}
	fn   func(*gocui.Gui, *gocui.View) error
	desc string
}

// leaderChord marks a global keybinding's second key in a <leader>-prefixed
// chord (Space, matching Vim's default leader, then the rune below). It
// exists only for display: formatKey renders it as "<leader>x", while
// BindKeys registers the underlying rune with gocui like any other key.
type leaderChord rune

const leaderTimeout = time.Second

func (a *App) globalKeyBindings() []appKeyBinding {
	return []appKeyBinding{
		{gocui.KeyCtrlC, a.quit, "Quit"},
		{'q', a.quit, "Quit"},
		{'r', a.rescan, "Rescan"},
		{'1', a.focusCategories, "Focus categories pane"},
		{'2', a.focusItems, "Focus items pane"},
		{leaderChord('e'), a.toggleCategories, "Toggle categories pane"},
		{'c', a.openConfigForm, "Edit project config"},
		{'?', a.openHelp, "Show keybinding help"},
	}
}

func (a *App) categoryKeyBindings() []appKeyBinding {
	return []appKeyBinding{
		{gocui.KeyArrowDown, a.categoryMoveDown, "Move down"},
		{'j', a.categoryMoveDown, "Move down"},
		{gocui.KeyArrowUp, a.categoryMoveUp, "Move up"},
		{'k', a.categoryMoveUp, "Move up"},
	}
}

func (a *App) itemKeyBindings() []appKeyBinding {
	return []appKeyBinding{
		{gocui.KeyArrowDown, a.rowMoveDown, "Move down"},
		{'j', a.rowMoveDown, "Move down"},
		{gocui.KeyArrowUp, a.rowMoveUp, "Move up"},
		{'k', a.rowMoveUp, "Move up"},
		{gocui.KeyEnter, a.activateSelectedRow, "Open item / toggle folder"},
		{'o', a.activateSelectedRow, "Open item / toggle folder"},
	}
}

func (a *App) BindKeys(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, a.armLeader); err != nil {
		return err
	}
	for _, b := range a.globalKeyBindings() {
		key := b.key
		if lc, ok := key.(leaderChord); ok {
			key = rune(lc)
		}
		if err := g.SetKeybinding("", key, gocui.ModNone, b.fn); err != nil {
			return err
		}
	}

	for _, viewname := range []string{a.Categories.WrapperName(), a.Categories.ContentViewName()} {
		if err := g.SetKeybinding(viewname, gocui.MouseLeft, gocui.ModNone, a.focusCategories); err != nil {
			return err
		}
	}
	for _, b := range a.categoryKeyBindings() {
		if err := g.SetKeybinding(a.Categories.WrapperName(), b.key, gocui.ModNone, b.fn); err != nil {
			return err
		}
	}

	for _, viewname := range []string{a.Content.WrapperName(), a.Content.ContentViewName(), a.Content.ScrollViewName()} {
		if err := g.SetKeybinding(viewname, gocui.MouseLeft, gocui.ModNone, a.focusItems); err != nil {
			return err
		}
	}
	for _, b := range a.itemKeyBindings() {
		if err := g.SetKeybinding(a.Content.WrapperName(), b.key, gocui.ModNone, b.fn); err != nil {
			return err
		}
	}

	if err := a.Preview.BindKeys(g, a.closePreview); err != nil {
		return err
	}
	if err := a.Help.BindKeys(g, a.closeHelp); err != nil {
		return err
	}
	if err := a.ConfigForm.BindKeys(g, a.submitConfigForm, a.cancelConfigForm); err != nil {
		return err
	}

	return nil
}

func (a *App) modalOpen() bool {
	return a.Preview.IsOpen() || a.Help.IsOpen() || a.ConfigForm.IsOpen()
}

func (a *App) quit(g *gocui.Gui, v *gocui.View) error {
	if a.modalOpen() {
		return nil
	}
	return gocui.ErrQuit
}

func (a *App) focusCategories(g *gocui.Gui, v *gocui.View) error {
	if a.modalOpen() || !a.categoriesVisible {
		return nil
	}
	a.focused = a.Categories.WrapperName()
	_, err := g.SetCurrentView(a.focused)
	return err
}

// armLeader starts the <leader> chord's timeout window on Space, Vim's
// default leader key. It's registered directly rather than through
// globalKeyBindings so it stays out of the help modal - by itself it does
// nothing a user would want listed as a binding.
func (a *App) armLeader(g *gocui.Gui, v *gocui.View) error {
	a.leaderArmedAt = time.Now()
	return nil
}

// consumeLeader reports whether Space was pressed within leaderTimeout of
// now, and clears the pending state either way so a chord can't be
// completed twice or after the window closes.
func (a *App) consumeLeader() bool {
	armed := !a.leaderArmedAt.IsZero() && time.Since(a.leaderArmedAt) <= leaderTimeout
	a.leaderArmedAt = time.Time{}
	return armed
}

// toggleCategories shows or hides the Categories pane so Content can
// expand to fill the freed width. modalOpen() guards against leaving a
// modal's previousFocused pointing at a view we're about to delete. The
// synchronous a.layout(g) call is required: run.go's manager renders
// before it lays out, so without this the re-shown pane would draw one
// empty frame before its views exist.
func (a *App) toggleCategories(g *gocui.Gui, v *gocui.View) error {
	if !a.consumeLeader() || a.modalOpen() {
		return nil
	}
	a.categoriesVisible = !a.categoriesVisible
	if !a.categoriesVisible {
		a.Categories.Delete(g)
		if a.focused == a.Categories.WrapperName() {
			a.focused = a.Content.WrapperName()
		}
	}

	if g == nil {
		return nil
	}
	if err := a.layout(g); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("categories: unable to focus view:", err)
	}
	return nil
}

func (a *App) focusItems(g *gocui.Gui, v *gocui.View) error {
	if a.modalOpen() {
		return nil
	}
	a.focused = a.Content.WrapperName()
	_, err := g.SetCurrentView(a.focused)
	return err
}

func (a *App) categoryMoveDown(g *gocui.Gui, v *gocui.View) error {
	if a.CategorySelected < len(categoryOrder)-1 {
		a.CategorySelected++
		a.RowSelected = 0
	}
	return nil
}

func (a *App) categoryMoveUp(g *gocui.Gui, v *gocui.View) error {
	if a.CategorySelected > 0 {
		a.CategorySelected--
		a.RowSelected = 0
	}
	return nil
}

func (a *App) rowMoveDown(g *gocui.Gui, v *gocui.View) error {
	if a.RowSelected < len(a.visibleRows())-1 {
		a.RowSelected++
	}
	return nil
}

func (a *App) rowMoveUp(g *gocui.Gui, v *gocui.View) error {
	if a.RowSelected > 0 {
		a.RowSelected--
	}
	return nil
}

func (a *App) activateSelectedRow(g *gocui.Gui, v *gocui.View) error {
	r, ok := a.selectedRow()
	if !ok {
		return nil
	}
	if r.kind == rowKindItem {
		return a.openSelectedInPreview(g, v)
	}
	a.collapsed[r.path] = !a.collapsed[r.path]
	return nil
}

func (a *App) rescan(g *gocui.Gui, v *gocui.View) error {
	if a.modalOpen() {
		return nil
	}
	items, err := a.Service.Scan()
	if err != nil {
		log.Println("scan:", err)
		return nil
	}
	a.setItems(items)
	return nil
}
