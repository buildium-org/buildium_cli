package templates

import (
	"embed"
	"io/fs"
)

// FS holds every vendored template tree, packaged into the binary at build
// time. Paths inside it are rooted at "files/", e.g. "files/go/Dockerfile".
//
// The `all:` prefix is required so dotfiles such as each template's
// `.gitignore` are embedded too (a bare `//go:embed files` would skip any
// path element beginning with `.` or `_`).
//
// Vendoring convention: every Go source file in a template — both `go.mod`
// and any `.go` file — is stored with a trailing `.tmpl` suffix
// (`go.mod.tmpl`, `main.go.tmpl`, ...). Two distinct toolchain problems force
// this:
//
//   - A directory containing a `go.mod` is treated as a nested module, which
//     the embed machinery silently excludes from the embedded tree.
//   - A `.go` file anywhere under the parent module is compiled and vetted by
//     `go build ./...` / `go vet ./...`; the template sources reference the
//     external Buildium harness and would not compile here.
//
// Renaming them to `.tmpl` hides them from the toolchain while keeping the
// bytes embeddable. The generator strips the trailing `.tmpl` from any
// embedded path when materializing a project, restoring the real filename.
//
//go:embed all:files
var FS embed.FS

// TmplSuffix marks template files whose name must be transformed on output:
// the generator removes this suffix when writing the file to disk (so
// `main.go.tmpl` is written as `main.go`). See the FS docs for why template
// Go sources carry it.
const TmplSuffix = ".tmpl"

// Sub returns an fs.FS rooted at a single template's directory, identified by
// its key ("tutorial", "go", "ts"). The returned filesystem's paths are
// relative to that template root (e.g. "Dockerfile", "manifest/info.json").
func Sub(key string) (fs.FS, error) {
	return fs.Sub(FS, "files/"+key)
}
