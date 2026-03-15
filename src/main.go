package main

import (
	"log"

	"notesmd.com/cli/internal/app"
	"notesmd.com/cli/internal/content"
)

func main() {
	pages, err := content.LoadPages("./pages")
	if err != nil {
		log.Fatalf("load pages: %v", err)
	}

	a := app.New(pages)

	if err := a.Run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}
