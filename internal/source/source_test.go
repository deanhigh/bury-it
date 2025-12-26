package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantType    Type
		wantName    string
		wantPathSfx string // suffix to check for path (for URLs) or empty for local
		wantErr     bool
	}{
		{
			name:        "github url with https",
			input:       "https://github.com/owner/repo",
			wantType:    TypeRemote,
			wantName:    "repo",
			wantPathSfx: "https://github.com/owner/repo",
		},
		{
			name:        "github url with .git suffix",
			input:       "https://github.com/owner/repo.git",
			wantType:    TypeRemote,
			wantName:    "repo",
			wantPathSfx: "https://github.com/owner/repo.git",
		},
		{
			name:        "github url with trailing slash",
			input:       "https://github.com/owner/repo/",
			wantType:    TypeRemote,
			wantName:    "repo",
			wantPathSfx: "https://github.com/owner/repo/",
		},
		{
			name:        "owner/repo shorthand",
			input:       "deanhigh/bury-it",
			wantType:    TypeRemote,
			wantName:    "bury-it",
			wantPathSfx: "https://github.com/deanhigh/bury-it",
		},
		{
			name:        "owner/repo with dots and dashes",
			input:       "some-org/my.project-name",
			wantType:    TypeRemote,
			wantName:    "my.project-name",
			wantPathSfx: "https://github.com/some-org/my.project-name",
		},
		{
			name:     "relative path with dot",
			input:    "./my-project",
			wantType: TypeLocal,
			wantName: "my-project",
		},
		{
			name:     "absolute path",
			input:    "/tmp/my-project",
			wantType: TypeLocal,
			wantName: "my-project",
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, err := Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse(%q) unexpected error: %v", tt.input, err)
				return
			}

			if src.Type != tt.wantType {
				t.Errorf("Parse(%q) Type = %v, want %v", tt.input, src.Type, tt.wantType)
			}

			if src.Name != tt.wantName {
				t.Errorf("Parse(%q) Name = %q, want %q", tt.input, src.Name, tt.wantName)
			}

			if tt.wantPathSfx != "" && src.Path != tt.wantPathSfx {
				t.Errorf("Parse(%q) Path = %q, want %q", tt.input, src.Path, tt.wantPathSfx)
			}
		})
	}
}

func TestSource_Validate(t *testing.T) {
	// Create a temporary directory to simulate repos
	tempDir, err := os.MkdirTemp("", "source-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	// Create a valid git repo
	validRepo := filepath.Join(tempDir, "valid-repo")
	if err := os.MkdirAll(filepath.Join(validRepo, ".git"), 0755); err != nil {
		t.Fatalf("Failed to create valid repo: %v", err)
	}

	// Create a non-git directory
	nonGitDir := filepath.Join(tempDir, "non-git")
	if err := os.MkdirAll(nonGitDir, 0755); err != nil {
		t.Fatalf("Failed to create non-git dir: %v", err)
	}

	// Create a file (not a directory)
	filePath := filepath.Join(tempDir, "a-file")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	tests := []struct {
		name    string
		source  *Source
		wantErr bool
	}{
		{
			name: "valid local git repo",
			source: &Source{
				Type: TypeLocal,
				Path: validRepo,
			},
			wantErr: false,
		},
		{
			name: "non-existent path",
			source: &Source{
				Type: TypeLocal,
				Path: filepath.Join(tempDir, "does-not-exist"),
			},
			wantErr: true,
		},
		{
			name: "path is a file not directory",
			source: &Source{
				Type: TypeLocal,
				Path: filePath,
			},
			wantErr: true,
		},
		{
			name: "directory without .git",
			source: &Source{
				Type: TypeLocal,
				Path: nonGitDir,
			},
			wantErr: true,
		},
		{
			name: "remote type skips local validation",
			source: &Source{
				Type: TypeRemote,
				Path: "https://github.com/owner/repo",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
