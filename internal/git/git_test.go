package git

import (
	"os"
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

func TestCopyFiles(t *testing.T) {
	// Create source directory with files
	sourceDir, err := os.MkdirTemp("", "git-copy-source-*")
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(sourceDir) })

	// Create some files
	if err := os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644); err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}

	subDir := filepath.Join(sourceDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "file2.txt"), []byte("content2"), 0644); err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}

	// Create .git directory (should be excluded)
	gitDir := filepath.Join(sourceDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(gitDir, "config"), []byte("git config"), 0644); err != nil {
		t.Fatalf("Failed to create git config: %v", err)
	}

	// Create destination directory
	destDir, err := os.MkdirTemp("", "git-copy-dest-*")
	if err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(destDir) })

	// Copy files
	if err := CopyFiles(sourceDir, destDir); err != nil {
		t.Fatalf("CopyFiles() error = %v", err)
	}

	// Verify files were copied
	tests := []struct {
		path        string
		shouldExist bool
	}{
		{filepath.Join(destDir, "file1.txt"), true},
		{filepath.Join(destDir, "subdir", "file2.txt"), true},
		{filepath.Join(destDir, ".git"), false}, // .git should be excluded
		{filepath.Join(destDir, ".git", "config"), false},
	}

	for _, tt := range tests {
		_, err := os.Stat(tt.path)
		exists := err == nil
		if exists != tt.shouldExist {
			t.Errorf("Path %q exists = %v, want %v", tt.path, exists, tt.shouldExist)
		}
	}

	// Verify content
	content, err := os.ReadFile(filepath.Join(destDir, "file1.txt"))
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}
	if string(content) != "content1" {
		t.Errorf("File content = %q, want %q", string(content), "content1")
	}
}
