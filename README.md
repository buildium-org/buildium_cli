# Buildium CLI

An interactive terminal app for scaffolding new Buildium **tutorials** and
**solution templates**. The starter templates are packaged inside the binary, so
generation works completely offline — no `git clone`, no network. Pick a
template in the wizard, fill in a few fields, and the CLI writes a ready-to-edit
project to disk with those values substituted in.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Download

Pre-built binaries for every supported platform are attached to each
[GitHub release](https://github.com/buildium-org/buildium_cli/releases). The
links below always resolve to the **latest** release, so they're safe to bookmark
or host elsewhere:

| Platform | Architecture | Download |
|----------|--------------|----------|
| macOS (Apple Silicon) | arm64 | [buildium_darwin_arm64.tar.gz](https://github.com/buildium-org/buildium_cli/releases/latest/download/buildium_darwin_arm64.tar.gz) |
| macOS (Intel) | amd64 | [buildium_darwin_amd64.tar.gz](https://github.com/buildium-org/buildium_cli/releases/latest/download/buildium_darwin_amd64.tar.gz) |
| Linux | amd64 | [buildium_linux_amd64.tar.gz](https://github.com/buildium-org/buildium_cli/releases/latest/download/buildium_linux_amd64.tar.gz) |
| Linux | arm64 | [buildium_linux_arm64.tar.gz](https://github.com/buildium-org/buildium_cli/releases/latest/download/buildium_linux_arm64.tar.gz) |
| Windows | amd64 | [buildium_windows_amd64.zip](https://github.com/buildium-org/buildium_cli/releases/latest/download/buildium_windows_amd64.zip) |
| Windows | arm64 | [buildium_windows_arm64.zip](https://github.com/buildium-org/buildium_cli/releases/latest/download/buildium_windows_arm64.zip) |

Checksums for every asset are published alongside them as
[checksums.txt](https://github.com/buildium-org/buildium_cli/releases/latest/download/checksums.txt).

After downloading, extract the archive and run the `buildium` binary. On
**macOS** and **Linux** make it executable first, and on macOS clear the
quarantine flag so Gatekeeper doesn't block the unsigned binary:

```bash
tar -xzf buildium_darwin_arm64.tar.gz
chmod +x buildium
xattr -d com.apple.quarantine buildium   # macOS only
./buildium
```

Confirm the version with `./buildium --version`.

> Want the binary from an in-flight change? Every pull request builds the same
> set of binaries and uploads them as **artifacts** under that PR's run in the
> [Actions tab](https://github.com/buildium-org/buildium_cli/actions) (look for
> `buildium-binaries`). These are unversioned snapshot builds for testing.

## Prerequisites

- **Go** 1.25 or later (only to build the CLI from source)

That's it — because the templates are embedded, you do **not** need `git` or
network access to generate a project.

> Most users should grab a [pre-built binary](#download) instead — building from
> source is only needed if you're modifying the CLI or its templates.

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

## Releasing

Releases are cut by pushing a version tag. The
[`release` workflow](.github/workflows/release.yml) then cross-compiles every
platform with [GoReleaser](https://goreleaser.com) and publishes a GitHub Release
with all the archives plus `checksums.txt` attached:

```bash
git tag v0.1.0
git push origin v0.1.0
```

The version you tag is stamped into the binary (`buildium --version`). Because
the release assets keep stable names, the
[download links above](#download) automatically point at the new release once it
finishes publishing — nothing on your website needs to change.

## Troubleshooting

**"destination already exists and is not empty"**
Choose a destination directory that doesn't exist yet (or is empty).

**The UI doesn't render / exits immediately**
The wizard needs an interactive terminal (a TTY). Run it directly in your
terminal rather than through a pipe or non-interactive shell.
