// Package archive orchestrates the archiving process.
package archive

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/deanhigh/bury-it/internal/git"
	"github.com/deanhigh/bury-it/internal/graveyard"
	"github.com/deanhigh/bury-it/internal/metadata"
	"github.com/deanhigh/bury-it/internal/source"
)

// Options contains the options for the archive operation.
type Options struct {
	// Source is the source repository string (URL, owner/repo, or path).
	Source string
	// Graveyard is the path to the graveyard repository.
	Graveyard string
	// Name is an optional override for the project name in the graveyard.
	Name string
	// DropHistory indicates whether to drop git history.
	DropHistory bool
}

// Result contains the result of the archive operation.
type Result struct {
	// ProjectName is the name of the archived project.
	ProjectName string
	// ProjectPath is the path to the archived project in the graveyard.
	ProjectPath string
	// HistoryPreserved indicates whether git history was preserved.
	HistoryPreserved bool
}

// Archive archives a source repository into a graveyard.
func Archive(opts Options) (*Result, error) {
	// Parse source
	src, err := source.Parse(opts.Source)
	if err != nil {
		return nil, fmt.Errorf("invalid source: %w", err)
	}

	// Parse graveyard
	gy, err := graveyard.New(opts.Graveyard)
	if err != nil {
		return nil, fmt.Errorf("invalid graveyard: %w", err)
	}

	// Validate graveyard
	if err := gy.Validate(); err != nil {
		return nil, err
	}

	// Determine project name
	projectName := src.Name
	if opts.Name != "" {
		projectName = opts.Name
	}

	// Validate project name
	if err := gy.ValidateProjectName(projectName); err != nil {
		return nil, err
	}

	// Handle remote repositories
	var localSourcePath string
	var tempDir string
	if src.Type == source.TypeRemote {
		// Clone to temp directory
		tempDir, err = os.MkdirTemp("", "bury-it-*")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
		defer func() { _ = os.RemoveAll(tempDir) }()

		clonePath := filepath.Join(tempDir, projectName)
		fmt.Printf("Cloning %s...\n", src.Path)
		if err := git.Clone(src.Path, clonePath); err != nil {
			return nil, fmt.Errorf("failed to clone repository: %w", err)
		}
		localSourcePath = clonePath
	} else {
		// Validate local source
		if err := src.Validate(); err != nil {
			return nil, err
		}
		localSourcePath = src.Path
	}

	// Get display path for metadata before any operations
	displayPath := src.DisplayPath()

	// Archive the project
	projectPath := gy.ProjectPath(projectName)
	historyPreserved := !opts.DropHistory

	if opts.DropHistory {
		// Copy only tracked files (respects .gitignore)
		fmt.Printf("Copying tracked files (without history) to %s...\n", projectName)
		if err := git.CopyTrackedFiles(localSourcePath, projectPath); err != nil {
			return nil, fmt.Errorf("failed to copy files: %w", err)
		}
	} else {
		// Use subtree to preserve history
		fmt.Printf("Adding %s with full history...\n", projectName)
		if err := git.SubtreeAdd(gy.Path, localSourcePath, projectName); err != nil {
			return nil, fmt.Errorf("failed to add subtree: %w", err)
		}
	}

	// Generate and write metadata
	meta := &metadata.Metadata{
		OriginalSource:   displayPath,
		BuriedAt:         time.Now(),
		HistoryPreserved: historyPreserved,
	}
	if err := meta.Write(projectPath); err != nil {
		return nil, err
	}

	// Stage the metadata file (and all files if drop-history was used)
	if opts.DropHistory {
		if err := git.StageAll(gy.Path); err != nil {
			return nil, fmt.Errorf("failed to stage files: %w", err)
		}
	} else {
		// For subtree, only stage the metadata file
		metaPath := filepath.Join(projectName, metadata.FileName)
		if err := git.StageFile(gy.Path, metaPath); err != nil {
			return nil, fmt.Errorf("failed to stage metadata: %w", err)
		}
	}

	// Auto-commit the archived project
	commitMsg := fmt.Sprintf("docs: bury-it - archived %s", projectName)
	fmt.Printf("Committing to graveyard...\n")
	if err := git.Commit(gy.Path, commitMsg); err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	return &Result{
		ProjectName:      projectName,
		ProjectPath:      projectPath,
		HistoryPreserved: historyPreserved,
	}, nil
}
