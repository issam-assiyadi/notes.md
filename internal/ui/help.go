package ui

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/internal/ui/components/confirmmodal"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/formpage"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/modal"
)

type helpLine struct {
	key  interface{}
	desc string
}

type helpSection struct {
	title string
	lines []helpLine
}

func (a *App) helpSections() []helpSection {
	return []helpSection{
		{"Global", toHelpLines(a.globalKeyBindings())},
		{"Categories pane", toHelpLines(a.categoryKeyBindings())},
		{"Items pane", toHelpLines(a.itemKeyBindings())},
		{"Modal controls", fromModalEntries(a.Preview.HelpEntries())},
		{"Config form", fromFormPageEntries(a.ConfigForm.HelpEntries())},
		{"Config confirm", fromConfirmModalEntries(a.ConfirmConfig.HelpEntries())},
	}
}

func toHelpLines(bindings []appKeyBinding) []helpLine {
	lines := make([]helpLine, len(bindings))
	for i, b := range bindings {
		lines[i] = helpLine{key: b.key, desc: b.desc}
	}
	return lines
}

func fromModalEntries(entries []modal.HelpEntry) []helpLine {
	lines := make([]helpLine, len(entries))
	for i, e := range entries {
		lines[i] = helpLine{key: e.Key, desc: e.Desc}
	}
	return lines
}

func fromFormPageEntries(entries []formpage.HelpEntry) []helpLine {
	lines := make([]helpLine, len(entries))
	for i, e := range entries {
		lines[i] = helpLine{key: e.Key, desc: e.Desc}
	}
	return lines
}

func fromConfirmModalEntries(entries []confirmmodal.HelpEntry) []helpLine {
	lines := make([]helpLine, len(entries))
	for i, e := range entries {
		lines[i] = helpLine{key: e.Key, desc: e.Desc}
	}
	return lines
}

var keyLabels = map[gocui.Key]string{
	gocui.KeyArrowDown: "↓",
	gocui.KeyArrowUp:   "↑",
	gocui.KeyEnter:     "Enter",
	gocui.KeySpace:     "Space",
	gocui.KeyEsc:       "Esc",
	gocui.KeyPgdn:      "PgDn",
	gocui.KeyPgup:      "PgUp",
	gocui.KeyHome:      "Home",
	gocui.KeyEnd:       "End",
	gocui.KeyCtrlC:     "Ctrl+C",
	gocui.KeyTab:       "Tab",
	gocui.KeyBacktab:   "Shift+Tab",
	gocui.KeyCtrlS:     "Ctrl+S",
}

func formatKey(key interface{}) string {
	switch k := key.(type) {
	case rune:
		return string(k)
	case leaderChord:
		return "<leader>" + string(rune(k))
	case gocui.Key:
		if label, ok := keyLabels[k]; ok {
			return label
		}
		return fmt.Sprintf("Key(%d)", k)
	default:
		return fmt.Sprintf("%v", key)
	}
}

// renderHelp draws the help modal's sections. It's shared between the
// browse view's Render and the config page's Render (see render.go) since
// Help can open as an overlay on top of either.
func (a *App) renderHelp(g *gocui.Gui) error {
	return a.Help.Render(g, func(v *gocui.View, contentWidth int) error {
		for _, section := range a.helpSections() {
			_, _ = fmt.Fprintf(v, "%s\n", section.title)
			for _, line := range section.lines {
				_, _ = fmt.Fprintf(v, "  %-10s %s\n", formatKey(line.key), line.desc)
			}
			_, _ = fmt.Fprintln(v)
		}
		return nil
	})
}

// openHelp deliberately excludes ConfigForm from its guard: Help is meant
// to be reachable from the config page too (its own status bar always
// advertises "? Help"), as an overlay on top of it rather than a
// replacement for it.
func (a *App) openHelp(g *gocui.Gui, v *gocui.View) error {
	if a.Preview.IsOpen() || a.Help.IsOpen() || a.ConfirmConfig.IsOpen() {
		return nil
	}
	a.focused = a.Help.Open(a.focused, " Keybindings ", -1)

	if g == nil {
		return nil
	}
	if err := a.layout(g); err != nil {
		return err
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("help: unable to focus view:", err)
	}
	return nil
}

func (a *App) closeHelp(g *gocui.Gui, v *gocui.View) error {
	if !a.Help.IsOpen() {
		return nil
	}
	a.focused = a.Help.Close(g)

	if g == nil {
		return nil
	}
	if _, err := g.SetCurrentView(a.focused); err != nil {
		log.Println("help: unable to restore focus:", err)
	}
	return nil
}
