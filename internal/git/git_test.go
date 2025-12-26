package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIsValidRepo(t *testing.T) {
	// Create a temp directory
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	// Create a valid repo structure
	validRepo := filepath.Join(tempDir, "valid-repo")
	if err := os.MkdirAll(filepath.Join(validRepo, ".git"), 0755); err != nil {
		t.Fatalf("Failed to create valid repo: %v", err)
	}

	// Create an invalid directory (no .git)
	invalidDir := filepath.Join(tempDir, "invalid-dir")
	if err := os.MkdirAll(invalidDir, 0755); err != nil {
		t.Fatalf("Failed to create invalid dir: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "valid repo with .git directory",
			path: validRepo,
			want: true,
		},
		{
			name: "directory without .git",
			path: invalidDir,
			want: false,
		},
		{
			name: "non-existent path",
			path: filepath.Join(tempDir, "does-not-exist"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidRepo(tt.path)
			if got != tt.want {
				t.Errorf("IsValidRepo(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestCopyTrackedFiles(t *testing.T) {
	// Create a real git repo to test with
	sourceDir, err := os.MkdirTemp("", "git-copy-source-*")
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sourceDir) })

	// Initialize git repo
	if err := runGit(sourceDir, "init"); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	if err := runGit(sourceDir, "config", "user.email", "test@test.com"); err != nil {
		t.Fatalf("Failed to set git email: %v", err)
	}
	if err := runGit(sourceDir, "config", "user.name", "Test"); err != nil {
		t.Fatalf("Failed to set git name: %v", err)
	}

	// Create tracked files
	if err := os.WriteFile(filepath.Join(sourceDir, "tracked.txt"), []byte("tracked content"), 0644); err != nil {
		t.Fatalf("Failed to create tracked file: %v", err)
	}

	subDir := filepath.Join(sourceDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested content"), 0644); err != nil {
		t.Fatalf("Failed to create nested file: %v", err)
	}

	// Create .gitignore and ignored file
	if err := os.WriteFile(filepath.Join(sourceDir, ".gitignore"), []byte("ignored.txt\n"), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "ignored.txt"), []byte("should be ignored"), 0644); err != nil {
		t.Fatalf("Failed to create ignored file: %v", err)
	}

	// Add and commit tracked files
	if err := runGit(sourceDir, "add", "tracked.txt", "subdir/nested.txt", ".gitignore"); err != nil {
		t.Fatalf("Failed to add files: %v", err)
	}
	if err := runGit(sourceDir, "commit", "-m", "initial commit"); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Create destination directory
	destDir, err := os.MkdirTemp("", "git-copy-dest-*")
	if err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(destDir) })

	// Copy tracked files
	if err := CopyTrackedFiles(sourceDir, destDir); err != nil {
		t.Fatalf("CopyTrackedFiles() error = %v", err)
	}

	// Verify tracked files were copied, ignored files were not
	tests := []struct {
		path        string
		shouldExist bool
	}{
		{filepath.Join(destDir, "tracked.txt"), true},
		{filepath.Join(destDir, "subdir", "nested.txt"), true},
		{filepath.Join(destDir, ".gitignore"), true},
		{filepath.Join(destDir, "ignored.txt"), false}, // should be excluded
		{filepath.Join(destDir, ".git"), false},        // .git should never be copied
	}

	for _, tt := range tests {
		_, err := os.Stat(tt.path)
		exists := err == nil
		if exists != tt.shouldExist {
			t.Errorf("Path %q exists = %v, want %v", tt.path, exists, tt.shouldExist)
		}
	}

	// Verify content
	content, err := os.ReadFile(filepath.Join(destDir, "tracked.txt"))
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}
	if string(content) != "tracked content" {
		t.Errorf("File content = %q, want %q", string(content), "tracked content")
	}
}

// runGit is a helper to run git commands in tests.
func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	return cmd.Run()
}
