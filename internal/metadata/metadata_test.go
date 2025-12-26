package metadata

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMetadata_Generate(t *testing.T) {
	fixedTime := time.Date(2025, 12, 26, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name            string
		meta            *Metadata
		wantContains    []string
		wantNotContains []string
	}{
		{
			name: "with history preserved",
			meta: &Metadata{
				OriginalSource:   "https://github.com/owner/repo",
				BuriedAt:         fixedTime,
				HistoryPreserved: true,
			},
			wantContains: []string{
				"https://github.com/owner/repo",
				"2025-12-26T10:30:00Z",
				"**History Preserved** | Yes",
			},
		},
		{
			name: "without history preserved",
			meta: &Metadata{
				OriginalSource:   "/path/to/local/repo",
				BuriedAt:         fixedTime,
				HistoryPreserved: false,
			},
			wantContains: []string{
				"/path/to/local/repo",
				"2025-12-26T10:30:00Z",
				"**History Preserved** | No",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.meta.Generate()

			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("Generate() missing expected content: %q\n\nGot:\n%s", want, got)
				}
			}

			for _, notWant := range tt.wantNotContains {
				if strings.Contains(got, notWant) {
					t.Errorf("Generate() contains unexpected content: %q\n\nGot:\n%s", notWant, got)
				}
			}
		})
	}
}

func TestMetadata_Write(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "metadata-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tempDir) })

	meta := &Metadata{
		OriginalSource:   "https://github.com/owner/repo",
		BuriedAt:         time.Now(),
		HistoryPreserved: true,
	}

	// Write metadata
	if err := meta.Write(tempDir); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	// Verify file exists
	metaPath := filepath.Join(tempDir, FileName)
	if _, err := os.Stat(metaPath); os.IsNotExist(err) {
		t.Fatalf("Metadata file was not created")
	}

	// Verify content
	content, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("Failed to read metadata file: %v", err)
	}

	if !strings.Contains(string(content), "https://github.com/owner/repo") {
		t.Errorf("Metadata file missing expected content")
	}
}

func TestMetadata_Write_NonExistentDir(t *testing.T) {
	meta := &Metadata{
		OriginalSource:   "https://github.com/owner/repo",
		BuriedAt:         time.Now(),
		HistoryPreserved: true,
	}

	// Try to write to non-existent directory
	err := meta.Write("/path/that/does/not/exist")
	if err == nil {
		t.Errorf("Write() expected error for non-existent directory, got nil")
	}
}
