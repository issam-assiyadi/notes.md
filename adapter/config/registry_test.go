package config_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/issam-assiyadi/leftmark/adapter/config"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	want := config.Registry{
		Projects: map[string]config.ProjectConfig{
			"/repo": {Ignore: []string{"vendor/**", "*.gen.go"}},
		},
	}

	if err := config.Save(path, want); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load = %+v, want %+v", got, want)
	}
}

func TestLoadMissingFileReturnsEmptyRegistry(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.json")

	got, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.Projects == nil || len(got.Projects) != 0 {
		t.Errorf("Load = %+v, want empty registry", got)
	}
}

func TestLookup(t *testing.T) {
	reg := config.Registry{
		Projects: map[string]config.ProjectConfig{
			"/foo/bar":     {Ignore: []string{"outer"}},
			"/foo/bar/baz": {Ignore: []string{"inner"}},
		},
	}

	tests := []struct {
		name      string
		cwd       string
		wantRoot  string
		wantFound bool
	}{
		{"exact match on shallower root", "/foo/bar", "/foo/bar", true},
		{"subdirectory picks longest match", "/foo/bar/baz/qux", "/foo/bar/baz", true},
		{"subdirectory of shallower root only", "/foo/bar/other", "/foo/bar", true},
		{"sibling with shared prefix does not match", "/foo/barbaz", "", false},
		{"unrelated path does not match", "/elsewhere", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, _, found := config.Lookup(reg, tt.cwd)
			if found != tt.wantFound || root != tt.wantRoot {
				t.Errorf("Lookup(%q) = (%q, %v), want (%q, %v)", tt.cwd, root, found, tt.wantRoot, tt.wantFound)
			}
		})
	}
}
