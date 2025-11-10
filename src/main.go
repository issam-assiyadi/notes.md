package main

import (
	"log"

	"portoflio.com/cli/internal/app"
	"portoflio.com/cli/internal/content"
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
