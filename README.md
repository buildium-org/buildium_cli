# Buildium CLI

A command-line interface for interacting with the Buildium platform. Create project templates and manage your Buildium workflow directly from your terminal.

## Prerequisites

- **Go** 1.25.5 or later
- **Git** (required for cloning project templates)

## Building

Build the CLI from source:

```bash
make build
```

This compiles the binary to `./buildium` in the current directory.

## Installation

After building, add the CLI to your PATH for global access:

```bash
export PATH=$PATH:/path/to/buildium_cli
```

To make this permanent, add the above line to your shell configuration file (`~/.zshrc`, `~/.bashrc`, etc.).

## Commands

### `login`

Authenticate with the Buildium platform. Your credentials are stored locally for subsequent commands.

```bash
buildium login -email <your-email> -password <your-password>
```

**Flags:**
| Flag | Required | Description |
|------|----------|-------------|
| `-email` | Yes | Your Buildium account email |
| `-password` | Yes | Your Buildium account password |

**Example:**
```bash
buildium login -email user@example.com -password mysecretpassword
```

---

### `project create-template`

Scaffold a new project from a starter template. This clones the appropriate language template and configures it with your project details.

```bash
buildium project create-template -projectid <id> -lang <language> -name <repo-name>
```

**Flags:**
| Flag | Required | Description |
|------|----------|-------------|
| `-projectid` | Yes | Your Buildium project ID |
| `-lang` | Yes | Programming language (`go` or `typescript`) |
| `-name` | Yes | Name for the cloned repository directory |

**Supported Languages:**
- `go` - Clones from [buildium-org/go_template](https://github.com/buildium-org/go_template)
- `typescript` - Clones from [buildium-org/ts_template](https://github.com/buildium-org/ts_template)

**Example:**
```bash
buildium project create-template -projectid abc123 -lang go -name my-redis-clone
```

This will:
1. Clone the Go template repository into `./my-redis-clone`
2. Configure the project ID in all relevant files
3. Set up the Docker image name based on your project's tutorial

## Configuration

The CLI stores its configuration in `.buildium/config.json` adjacent to the executable. This file contains:

- `Environment` - The target environment (`PROD` by default)
- `AuthToken` - Your authentication token (set after login)

**Note:** You must run `buildium login` before using commands that require authentication (like `project create-template`).

## Usage Flow

1. **Login** to authenticate with Buildium:
   ```bash
   buildium login -email you@example.com -password yourpassword
   ```

2. **Create a project** on the Buildium web platform and note your project ID

3. **Generate your starter template**:
   ```bash
   buildium project create-template -projectid <your-project-id> -lang go -name my-project
   ```

4. **Start coding** in your new project directory!

## Troubleshooting

**"Not logged in" error**  
Run `buildium login` with your credentials before using other commands.

**"Failed to get project" error**  
Verify your project ID is correct and that you have access to the project.

**Clone fails**  
Ensure you have git installed and network access to GitHub.

