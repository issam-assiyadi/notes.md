package version

import "testing"

func TestWithVPrefix(t *testing.T) {
	cases := map[string]string{
		"0.1.0":                  "v0.1.0",
		"v0.1.0":                 "v0.1.0",
		"0.0.0-SNAPSHOT-13c6520": "v0.0.0-SNAPSHOT-13c6520",
		"dev":                    "dev",
		"":                       "",
	}
	for in, want := range cases {
		if got := withVPrefix(in); got != want {
			t.Errorf("withVPrefix(%q) = %q, want %q", in, got, want)
		}
	}
}
