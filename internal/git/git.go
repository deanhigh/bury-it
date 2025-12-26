// Package git provides a wrapper around git commands.
package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsValidRepo checks if the given path is a valid git repository.
func IsValidRepo(path string) bool {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Clone clones a remote repository to the destination path.
func Clone(url, dest string) error {
	cmd := exec.Command("git", "clone", url, dest)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
}

// GetRemoteURL returns the origin remote URL for a repository.
func GetRemoteURL(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "remote", "get-url", "origin")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		// No remote is not necessarily an error for local repos
		return "", nil
	}
	return strings.TrimSpace(stdout.String()), nil
}

// GetDefaultBranch returns the default branch name for a repository.
func GetDefaultBranch(repoPath string) (string, error) {
	// Try to get the current branch first
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to get branch: %s", strings.TrimSpace(stderr.String()))
	}
	branch := strings.TrimSpace(stdout.String())
	if branch == "" || branch == "HEAD" {
		// Detached HEAD, try common branch names
		for _, name := range []string{"main", "master"} {
			cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--verify", name)
			if cmd.Run() == nil {
				return name, nil
			}
		}
		return "", fmt.Errorf("unable to determine default branch")
	}
	return branch, nil
}

// SubtreeAdd adds a repository as a subtree with full history.
func SubtreeAdd(graveyardPath, sourceRepoPath, prefix string) error {
	// Get the default branch of the source repo
	branch, err := GetDefaultBranch(sourceRepoPath)
	if err != nil {
		return fmt.Errorf("failed to get source branch: %w", err)
	}

	// Get absolute path to source repo
	absSourcePath, err := filepath.Abs(sourceRepoPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Add as subtree
	cmd := exec.Command("git", "-C", graveyardPath, "subtree", "add",
		"--prefix="+prefix, absSourcePath, branch)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git subtree add failed: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
}

// CopyTrackedFiles copies only git-tracked files from source to destination.
// This respects .gitignore by using git archive to export only tracked files.
func CopyTrackedFiles(sourcePath, destPath string) error {
	// Create destination directory
	if err := os.MkdirAll(destPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Use git archive to create a tar of tracked files, then extract
	// This automatically respects .gitignore since only tracked files are included
	archiveCmd := exec.Command("git", "-C", sourcePath, "archive", "--format=tar", "HEAD")
	extractCmd := exec.Command("tar", "-xf", "-", "-C", destPath)

	// Pipe archive output to tar extract
	var archiveStderr, extractStderr bytes.Buffer
	archiveCmd.Stderr = &archiveStderr
	extractCmd.Stderr = &extractStderr

	pipe, err := archiveCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}
	extractCmd.Stdin = pipe

	if err := archiveCmd.Start(); err != nil {
		return fmt.Errorf("git archive failed to start: %w", err)
	}
	if err := extractCmd.Start(); err != nil {
		return fmt.Errorf("tar extract failed to start: %w", err)
	}

	if err := archiveCmd.Wait(); err != nil {
		return fmt.Errorf("git archive failed: %s", strings.TrimSpace(archiveStderr.String()))
	}
	if err := extractCmd.Wait(); err != nil {
		return fmt.Errorf("tar extract failed: %s", strings.TrimSpace(extractStderr.String()))
	}

	return nil
}

// StageAll stages all changes in the repository.
func StageAll(repoPath string) error {
	cmd := exec.Command("git", "-C", repoPath, "add", "-A")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
}

// StageFile stages a specific file in the repository.
func StageFile(repoPath, filePath string) error {
	cmd := exec.Command("git", "-C", repoPath, "add", filePath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
}

// Commit creates a commit with the given message.
func Commit(repoPath, message string) error {
	cmd := exec.Command("git", "-C", repoPath, "commit", "-m", message)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %s", strings.TrimSpace(stderr.String()))
	}
	return nil
}
