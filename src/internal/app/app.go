package app

import (
	"log"

	"github.com/jroimartin/gocui"
	"portoflio.com/cli/internal/content"
)

type App struct {
	Pages       content.Pages
	CurrentPage int

	// I don't understand yet the role of those attribues
	SidebarView string
	ContentView string
}

func New(pages content.Pages) *App {
	return &App{
		Pages:       pages,
		CurrentPage: 0,

		SidebarView: "sidebar",
		ContentView: "content",
	}
}

func (app *App) Run() error {

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.SetManagerFunc(app.layout)

	// TODO: define keybinding in app/keybinding, this is just a minimal one.
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error { return gocui.ErrQuit }); err != nil {
		return err
	}

	if _, err := g.SetCurrentView(app.SidebarView); err != nil {
		log.Println("unable to set initial view:", err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
