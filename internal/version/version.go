// Package version exposes the build-time version of the buildium CLI.
//
// The values default to development placeholders and are overridden at build
// time via the linker, e.g.:
//
//	go build -ldflags "-X buildium_cli/internal/version.Version=1.2.3" main.go
//
// The release tooling (GoReleaser) populates these from the git tag.
package version

var (
	// Version is the semantic version of the binary (e.g. "1.2.3"), or "dev"
	// for an unstamped local build.
	Version = "dev"
	// Commit is the git commit the binary was built from, when stamped.
	Commit = "none"
	// Date is the build timestamp, when stamped.
	Date = "unknown"
)

// String returns a human-readable version line for the CLI.
func String() string {
	return "buildium " + Version + " (commit " + Commit + ", built " + Date + ")"
}
