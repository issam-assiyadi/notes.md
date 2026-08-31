package ui

import (
	"log"

	"github.com/awesome-gocui/gocui"
)

type appKeyBinding struct {
	key  interface{}
	fn   func(*gocui.Gui, *gocui.View) error
	desc string
}

func (a *App) globalKeyBindings() []appKeyBinding {
	return []appKeyBinding{
		{gocui.KeyCtrlC, a.quit, "Quit"},
		{'q', a.quit, "Quit"},
		{'r', a.rescan, "Rescan"},
		{'1', a.focusCategories, "Focus categories pane"},
		{'2', a.focusItems, "Focus items pane"},
		{'c', a.toggleCategories, "Toggle categories pane"},
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
		{gocui.KeySpace, a.activateSelectedRow, "Open item / toggle folder"},
	}
}

func (a *App) BindKeys(g *gocui.Gui) error {
	for _, b := range a.globalKeyBindings() {
		if err := g.SetKeybinding("", b.key, gocui.ModNone, b.fn); err != nil {
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

	return nil
}

func (a *App) modalOpen() bool { return a.Preview.IsOpen() || a.Help.IsOpen() }

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

// toggleCategories shows or hides the Categories pane so Content can
// expand to fill the freed width. modalOpen() guards against leaving a
// modal's previousFocused pointing at a view we're about to delete. The
// synchronous a.layout(g) call is required: run.go's manager renders
// before it lays out, so without this the re-shown pane would draw one
// empty frame before its views exist.
func (a *App) toggleCategories(g *gocui.Gui, v *gocui.View) error {
	if a.modalOpen() {
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
