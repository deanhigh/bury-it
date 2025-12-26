package graveyard

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "absolute path",
			path:    "/tmp/graveyard",
			wantErr: false,
		},
		{
			name:    "relative path",
			path:    "./graveyard",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gy, err := New(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("New(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
				return
			}
			if !tt.wantErr && gy == nil {
				t.Errorf("New(%q) returned nil graveyard", tt.path)
			}
		})
	}
}

func TestGraveyard_Validate(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "graveyard-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	// Create a valid graveyard (git repo)
	validGraveyard := filepath.Join(tempDir, "valid-graveyard")
	if err := os.MkdirAll(filepath.Join(validGraveyard, ".git"), 0755); err != nil {
		t.Fatalf("Failed to create valid graveyard: %v", err)
	}

	// Create a non-git directory
	nonGitDir := filepath.Join(tempDir, "non-git")
	if err := os.MkdirAll(nonGitDir, 0755); err != nil {
		t.Fatalf("Failed to create non-git dir: %v", err)
	}

	// Create a file
	filePath := filepath.Join(tempDir, "a-file")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "valid git graveyard",
			path:    validGraveyard,
			wantErr: false,
		},
		{
			name:    "non-existent path",
			path:    filepath.Join(tempDir, "does-not-exist"),
			wantErr: true,
		},
		{
			name:    "path is a file",
			path:    filePath,
			wantErr: true,
		},
		{
			name:    "directory without .git",
			path:    nonGitDir,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gy, err := New(tt.path)
			if err != nil {
				t.Fatalf("New(%q) unexpected error: %v", tt.path, err)
			}
			err = gy.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGraveyard_ValidateProjectName(t *testing.T) {
	// Create temp graveyard
	tempDir, err := os.MkdirTemp("", "graveyard-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	// Create .git directory
	if err := os.MkdirAll(filepath.Join(tempDir, ".git"), 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Create an existing project
	existingProject := filepath.Join(tempDir, "existing-project")
	if err := os.MkdirAll(existingProject, 0755); err != nil {
		t.Fatalf("Failed to create existing project: %v", err)
	}

	gy, err := New(tempDir)
	if err != nil {
		t.Fatalf("Failed to create graveyard: %v", err)
	}

	tests := []struct {
		name        string
		projectName string
		wantErr     bool
	}{
		{
			name:        "valid new project name",
			projectName: "new-project",
			wantErr:     false,
		},
		{
			name:        "project already exists",
			projectName: "existing-project",
			wantErr:     true,
		},
		{
			name:        "empty name",
			projectName: "",
			wantErr:     true,
		},
		{
			name:        "dot",
			projectName: ".",
			wantErr:     true,
		},
		{
			name:        "double dot",
			projectName: "..",
			wantErr:     true,
		},
		{
			name:        "name with slash",
			projectName: "foo/bar",
			wantErr:     true,
		},
		{
			name:        "name with backslash",
			projectName: "foo\\bar",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gy.ValidateProjectName(tt.projectName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName(%q) error = %v, wantErr %v", tt.projectName, err, tt.wantErr)
			}
		})
	}
}

func TestGraveyard_ProjectPath(t *testing.T) {
	gy := &Graveyard{Path: "/path/to/graveyard"}

	got := gy.ProjectPath("my-project")
	want := "/path/to/graveyard/my-project"

	if got != want {
		t.Errorf("ProjectPath() = %q, want %q", got, want)
	}
}

func TestGraveyard_ProjectExists(t *testing.T) {
	// Create temp graveyard
	tempDir, err := os.MkdirTemp("", "graveyard-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	// Create an existing project
	existingProject := filepath.Join(tempDir, "existing")
	if err := os.MkdirAll(existingProject, 0755); err != nil {
		t.Fatalf("Failed to create existing project: %v", err)
	}

	gy := &Graveyard{Path: tempDir}

	tests := []struct {
		name        string
		projectName string
		want        bool
	}{
		{
			name:        "existing project",
			projectName: "existing",
			want:        true,
		},
		{
			name:        "non-existing project",
			projectName: "non-existing",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gy.ProjectExists(tt.projectName)
			if got != tt.want {
				t.Errorf("ProjectExists(%q) = %v, want %v", tt.projectName, got, tt.want)
			}
		})
	}
}
