// Package confirmmodal is a small floating overlay with a message and two
// action buttons, used to ask the user a yes/no question without leaving
// whatever they were looking at. Built from the reusable button component.
package confirmmodal

import (
	"github.com/awesome-gocui/gocui"

	"github.com/issam-assiyadi/leftmark/internal/ui/components/button"
)

type Config struct {
	BaseName     string
	ConfirmLabel string
	DismissLabel string
}

type field int

const (
	fieldConfirm field = iota
	fieldDismiss
	fieldCount
)

type Modal struct {
	baseName string

	confirm *button.Button
	dismiss *button.Button

	open            bool
	title           string
	message         string
	previousFocused string
	focus           field
}

func New(cfg Config) *Modal {
	base := cfg.BaseName
	return &Modal{
		baseName: base,
		confirm: button.New(button.Config{
			BaseName: base + "-confirm",
			Caption:  cfg.ConfirmLabel,
			Color:    gocui.ColorCyan,
		}),
		dismiss: button.New(button.Config{
			BaseName: base + "-dismiss",
			Caption:  cfg.DismissLabel,
			Color:    gocui.ColorYellow,
		}),
	}
}

func (m *Modal) IsOpen() bool { return m.open }

func (m *Modal) wrapperName() string     { return m.baseName + "-wrapper" }
func (m *Modal) messageViewName() string { return m.baseName + "-message" }

// stops lists the modal's focusable views in Tab order.
func (m *Modal) stops() []string {
	return []string{m.confirm.Name(), m.dismiss.Name()}
}

// Open switches the modal to its open state and returns the view name
// that should receive focus first (the confirm button). currentFocused
// is remembered so Close can restore it.
func (m *Modal) Open(currentFocused, title, message string) string {
	m.title = title
	m.message = message
	m.previousFocused = currentFocused
	m.focus = fieldConfirm
	m.open = true
	return m.confirm.Name()
}

// Close switches the modal to its closed state, deletes all of its views
// (Layout recreates them fresh next time Open is called), and returns the
// view name to restore focus to.
func (m *Modal) Close(g *gocui.Gui) string {
	m.open = false
	if g != nil {
		m.confirm.Delete(g)
		m.dismiss.Delete(g)
		_ = g.DeleteView(m.wrapperName())
		_ = g.DeleteView(m.messageViewName())
	}
	return m.previousFocused
}

// HelpEntry is the key/description pair this modal exposes for help
// screens.
type HelpEntry struct {
	Key  interface{}
	Desc string
}

func (m *Modal) HelpEntries() []HelpEntry {
	return []HelpEntry{
		{gocui.KeyTab, "Next option"},
		{gocui.KeyBacktab, "Previous option"},
		{gocui.KeyEnter, "Choose"},
		{gocui.KeyEsc, "Dismiss"},
	}
}
