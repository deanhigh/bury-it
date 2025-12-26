# bury-it

A CLI tool to sunset experimental projects by archiving them into a local "graveyard" repository while preserving their git history.

## Features

- Archive GitHub repositories or local git repos to a graveyard
- Preserve full git history (or optionally drop it)
- Creates metadata file for each buried project
- Works on macOS, Linux, and Windows

## Installation

### From Source

```bash
go install github.com/deanhigh/bury-it@latest
```

### Build Locally

```bash
git clone https://github.com/deanhigh/bury-it.git
cd bury-it
make build
```

## Usage

```bash
# Show help
bury-it --help

# Bury a GitHub repository
bury-it --source deanhigh/old-project --graveyard ~/graveyard

# Bury a local repository
bury-it --source ./my-experiment --graveyard ~/graveyard

# Bury without preserving history
bury-it --source ./my-experiment --graveyard ~/graveyard --drop-history
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--source` | `-s` | Source repository (GitHub URL, owner/repo, or local path) |
| `--graveyard` | `-g` | Local path to the graveyard repository |
| `--drop-history` | | Archive only the latest state, discard git history |
| `--help` | `-h` | Show help message |
| `--version` | `-v` | Show version |

## How It Works

1. Validates the source repository exists and is a valid git repo
2. Checks the graveyard location (creates if needed)
3. Archives the project as a subdirectory in the graveyard
4. Creates a `.bury-it.md` metadata file with archive details
5. Reminds you to commit the graveyard and archive the original

**Note**: bury-it does not delete the original repository. After burying, you should manually commit the graveyard changes and archive/delete the original.

## License

MIT License - see [LICENSE](LICENSE) for details.
