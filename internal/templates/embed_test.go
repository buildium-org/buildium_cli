package templates

import (
	"io/fs"
	"testing"
)

// wantFiles lists a few representative paths that must be present in each
// embedded template root. It deliberately includes files that go:embed would
// drop without the `all:` prefix (.gitignore) and a nested go.mod, which is the
// case most likely to be silently excluded.
var wantFiles = map[string][]string{
	"tutorial": {"Dockerfile", "Makefile", "go.mod.tmpl", "main.go.tmpl", "manifest/info.json", "steps/steps.go.tmpl"},
	"go":       {"Dockerfile", "Makefile", "go.mod.tmpl", "main.go.tmpl", "meta.json", ".gitignore"},
	"ts":       {"Dockerfile", "Makefile", "package.json", "meta.json", "src/main.ts", ".gitignore"},
}

func TestEmbeddedTemplatesPresent(t *testing.T) {
	for key, want := range wantFiles {
		sub, err := Sub(key)
		if err != nil {
			t.Fatalf("Sub(%q): %v", key, err)
		}

		got := map[string]bool{}
		if err := fs.WalkDir(sub, ".", func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				got[p] = true
			}
			return nil
		}); err != nil {
			t.Fatalf("walk %q: %v", key, err)
		}

		if len(got) == 0 {
			t.Fatalf("template %q embedded no files", key)
		}
		for _, f := range want {
			if !got[f] {
				t.Errorf("template %q missing embedded file %q (walk found %d files)", key, f, len(got))
			}
		}
	}
}
