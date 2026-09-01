package ui

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/issam-assiyadi/leftmark"
	"github.com/issam-assiyadi/leftmark/adapter/config"
)

// consumeConfigFormChord is the Root/Ignore fields' ConsumeChord hook
// (wired in app.go's New, see textinput's doc comment): it lets
// <leader>s/<leader>q complete a save or exit even while a field has
// focus. armLeader already wrote a placeholder space into the view by
// the time this runs (see its own comment); consumeLeaderChordInView
// confirms that space landed in this exact view within the timeout
// before EditDelete retracts it, so an unrelated 's' or 'q' typed on its
// own - or one typed a beat too late - is left as ordinary text.
func (a *App) consumeConfigFormChord(g *gocui.Gui, v *gocui.View, ch rune) bool {
	if ch != 's' && ch != 'q' {
		return false
	}
	if v == nil || !a.consumeLeaderChordInView(v.Name()) {
		return false
	}
	v.EditDelete(true)
	if ch == 's' {
		_ = a.submitConfigForm(g, v)
	} else {
		_ = a.cancelConfigForm(g, v)
	}
	return true
}

// openConfigConfirm asks, via a small overlay, whether an unconfigured
// project should be configured now. It's shown once at startup instead of
// forcing the user straight into the config page - see run.go.
func (a *App) openConfigConfirm(g *gocui.Gui, v *gocui.View) error {
	if a.modalOpen() {
		return nil
	}
	a.focused = a.ConfirmConfig.Open(a.focused, " Project Config ", "This project isn't configured yet. Configure it now?")
	if g == nil {
		return nil
	}
	if err := a.layout(g); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("configconfirm: unable to focus view:", err)
	}
	return nil
}

// confirmConfigureYes closes the confirm overlay and opens the config
// page in its place, carrying its previousFocused forward so the config
// page's own Close still restores focus to wherever the user started.
func (a *App) confirmConfigureYes(g *gocui.Gui, v *gocui.View) error {
	a.focused = a.ConfirmConfig.Close(g)
	return a.openConfigForm(g, v)
}

// confirmConfigureSkip dismisses the confirm overlay and leaves the user
// browsing with whatever was already scanned from their current
// directory - see tui/leftmark/main.go's unregistered fallback.
func (a *App) confirmConfigureSkip(g *gocui.Gui, v *gocui.View) error {
	a.focused = a.ConfirmConfig.Close(g)
	if g == nil {
		return nil
	}
	if err := a.layout(g); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("configconfirm: unable to restore focus:", err)
	}
	return nil
}

// openConfigForm switches to the full-page config form. Categories and
// Content are deleted rather than merely left unlaid-out, matching
// scrollview.View.Delete's documented requirement for a pane that comes
// and goes: gocui keeps drawing already-registered views regardless of
// whether Layout is still called for them.
func (a *App) openConfigForm(g *gocui.Gui, v *gocui.View) error {
	if a.modalOpen() {
		return nil
	}
	a.Categories.Delete(g)
	a.Content.Delete(g)
	a.focused = a.ConfigForm.Open(g, a.focused, " Project Config ", a.projectRoot, strings.Join(a.ignorePatterns, ", "))
	if g == nil {
		return nil
	}
	if err := a.layout(g); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("configform: unable to focus view:", err)
	}
	return nil
}

// saveConfigForm is <leader>s's handler, registered globally (see
// keybindings.go's BindKeys) since it needs to reach the page regardless
// of which of its fields or buttons currently has focus. It defers to
// Help when Help is open on top of the config page, so a stray 's' meant
// for reading the help screen can't accidentally save.
func (a *App) saveConfigForm(g *gocui.Gui, v *gocui.View) error {
	if !a.consumeLeader() || !a.ConfigForm.IsOpen() || a.Help.IsOpen() {
		return nil
	}
	return a.submitConfigForm(g, v)
}

func (a *App) cancelConfigForm(g *gocui.Gui, v *gocui.View) error {
	if !a.ConfigForm.IsOpen() {
		return nil
	}
	a.focused = a.ConfigForm.Close(g)
	if g == nil {
		return nil
	}
	if err := a.layout(g); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("configform: unable to restore focus:", err)
	}
	return nil
}

func parseIgnoreList(text string) []string {
	var ignore []string
	for _, pattern := range strings.Split(text, ",") {
		if pattern = strings.TrimSpace(pattern); pattern != "" {
			ignore = append(ignore, pattern)
		}
	}
	return ignore
}

func (a *App) submitConfigForm(g *gocui.Gui, v *gocui.View) error {
	rootText, ignoreText := a.ConfigForm.Values(g)

	root := filepath.Clean(rootText)
	if !filepath.IsAbs(root) {
		a.ConfigForm.SetError("root must be an absolute path")
		return nil
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		a.ConfigForm.SetError("root does not exist or is not a directory")
		return nil
	}

	ignore := parseIgnoreList(ignoreText)

	if a.registry.Projects == nil {
		a.registry.Projects = map[string]config.ProjectConfig{}
	}
	if a.registered && a.projectRoot != root {
		delete(a.registry.Projects, a.projectRoot)
	}
	a.registry.Projects[root] = config.ProjectConfig{Ignore: ignore}

	if err := config.Save(a.registryPath, a.registry); err != nil {
		a.ConfigForm.SetError("save failed: " + err.Error())
		return nil
	}

	a.Service = leftmark.New(root, ignore...)
	a.projectRoot, a.ignorePatterns, a.registered = root, ignore, true

	items, err := a.Service.Scan()
	if err != nil {
		a.ConfigForm.SetError("scan failed: " + err.Error())
		return nil
	}
	a.setItems(items)

	a.focused = a.ConfigForm.Close(g)
	if g == nil {
		return nil
	}
	if err := a.layout(g); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("configform: unable to restore focus:", err)
	}
	return nil
}
