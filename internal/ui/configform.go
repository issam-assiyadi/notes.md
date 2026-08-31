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

func (a *App) openConfigForm(g *gocui.Gui, v *gocui.View) error {
	if a.modalOpen() {
		return nil
	}
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

func (a *App) cancelConfigForm(g *gocui.Gui, v *gocui.View) error {
	if !a.ConfigForm.IsOpen() {
		return nil
	}
	a.focused = a.ConfigForm.Close(g)
	_, err := g.SetCurrentView(a.focused)
	return err
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
	_, err = g.SetCurrentView(a.focused)
	return err
}
