package mover

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marce/download-cleaner/scanner"
)

// Move performs the actual filesystem moves for all non-skipped entries.
// It creates destination directories as needed and resolves filename conflicts.
// Returns the subset of entries that were successfully moved, with DestPath
// updated to reflect any conflict-resolution renaming.
func Move(entries []scanner.FileEntry) ([]scanner.FileEntry, error) {
	var moved []scanner.FileEntry
	var firstErr error

	for i := range entries {
		if entries[i].IsSkipped {
			continue
		}

		if err := ensureDir(entries[i].DestDir); err != nil {
			fmt.Fprintf(os.Stderr, "  [error] cannot create %s: %v\n", entries[i].DestDir, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		// Resolve any filename conflict before moving
		entries[i].DestPath = resolveConflict(entries[i].DestPath)

		if err := os.Rename(entries[i].SourcePath, entries[i].DestPath); err != nil {
			fmt.Fprintf(os.Stderr, "  [error] cannot move %s: %v\n", entries[i].Name, err)
			if firstErr == nil {
				firstErr = err
			}
			continue
		}

		moved = append(moved, entries[i])
	}

	return moved, firstErr
}

// ensureDir creates dir and any necessary parents if they do not exist.
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0o755)
}

// resolveConflict checks whether destPath already exists. If so, it appends
// _1, _2, … to the stem (before the extension) until a free slot is found.
// Handles .tar.gz double extensions correctly.
func resolveConflict(destPath string) string {
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		return destPath
	}

	dir := filepath.Dir(destPath)
	base := filepath.Base(destPath)

	// Determine stem and extension, handling .tar.gz
	var stem, ext string
	lowerBase := strings.ToLower(base)
	if strings.HasSuffix(lowerBase, ".tar.gz") {
		ext = base[len(base)-7:] // preserve original casing of ".tar.gz"
		stem = base[:len(base)-7]
	} else {
		rawExt := filepath.Ext(base)
		ext = rawExt
		stem = strings.TrimSuffix(base, rawExt)
	}

	for i := 1; ; i++ {
		candidate := filepath.Join(dir, fmt.Sprintf("%s_%d%s", stem, i, ext))
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
	}
}
