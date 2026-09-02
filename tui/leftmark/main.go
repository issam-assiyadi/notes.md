package main

import (
	"log"
	"os"

	"github.com/issam-assiyadi/leftmark"
	"github.com/issam-assiyadi/leftmark/adapter/config"
	"github.com/issam-assiyadi/leftmark/internal/cliadapter"
	"github.com/issam-assiyadi/leftmark/internal/ui"
)

func main() {
	if handled, code := cliadapter.Dispatch(os.Args[1:], os.Stdout, os.Stderr); handled {
		os.Exit(code)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("getwd: %v", err)
	}

	registryPath, err := config.DefaultPath()
	if err != nil {
		log.Fatalf("config path: %v", err)
	}
	registry, err := config.Load(registryPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	root, projectCfg, registered := config.Lookup(registry, cwd)
	if !registered {
		root = cwd
	}

	svc := leftmark.New(root, projectCfg.Ignore...)

	a, err := ui.New(svc, ui.StartupConfig{
		RegistryPath: registryPath,
		Registry:     registry,
		Root:         root,
		Ignore:       projectCfg.Ignore,
		Registered:   registered,
	})
	if err != nil {
		log.Fatalf("ui: %v", err)
	}

	if err := a.Run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}
