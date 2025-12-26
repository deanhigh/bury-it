// Package metadata handles generation of .bury-it.md files.
package metadata

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Metadata contains information about an archived project.
type Metadata struct {
	// OriginalSource is the original source location.
	OriginalSource string
	// BuriedAt is the timestamp when the project was buried.
	BuriedAt time.Time
	// HistoryPreserved indicates whether git history was preserved.
	HistoryPreserved bool
}

// FileName is the name of the metadata file.
const FileName = ".bury-it.md"

// Generate generates the metadata content as a string.
func (m *Metadata) Generate() string {
	historyStr := "Yes"
	if !m.HistoryPreserved {
		historyStr = "No"
	}

	return fmt.Sprintf(`# Archived Project

| Field | Value |
|-------|-------|
| **Original Source** | %s |
| **Buried On** | %s |
| **History Preserved** | %s |

---

*This project was archived using [bury-it](https://github.com/deanhigh/bury-it).*
`, m.OriginalSource, m.BuriedAt.Format(time.RFC3339), historyStr)
}

// Write writes the metadata file to the specified directory.
func (m *Metadata) Write(dir string) error {
	filePath := filepath.Join(dir, FileName)
	content := m.Generate()
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}
	return nil
}
