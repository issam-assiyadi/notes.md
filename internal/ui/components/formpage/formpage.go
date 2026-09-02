// Package formpage is a full-page form with two editable text fields and
// Save/Cancel actions, used to let the user register or edit a project's
// root and ignore globs. It takes over the whole terminal rather than
// overlaying the categories/items panes, since project config is visited
// deliberately and rarely, not glanced at mid-browse. Its border, title,
// and sticky action row are drawn by actionview; the fields and actions
// themselves are the reusable textinput and button components. Render
// must never touch either field's own buffer, since those hold live,
// in-progress user input redrawn only by gocui's own editing, not by this
// page's per-frame Render.
package formpage

import (
	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/internal/ui/components/actionview"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/button"
	"github.com/issam-assiyadi/leftmark/internal/ui/components/textinput"
)

type Config struct {
	BaseName string
	// ConsumeChord is passed straight through to both the Root and
	// Ignore fields' own textinput.Config - see its doc comment.
	ConsumeChord func(g *gocui.Gui, v *gocui.View, ch rune) bool
}

type field int

const (
	fieldRoot field = iota
	fieldIgnore
	fieldSave
	fieldCancel
	fieldCount
)

type Page struct {
	baseName string

	av     *actionview.View
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

func New(cfg Config) *Page {
	base := cfg.BaseName
	return &Page{
		baseName: base,
		av:       actionview.New(actionview.Config{BaseName: base}),
		root: textinput.New(textinput.Config{
			BaseName:     base + "-root",
			Title:        " Root (absolute path) ",
			ConsumeChord: cfg.ConsumeChord,
		}),
		ignore: textinput.New(textinput.Config{
			BaseName:     base + "-ignore",
			Title:        " Ignore (comma-separated globs) ",
			ConsumeChord: cfg.ConsumeChord,
		}),
		save: button.New(button.Config{
			BaseName: base + "-save",
			Caption:  "Save (<leader>s)",
			Color:    gocui.ColorCyan,
		}),
		cancel: button.New(button.Config{
			BaseName: base + "-cancel",
			Caption:  "Cancel (<leader>q)",
			Color:    gocui.ColorRed,
		}),
	}
}

func (p *Page) IsOpen() bool { return p.open }

func (p *Page) errorLineName() string { return p.baseName + "-error" }

// stops lists the page's focusable views in Tab order.
func (p *Page) stops() []string {
	return []string{p.root.Name(), p.ignore.Name(), p.save.Name(), p.cancel.Name()}
}

// Open switches the page to its open state, pre-filling both fields, and
// returns the view name that should receive focus first (the Root field).
// currentFocused is remembered so Close can restore it.
func (p *Page) Open(g *gocui.Gui, currentFocused, title, rootValue, ignoreValue string) string {
	p.title = title
	p.previousFocused = currentFocused
	p.rootValue = rootValue
	p.ignoreValue = ignoreValue
	p.errorText = ""
	p.focus = fieldRoot
	p.open = true
	if g != nil {
		g.Cursor = true // Root, the initial focus, is always a text field.
	}
	return p.root.Name()
}

// Close switches the page to its closed state, deletes all of its views
// (Layout recreates them fresh, with a fresh pre-fill, next time Open is
// called), and returns the view name to restore focus to.
func (p *Page) Close(g *gocui.Gui) string {
	p.open = false
	if g != nil {
		g.Cursor = false
		p.root.Delete(g)
		p.ignore.Delete(g)
		p.save.Delete(g)
		p.cancel.Delete(g)
		p.av.Delete(g)
		_ = g.DeleteView(p.errorLineName())
	}
	return p.previousFocused
}

// Values reads back the current contents of the two text fields directly
// from their gocui buffers.
func (p *Page) Values(g *gocui.Gui) (root, ignore string) {
	return p.root.Value(g), p.ignore.Value(g)
}

// SetError sets the message the next Render call will show on the form's
// error line, without closing the page. Pass "" to clear it.
func (p *Page) SetError(text string) { p.errorText = text }

// HelpEntry is the key/description pair this page exposes for help
// screens.
type HelpEntry struct {
	Key  interface{}
	Desc string
}

// HelpEntries advertises <leader>s/<leader>q as plain strings rather than
// the app-level leaderChord type, since that type lives in package ui and
// importing it here would cycle back to this package. formatKey's
// default case renders any string key verbatim, so "<leader>s" already
// prints correctly without it. Ctrl+S and Esc still work too (see
// BindKeys) - kept as an unadvertised fallback for while a field is
// being actively typed into, since gocui defers rune keybindings like
// <leader>s's second key to the field's own Editor there.
func (p *Page) HelpEntries() []HelpEntry {
	return []HelpEntry{
		{gocui.KeyTab, "Next field"},
		{gocui.KeyBacktab, "Previous field"},
		{gocui.KeyEnter, "Save"},
		{"<leader>s", "Save"},
		{"<leader>q", "Exit to main view"},
	}
}
