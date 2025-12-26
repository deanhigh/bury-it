# bury-it Requirements

## Overview

A CLI tool to sunset experimental projects by archiving them into a local "graveyard" repository while optionally preserving their full git history.

## Functional Requirements

### FR-1: Source Repository Support

- **FR-1.1**: Accept remote GitHub repositories (via URL or owner/repo format)
- **FR-1.2**: Accept local git repositories (via filesystem path)
- **FR-1.3**: Validate that the source is a valid git repository

### FR-2: Graveyard Repository

- **FR-2.1**: Archive projects into a specified local graveyard repository
- **FR-2.2**: Each archived project stored as a subdirectory (flat structure)
- **FR-2.3**: Fail with clear error if project name already exists in graveyard
- **FR-2.4**: Allow an alternative directory name in graveyard to address 2.3

### FR-3: History Management

- **FR-3.1**: By default, preserve full git history when archiving
- **FR-3.2**: Support `--drop-history` flag to archive only latest state
- **FR-3.3**: Respect `.gitignore` - do not archive ignored files

### FR-4: Metadata

- **FR-4.1**: Create `.bury-it.md` in each archived project with:
  - Original source location
  - Date/time buried
  - Whether history was preserved

### FR-5: CLI Interface

- **FR-5.1**: Display help when invoked with no arguments
- **FR-5.2**: Display help with `--help` or `-h` flag
- **FR-5.3**: Provide clear error messages for invalid inputs
- **FR-5.4**: Automatically commit the archived project with message: `docs: bury-it - archived <project-name>`

## Non-Functional Requirements

### NFR-1: Portability

- Run on macOS, Linux, and Windows

### NFR-2: No External Dependencies

- Single static binary, no runtime dependencies (git must be installed)


