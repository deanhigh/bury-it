// Package graveyard handles graveyard repository operations.
package graveyard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/deanhigh/bury-it/internal/git"
)

// Graveyard represents a graveyard repository.
type Graveyard struct {
	// Path is the absolute path to the graveyard repository.
	Path string
}

// New creates a new Graveyard instance from the given path.
func New(path string) (*Graveyard, error) {
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
		return nil, fmt.Errorf("failed to resolve graveyard path: %w", err)
	}

	return &Graveyard{Path: absPath}, nil
}

// Validate checks that the graveyard is a valid git repository.
func (g *Graveyard) Validate() error {
	// Check if path exists
	info, err := os.Stat(g.Path)
	if os.IsNotExist(err) {
		return fmt.Errorf("graveyard path does not exist: %s", g.Path)
	}
	if err != nil {
		return fmt.Errorf("failed to access graveyard path: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("graveyard path is not a directory: %s", g.Path)
	}

	// Check if it's a git repository
	if !git.IsValidRepo(g.Path) {
		return fmt.Errorf("graveyard is not a git repository: %s", g.Path)
	}

	return nil
}

// ProjectPath returns the path where a project would be archived.
func (g *Graveyard) ProjectPath(name string) string {
	return filepath.Join(g.Path, name)
}

// ProjectExists checks if a project already exists in the graveyard.
func (g *Graveyard) ProjectExists(name string) bool {
	projectPath := g.ProjectPath(name)
	info, err := os.Stat(projectPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ValidateProjectName checks if a project name can be used.
func (g *Graveyard) ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check for invalid characters
	if strings.ContainsAny(name, "/\\:*?\"<>|") {
		return fmt.Errorf("project name contains invalid characters: %s", name)
	}

	// Check for reserved names
	if name == "." || name == ".." {
		return fmt.Errorf("project name cannot be '.' or '..'")
	}

	// Check if project already exists
	if g.ProjectExists(name) {
		return fmt.Errorf("project already exists in graveyard: %s (use --name to specify an alternative name)", name)
	}

	return nil
}
