# Buildium CLI

An interactive terminal app for scaffolding new Buildium **tutorials** and
**solution templates**. The starter templates are packaged inside the binary, so
generation works completely offline — no `git clone`, no network. Pick a
template in the wizard, fill in a few fields, and the CLI writes a ready-to-edit
project to disk with those values substituted in.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Prerequisites

- **Go** 1.25 or later (only to build the CLI)

That's it — because the templates are embedded, you do **not** need `git` or
network access to generate a project.

## Building

```bash
make build
```

This compiles the binary to `./buildium` in the current directory.

## Installation

After building, add the CLI to your `PATH` for global access:

```bash
export PATH=$PATH:/path/to/buildium_cli
```

To make this permanent, add the line to your shell config (`~/.zshrc`,
`~/.bashrc`, etc.).

## Usage

Launch the wizard:

```bash
./buildium
```

The wizard walks you through four steps:

1. **Select a template** — choose Tutorial, Solution (Go), or Solution
   (TypeScript) from the list (`↑`/`↓` to move, `enter` to choose).
2. **Fill in the fields** — enter a destination directory plus the fields the
   chosen template needs (`tab`/`↑`/`↓` to move between fields, `enter` to
   advance / submit). Required fields can't be left blank.
3. **Review** — confirm the template, destination, and values (`enter`/`y` to
   generate, `esc`/`n` to go back).
4. **Done** — the project is written to your destination directory with every
   field substituted in. The CLI refuses to write into a directory that already
   exists and isn't empty.

Press `ctrl+c` to quit at any point.

## Templates and their fields

Each template ships with placeholder tokens (`<..._HERE>`) in its files; the
wizard replaces them with the values you enter.

| Template | Field | Replaces | Notes |
|----------|-------|----------|-------|
| **Tutorial** | Tutorial name | `<YOUR_IMAGE_NAME_HERE>` | The test-harness image is tagged `<name>_harness`. |
| **Solution (Go)** | Buildium project ID | `<PROJECT_ID_HERE>` | From the project page on the Buildium website. |
| | Docker image name | `<YOUR_IMAGE_NAME_HERE>` | Tag for the image you build and run. |
| | Test harness base image | `<TEST_HARNESS_IMAGE_HERE>` | The harness image published by the tutorial author. |
| **Solution (TypeScript)** | Buildium project ID | `<PROJECT_ID_HERE>` | From the project page on the Buildium website. |
| | Docker image name | `<YOUR_IMAGE_NAME_HERE>` | Tag for the image you build and run. |
| | Test harness base image | `<TEST_HARNESS_IMAGE_HERE>` | The harness image published by the tutorial author. |

The source templates live in these repositories and are vendored into the CLI:

- Tutorial — [buildium-org/tutorial_template](https://github.com/buildium-org/tutorial_template)
- Go solution — [buildium-org/go_template](https://github.com/buildium-org/go_template)
- TypeScript solution — [buildium-org/ts_template](https://github.com/buildium-org/ts_template)

## Usage flow

### For tutorial authors

1. Run `./buildium`, choose **Tutorial template**, and enter a tutorial name.
2. Customize the generated tutorial — edit the manifest, stages, and step
   harness in your new directory.

### For tutorial participants

1. Create a project on the Buildium website and note your project ID.
2. Run `./buildium`, choose **Solution (Go)** or **Solution (TypeScript)**, and
   enter your project ID, image name, and the harness image from the tutorial.
3. Start coding in your new project directory.

## Maintaining the templates

The vendored template trees live under
[`internal/templates/files/`](internal/templates/files/) (one directory per
template key: `tutorial`, `go`, `ts`). To update a template, edit the files
there and rebuild.

Two conventions apply to the vendored files:

- **Placeholder tokens** use the upstream `<UPPER_SNAKE_HERE>` markers. Add a
  matching `Field` in [`internal/templates/catalog.go`](internal/templates/catalog.go)
  for any new token so the wizard prompts for it.
- **Go source files** (`go.mod`, `*.go`) are stored with a trailing `.tmpl`
  suffix (e.g. `main.go.tmpl`). This keeps a nested `go.mod` from being excluded
  by `go:embed` and keeps stray `.go` files out of this module's build. The
  generator strips the `.tmpl` suffix when it writes the project.

## Troubleshooting

**"destination already exists and is not empty"**
Choose a destination directory that doesn't exist yet (or is empty).

**The UI doesn't render / exits immediately**
The wizard needs an interactive terminal (a TTY). Run it directly in your
terminal rather than through a pipe or non-interactive shell.
