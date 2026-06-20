package templates

import "io/fs"

// Token convention
//
// A template's replaceable values are marked in its files with literal
// upper-snake markers wrapped in angle brackets and suffixed with `_HERE`,
// e.g. <YOUR_IMAGE_NAME_HERE>. These are exactly the placeholders already
// shipped in the upstream Buildium template repos, so the vendored files stay
// byte-for-byte identical to their source and the markers are easy to spot by
// eye. The generator replaces every occurrence of a field's Token with the
// user-supplied value.
const (
	tokenImageName      = "<YOUR_IMAGE_NAME_HERE>"
	tokenHarnessImage   = "<TEST_HARNESS_IMAGE_HERE>"
	tokenDockerUsername = "<YOUR_DOCKER_USERNAME_HERE>"
)

// Field is one replaceable value in a template. The generator replaces every
// occurrence of Token (in file contents and in path names) with the value the
// user provides for this field.
type Field struct {
	// Key is a stable identifier for the field (e.g. "image_name"). It is the
	// map key the TUI collects values under and the engine reads them from.
	Key string
	// Token is the literal marker that appears in the template's files
	// (e.g. "<YOUR_IMAGE_NAME_HERE>").
	Token string
	// Label is the human-facing prompt shown in the wizard.
	Label string
	// Help is an optional one-line hint shown beneath the prompt.
	Help string
	// Required rejects an empty submission when true.
	Required bool
	// Default is the value used when the user leaves the field blank.
	Default string
}

// Template is one generatable project: a stable key (which is also its embedded
// FS root), a human label, and the ordered fields the user fills in.
type Template struct {
	Key    string
	Label  string
	Fields []Field
}

// FS returns the embedded file tree backing this template, rooted at its files
// (so paths look like "Dockerfile", "manifest/info.json").
func (t Template) FS() (fs.FS, error) { return Sub(t.Key) }

// Shared field definitions. Solutions (go, ts) take the same three fields; the
// tutorial only needs an image name (used to derive its harness image tag).
var (
	fieldTutorialName = Field{
		Key:      "image_name",
		Token:    tokenImageName,
		Label:    "Tutorial name",
		Help:     "Base name for the tutorial; the harness image is tagged <name>_harness",
		Required: true,
	}
	fieldDockerUsername = Field{
		Key:      "docker_username",
		Token:    tokenDockerUsername,
		Label:    "Docker username",
		Help:     "Your Docker username",
		Required: true,
	}
	fieldImageName = Field{
		Key:      "image_name",
		Token:    tokenImageName,
		Label:    "Docker image name",
		Help:     "Tag for the image you build and run (lowercase, no spaces)",
		Required: true,
	}
	fieldHarnessImage = Field{
		Key:      "harness_image",
		Token:    tokenHarnessImage,
		Label:    "Test harness base image",
		Help:     "The harness image published by the tutorial author (e.g. my-tutorial_harness)",
		Required: true,
	}
)

// Catalog returns every template the CLI can generate, in display order. It is
// the single source of truth shared by the wizard (which prompts for each
// template's fields) and the generator (which replaces each field's token).
func Catalog() []Template {
	return []Template{
		{
			Key:    "tutorial",
			Label:  "Tutorial template",
			Fields: []Field{fieldTutorialName, fieldDockerUsername},
		},
		{
			Key:    "go",
			Label:  "Solution template (Go)",
			Fields: []Field{fieldImageName, fieldHarnessImage},
		},
		{
			Key:    "ts",
			Label:  "Solution template (TypeScript)",
			Fields: []Field{fieldImageName, fieldHarnessImage},
		},
	}
}

// ByKey returns the template with the given key, and whether it was found.
func ByKey(key string) (Template, bool) {
	for _, t := range Catalog() {
		if t.Key == key {
			return t, true
		}
	}
	return Template{}, false
}
