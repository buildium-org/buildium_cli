package templates

import (
	"io/fs"
	"strings"
	"testing"
)

// readAllContents concatenates every file in fsys into one string, so tests can
// assert that a token actually appears somewhere in a template.
func readAllContents(t *testing.T, fsys fs.FS) string {
	t.Helper()
	var b strings.Builder
	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		b.Write(data)
		b.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}
	return b.String()
}

func TestCatalogResolvesAndTokensExist(t *testing.T) {
	cat := Catalog()
	if len(cat) != 3 {
		t.Fatalf("expected 3 templates, got %d", len(cat))
	}

	for _, tmpl := range cat {
		if tmpl.Key == "" || tmpl.Label == "" {
			t.Errorf("template %+v has empty key or label", tmpl)
		}

		fsys, err := tmpl.FS()
		if err != nil {
			t.Fatalf("template %q FS(): %v", tmpl.Key, err)
		}
		contents := readAllContents(t, fsys)

		if len(tmpl.Fields) == 0 {
			t.Errorf("template %q declares no fields", tmpl.Key)
		}

		seen := map[string]bool{}
		for _, f := range tmpl.Fields {
			if f.Key == "" || f.Token == "" || f.Label == "" {
				t.Errorf("template %q field %+v has empty key/token/label", tmpl.Key, f)
			}
			if seen[f.Key] {
				t.Errorf("template %q has duplicate field key %q", tmpl.Key, f.Key)
			}
			seen[f.Key] = true

			// The catalog must match reality: each declared token has to occur
			// in at least one of the template's vendored files, otherwise the
			// wizard would prompt for a value that gets substituted nowhere.
			if !strings.Contains(contents, f.Token) {
				t.Errorf("template %q field %q: token %q not present in any template file",
					tmpl.Key, f.Key, f.Token)
			}
		}
	}
}

func TestByKey(t *testing.T) {
	for _, key := range []string{"tutorial", "go", "ts"} {
		if _, ok := ByKey(key); !ok {
			t.Errorf("ByKey(%q): not found", key)
		}
	}
	if _, ok := ByKey("does-not-exist"); ok {
		t.Error("ByKey(\"does-not-exist\"): expected not found")
	}
}
