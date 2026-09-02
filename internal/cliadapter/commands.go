package cliadapter

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/issam-assiyadi/leftmark"
	"github.com/issam-assiyadi/leftmark/adapter/config"
	"github.com/issam-assiyadi/leftmark/adapter/githook"
	"github.com/issam-assiyadi/leftmark/application"
	"github.com/issam-assiyadi/leftmark/domain"
)

var scanCategories = []domain.Kind{domain.KindTODO, domain.KindFIXME, domain.KindNOTE, domain.KindQUESTION}

// parseCategory resolves a --category flag value to a domain.Kind,
// case-insensitively. An empty input means "no filter" and is always valid.
func parseCategory(s string) (kind domain.Kind, ok bool) {
	if s == "" {
		return "", true
	}
	upper := domain.Kind(strings.ToUpper(s))
	for _, k := range scanCategories {
		if k == upper {
			return k, true
		}
	}
	return "", false
}

func filterByKind(items []domain.Item, kind domain.Kind) []domain.Item {
	var filtered []domain.Item
	for _, item := range items {
		if item.Kind == kind {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func newService(root string) (*application.Service, error) {
	explicitRoot := root != ""
	if !explicitRoot {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		root = wd
	}

	var ignore []string
	if path, err := config.DefaultPath(); err == nil {
		if reg, err := config.Load(path); err == nil {
			if matchedRoot, proj, found := config.Lookup(reg, root); found {
				ignore = proj.Ignore
				// An explicit -root is the user's final word on where to
				// scan; only a root defaulted from cwd gets resolved up
				// to its registered ancestor, matching the TUI.
				if !explicitRoot {
					root = matchedRoot
				}
			}
		}
	}

	return leftmark.New(root, ignore...), nil
}

func printItems(items []domain.Item, asJSON bool, stdout io.Writer) int {
	if asJSON {
		out := make([]itemJSON, len(items))
		for i, item := range items {
			out[i] = toItemJSON(item)
		}
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			return 1
		}
		return 0
	}

	for _, item := range items {
		printf(stdout, "%-9s %s:%d %s\n", item.Kind, item.File, item.Line, item.Text)
	}
	return 0
}

func runScan(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", "", "repo root (default: current directory)")
	jsonOut := fs.Bool("json", false, "print JSON")
	category := fs.String("category", "", "filter by category: TODO, FIXME, NOTE, or QUESTION")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	kind, ok := parseCategory(*category)
	if !ok {
		printf(stderr, "unknown category %q, want one of TODO, FIXME, NOTE, QUESTION\n", *category)
		return 2
	}

	svc, err := newService(*root)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	items, err := svc.Scan()
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	if kind != "" {
		items = filterByKind(items, kind)
	}
	return printItems(items, *jsonOut, stdout)
}

func runReport(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", "", "repo root (default: current directory)")
	jsonOut := fs.Bool("json", false, "print JSON")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	svc, err := newService(*root)
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	items, err := svc.Scan()
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	summary := application.Report(items)
	if *jsonOut {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(summary); err != nil {
			printf(stderr, "%v\n", err)
			return 1
		}
		return 0
	}

	printf(stdout, "%d items total\n", summary.Total)
	kinds := make([]string, 0, len(summary.ByKind))
	for k := range summary.ByKind {
		kinds = append(kinds, string(k))
	}
	sort.Strings(kinds)
	for _, k := range kinds {
		printf(stdout, "  %s: %d\n", k, summary.ByKind[domain.Kind(k)])
	}
	return 0
}

func runHook(args []string, stdout, stderr io.Writer) int {
	if len(args) < 2 || args[0] != "install" {
		printf(stderr, "usage: leftmark hook install <pre-commit|pre-push>\n")
		return 2
	}

	hookName := args[1]
	if hookName != "pre-commit" && hookName != "pre-push" {
		printf(stderr, "unknown hook %q, want pre-commit or pre-push\n", hookName)
		return 2
	}

	wd, err := os.Getwd()
	if err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}

	dir := filepath.Join(wd, "scripts", "git-hooks")
	if err := githook.Install(dir, hookName); err != nil {
		printf(stderr, "%v\n", err)
		return 1
	}
	printf(stdout, "installed %s in %s\n", hookName, dir)
	return 0
}
