// Package formmodal is a small floating overlay with two editable text
// fields and Save/Cancel actions, used to let the user register or edit a
// project's root and ignore globs. It's built from the reusable textinput
// and button components; Render must never touch either field's own
// buffer, since those hold live, in-progress user input redrawn only by
// gocui's own editing, not by this modal's per-frame Render.
package formmodal

import (
	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/internal/ui/components/button"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/textinput"
)

type Config struct {
	BaseName string
}

type field int

const (
	fieldRoot field = iota
	fieldIgnore
	fieldSave
	fieldCancel
	fieldCount
)

type Modal struct {
	baseName string

	root   *textinput.Input
	ignore *textinput.Input
	save   *button.Button
	cancel *button.Button

	open            bool
	title           string
	previousFocused string
	focus           field

	rootValue   string
	ignoreValue string
	errorText   string
}

func New(cfg Config) *Modal {
	base := cfg.BaseName
	return &Modal{
		baseName: base,
		root: textinput.New(textinput.Config{
			BaseName: base + "-root",
			Title:    " Root (absolute path) ",
		}),
		ignore: textinput.New(textinput.Config{
			BaseName: base + "-ignore",
			Title:    " Ignore (comma-separated globs) ",
		}),
		save: button.New(button.Config{
			BaseName: base + "-save",
			Caption:  "Save",
			Color:    gocui.ColorCyan,
		}),
		cancel: button.New(button.Config{
			BaseName: base + "-cancel",
			Caption:  "Cancel",
			Color:    gocui.ColorRed,
		}),
	}
}

func (m *Modal) IsOpen() bool { return m.open }

func (m *Modal) wrapperName() string   { return m.baseName + "-wrapper" }
func (m *Modal) errorLineName() string { return m.baseName + "-error" }

// stops lists the modal's focusable views in Tab order.
func (m *Modal) stops() []string {
	return []string{m.root.Name(), m.ignore.Name(), m.save.Name(), m.cancel.Name()}
}

// Open switches the modal to its open state, pre-filling both fields, and
// returns the view name that should receive focus first (the Root field).
// currentFocused is remembered so Close can restore it.
func (m *Modal) Open(g *gocui.Gui, currentFocused, title, rootValue, ignoreValue string) string {
	m.title = title
	m.previousFocused = currentFocused
	m.rootValue = rootValue
	m.ignoreValue = ignoreValue
	m.errorText = ""
	m.focus = fieldRoot
	m.open = true
	g.Cursor = true // Root, the initial focus, is always a text field.
	return m.root.Name()
}

// Close switches the modal to its closed state, deletes all of its views
// (Layout recreates them fresh, with a fresh pre-fill, next time Open is
// called), and returns the view name to restore focus to.
func (m *Modal) Close(g *gocui.Gui) string {
	m.open = false
	g.Cursor = false
	m.root.Delete(g)
	m.ignore.Delete(g)
	m.save.Delete(g)
	m.cancel.Delete(g)
	_ = g.DeleteView(m.wrapperName())
	_ = g.DeleteView(m.errorLineName())
	return m.previousFocused
}

// Values reads back the current contents of the two text fields directly
// from their gocui buffers.
func (m *Modal) Values(g *gocui.Gui) (root, ignore string) {
	return m.root.Value(g), m.ignore.Value(g)
}

// SetError sets the message the next Render call will show on the form's
// error line, without closing the form. Pass "" to clear it.
func (m *Modal) SetError(text string) { m.errorText = text }

// HelpEntry is the key/description pair this modal exposes for help
// screens.
type HelpEntry struct {
	Key  interface{}
	Desc string
}

func (m *Modal) HelpEntries() []HelpEntry {
	return []HelpEntry{
		{gocui.KeyTab, "Next field"},
		{gocui.KeyBacktab, "Previous field"},
		{gocui.KeyEnter, "Save"},
		{gocui.KeyCtrlS, "Save"},
		{gocui.KeyEsc, "Cancel"},
	}
}
