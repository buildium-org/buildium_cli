package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"buildium_cli/internal/templates"
)

// ErrDestinationNotEmpty is returned when Generate is asked to write into a
// directory that already exists and contains files. Callers can test for it
// with errors.Is to give the user a friendly "pick another name" message.
var ErrDestinationNotEmpty = errors.New("destination already exists and is not empty")

// Generate materializes tmpl into destDir, replacing every field token with the
// value supplied in values (keyed by templates.Field.Key). It:
//
//   - refuses to write into an existing non-empty directory (ErrDestinationNotEmpty);
//   - recreates the template's directory structure under destDir;
//   - replaces tokens in both path segments and text file contents;
//   - strips the .tmpl suffix from output filenames (see templates.TmplSuffix);
//   - copies binary files byte-for-byte without attempting replacement.
//
// File modes: embed.FS reports every file as 0444, so source permissions can't
// be carried through. Files are written 0644, except text files beginning with
// a "#!" shebang, which are written 0755 so generated scripts stay executable.
func Generate(tmpl templates.Template, values map[string]string, destDir string) error {
	replacer := replacerFor(tmpl, values)

	if err := ensureEmptyDir(destDir); err != nil {
		return err
	}

	fsys, err := tmpl.FS()
	if err != nil {
		return fmt.Errorf("load template %q: %w", tmpl.Key, err)
	}

	return fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if p == "." {
			return nil
		}

		// Token-substitute the path, then drop the .tmpl suffix from its final
		// element, so e.g. "main.go.tmpl" -> "main.go".
		outRel := strings.TrimSuffix(replacer.Replace(p), templates.TmplSuffix)
		outPath := filepath.Join(destDir, filepath.FromSlash(outRel))

		if d.IsDir() {
			return os.MkdirAll(outPath, 0o755)
		}

		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}

		mode := os.FileMode(0o644)
		if !isBinary(data) {
			data = []byte(replacer.Replace(string(data)))
			if bytes.HasPrefix(data, []byte("#!")) {
				mode = 0o755
			}
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
			return err
		}
		return os.WriteFile(outPath, data, mode)
	})
}

// replacerFor builds a string replacer that maps each field's token to the
// supplied value (empty when a value is missing), falling back to the field's
// Default for blank entries.
func replacerFor(tmpl templates.Template, values map[string]string) *strings.Replacer {
	pairs := make([]string, 0, len(tmpl.Fields)*2)
	for _, f := range tmpl.Fields {
		v := values[f.Key]
		if v == "" {
			v = f.Default
		}
		pairs = append(pairs, f.Token, v)
	}
	return strings.NewReplacer(pairs...)
}

// ensureEmptyDir creates destDir if missing and errors if it exists with
// content (or exists as a non-directory).
func ensureEmptyDir(dir string) error {
	info, err := os.Stat(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return os.MkdirAll(dir, 0o755)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("destination %q exists and is not a directory", dir)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("%w: %s", ErrDestinationNotEmpty, dir)
	}
	return nil
}

// isBinary reports whether data looks non-textual (contains a NUL byte), in
// which case token replacement is skipped and the bytes are copied verbatim.
func isBinary(data []byte) bool {
	return bytes.IndexByte(data, 0) >= 0
}
