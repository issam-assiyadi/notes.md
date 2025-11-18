package app

import (
	"log"

	"github.com/jroimartin/gocui"
	"portoflio.com/cli/internal/content"
)

type App struct {
	Pages       content.Pages
	CurrentPage int

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

	g.SetManagerFunc(func(g *gocui.Gui) error {
		if err := app.layout(g); err != nil {
			return err
		}
		return app.Render(g)
	})

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
