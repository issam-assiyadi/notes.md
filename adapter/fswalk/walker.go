// Package fswalk implements application.FileWalker over the real
// filesystem.
package fswalk

import (
	"io/fs"
	"path/filepath"

	gitignore "github.com/sabhiram/go-gitignore"
)

// Walker walks a directory tree, honoring the root .gitignore (if any),
// extra ignore globs supplied at construction, and always skipping .git
// regardless of ignore rules. It only consults a root-level .gitignore, not
// git's full per-directory cascading semantics - enough for a personal dev
// tool without reimplementing git's ignore rules.
type Walker struct {
	extraIgnore []string
}

// New builds a Walker. extraIgnore is a set of additional .gitignore-style
// glob patterns applied on top of the root .gitignore.
func New(extraIgnore ...string) Walker { return Walker{extraIgnore: extraIgnore} }

func (w Walker) Walk(root string, fn func(path string) error) error {
	rootIgnore, _ := gitignore.CompileIgnoreFile(filepath.Join(root, ".gitignore"))

	var extraIgnore *gitignore.GitIgnore
	if len(w.extraIgnore) > 0 {
		extraIgnore = gitignore.CompileIgnoreLines(w.extraIgnore...)
	}

	matches := func(rel string) bool {
		return (rootIgnore != nil && rootIgnore.MatchesPath(rel)) ||
			(extraIgnore != nil && extraIgnore.MatchesPath(rel))
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			rel = path
		}

		if d.IsDir() {
			switch {
			case rel == ".":
				return nil
			case d.Name() == ".git":
				return filepath.SkipDir
			case matches(rel):
				return filepath.SkipDir
			default:
				return nil
			}
		}

		if matches(rel) {
			return nil
		}

		return fn(path)
	})
}
