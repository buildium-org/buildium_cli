package generator

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"buildium_cli/internal/templates"
)

// leftoverToken matches any unsubstituted <..._HERE> placeholder.
var leftoverToken = regexp.MustCompile(`<[A-Z_]+_HERE>`)

func mustTemplate(t *testing.T, key string) templates.Template {
	t.Helper()
	tmpl, ok := templates.ByKey(key)
	if !ok {
		t.Fatalf("template %q not in catalog", key)
	}
	return tmpl
}

// fullValues fills every field of a template with a recognizable value.
func fullValues(tmpl templates.Template) map[string]string {
	v := map[string]string{}
	for _, f := range tmpl.Fields {
		v[f.Key] = "val-" + f.Key
	}
	return v
}

func TestGenerateAllTemplatesLeaveNoTokens(t *testing.T) {
	for _, tmpl := range templates.Catalog() {
		t.Run(tmpl.Key, func(t *testing.T) {
			dest := filepath.Join(t.TempDir(), "out")
			if err := Generate(tmpl, fullValues(tmpl), dest); err != nil {
				t.Fatalf("Generate: %v", err)
			}

			var files int
			err := filepath.WalkDir(dest, func(p string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return err
				}
				files++
				// No .tmpl suffix should survive into the output tree.
				if strings.HasSuffix(p, templates.TmplSuffix) {
					t.Errorf("output retains .tmpl suffix: %s", p)
				}
				data, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				if m := leftoverToken.Find(data); m != nil {
					t.Errorf("%s still contains placeholder %q", p, m)
				}
				return nil
			})
			if err != nil {
				t.Fatalf("walk output: %v", err)
			}
			if files == 0 {
				t.Fatal("generated no files")
			}
		})
	}
}

func TestGenerateGoSubstitutesKnownFiles(t *testing.T) {
	tmpl := mustTemplate(t, "go")
	dest := filepath.Join(t.TempDir(), "proj")
	vals := map[string]string{
		"image_name":    "my-app",
		"harness_image": "my-tut_harness",
	}
	if err := Generate(tmpl, vals, dest); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// go.mod.tmpl -> go.mod, main.go.tmpl -> main.go
	for _, f := range []string{"go.mod", "main.go", ".gitignore", "meta.json"} {
		if _, err := os.Stat(filepath.Join(dest, f)); err != nil {
			t.Errorf("expected output file %q: %v", f, err)
		}
	}

	checks := map[string][]string{
		"Dockerfile":         {"my-tut_harness"},
		"docker-compose.yml": {"my-app"},
		"Makefile":           {"my-app"},
	}
	for file, wants := range checks {
		data, err := os.ReadFile(filepath.Join(dest, file))
		if err != nil {
			t.Fatalf("read %s: %v", file, err)
		}
		for _, w := range wants {
			if !strings.Contains(string(data), w) {
				t.Errorf("%s missing substituted value %q", file, w)
			}
		}
	}
}

func TestGenerateRestoresShebangExecutable(t *testing.T) {
	tmpl := mustTemplate(t, "ts")
	dest := filepath.Join(t.TempDir(), "proj")
	if err := Generate(tmpl, fullValues(tmpl), dest); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	info, err := os.Stat(filepath.Join(dest, "app"))
	if err != nil {
		t.Fatalf("stat app: %v", err)
	}
	if info.Mode()&0o111 == 0 {
		t.Errorf("app script is not executable (mode %v)", info.Mode())
	}

	// package.json carries the image name token too.
	data, err := os.ReadFile(filepath.Join(dest, "package.json"))
	if err != nil {
		t.Fatalf("read package.json: %v", err)
	}
	if !strings.Contains(string(data), "val-image_name") {
		t.Errorf("package.json name not substituted: %s", data)
	}
}

func TestGenerateRefusesNonEmptyDestination(t *testing.T) {
	tmpl := mustTemplate(t, "go")
	dest := t.TempDir() // exists and we put a file in it
	if err := os.WriteFile(filepath.Join(dest, "keep.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Generate(tmpl, fullValues(tmpl), dest)
	if !errors.Is(err, ErrDestinationNotEmpty) {
		t.Fatalf("expected ErrDestinationNotEmpty, got %v", err)
	}
}

func TestGenerateCreatesMissingDestination(t *testing.T) {
	tmpl := mustTemplate(t, "tutorial")
	dest := filepath.Join(t.TempDir(), "nested", "new-tutorial")
	if err := Generate(tmpl, fullValues(tmpl), dest); err != nil {
		t.Fatalf("Generate into missing dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dest, "Makefile")); err != nil {
		t.Errorf("expected Makefile in created dir: %v", err)
	}
}
