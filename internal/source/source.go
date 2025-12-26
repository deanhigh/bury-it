// Package source handles parsing and validation of source repositories.
package source

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/deanhigh/bury-it/internal/git"
)

// Type represents the type of source repository.
type Type int

const (
	// TypeLocal represents a local filesystem repository.
	TypeLocal Type = iota
	// TypeRemote represents a remote GitHub repository.
	TypeRemote
)

// Source represents a parsed source repository.
type Source struct {
	// Type is the source type (local or remote).
	Type Type
	// Path is the local filesystem path (for local repos) or the URL (for remote repos).
	Path string
	// Name is the extracted project name.
	Name string
	// OriginalInput is the original input string.
	OriginalInput string
}

// gitHubURLPattern matches GitHub URLs.
var gitHubURLPattern = regexp.MustCompile(`^https?://github\.com/([^/]+)/([^/]+?)(?:\.git)?/?$`)

// ownerRepoPattern matches owner/repo shorthand.
var ownerRepoPattern = regexp.MustCompile(`^([a-zA-Z0-9_.-]+)/([a-zA-Z0-9_.-]+)$`)

// Parse parses the input string and returns a Source.
func Parse(input string) (*Source, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("source cannot be empty")
	}

	// Check if it's a GitHub URL
	if matches := gitHubURLPattern.FindStringSubmatch(input); matches != nil {
		return &Source{
			Type:          TypeRemote,
			Path:          input,
			Name:          matches[2],
			OriginalInput: input,
		}, nil
	}

	// Check if it's owner/repo shorthand (but not a local path like ./foo or /foo)
	if !strings.HasPrefix(input, ".") && !strings.HasPrefix(input, "/") && !strings.HasPrefix(input, "~") {
		if matches := ownerRepoPattern.FindStringSubmatch(input); matches != nil {
			url := fmt.Sprintf("https://github.com/%s/%s", matches[1], matches[2])
			return &Source{
				Type:          TypeRemote,
				Path:          url,
				Name:          matches[2],
				OriginalInput: input,
			}, nil
		}
	}

	// Treat as local path
	path := input

	// Expand ~ to home directory
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to expand home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Extract project name from path
	name := filepath.Base(absPath)

	return &Source{
		Type:          TypeLocal,
		Path:          absPath,
		Name:          name,
		OriginalInput: input,
	}, nil
}

// Validate validates that the source is a valid git repository.
func (s *Source) Validate() error {
	switch s.Type {
	case TypeLocal:
		// Check if path exists
		info, err := os.Stat(s.Path)
		if os.IsNotExist(err) {
			return fmt.Errorf("source path does not exist: %s", s.Path)
		}
		if err != nil {
			return fmt.Errorf("failed to access source path: %w", err)
		}
		if !info.IsDir() {
			return fmt.Errorf("source path is not a directory: %s", s.Path)
		}
		// Check if it's a git repository
		if !git.IsValidRepo(s.Path) {
			return fmt.Errorf("source is not a git repository: %s", s.Path)
		}
	case TypeRemote:
		// Remote repos will be validated during clone
		// We could add a lightweight check here (e.g., git ls-remote) but that
		// would add latency for valid repos. We'll let clone fail with a clear error.
	}
	return nil
}

// DisplayPath returns a human-readable path for display purposes.
func (s *Source) DisplayPath() string {
	if s.Type == TypeRemote {
		return s.Path
	}
	// For local repos, try to get remote URL, otherwise use path
	if remoteURL, err := git.GetRemoteURL(s.Path); err == nil && remoteURL != "" {
		return remoteURL
	}
	return s.Path
}
