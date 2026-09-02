// Package version resolves the leftmark version string at runtime across
// two build paths: GoReleaser builds inject it via -ldflags, while a plain
// `go install pkg@vX.Y.Z` never passes ldflags but has the requested
// version stamped into the binary's build info automatically.
package version

import (
	"runtime/debug"
	"unicode"
)

// version is set via -ldflags "-X .../internal/version.version=X.Y.Z" by
// GoReleaser, using its {{.Version}} template value, which has no leading
// "v". It stays empty for a plain `go build` or local `go install`.
var version string

// Version returns the ldflags-injected value if set, otherwise the module
// version Go's build stamps into the binary for `go install pkg@vX.Y.Z`,
// otherwise "dev". The result always carries a "v" prefix when it's a
// version number, so it reads the same regardless of which build path
// produced the binary.
func Version() string {
	if version != "" {
		return withVPrefix(version)
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return "dev"
}

func withVPrefix(v string) string {
	if v == "" || v[0] == 'v' || !unicode.IsDigit(rune(v[0])) {
		return v
	}
	return "v" + v
}
