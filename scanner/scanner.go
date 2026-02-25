package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FileEntry holds all metadata about a file discovered in the Downloads folder.
type FileEntry struct {
	Name       string    // original filename
	SourcePath string    // absolute path of the file
	DestDir    string    // absolute path of destination directory
	DestPath   string    // provisional destination path (may be updated by mover for conflicts)
	Category   string    // "Executables" | "Documents" | "Archives" | "Images" | ""
	ModTime    time.Time // file modification time
	IsOld      bool      // true if file age exceeds threshold
	IsSkipped  bool      // true if no recognised extension (left in place)
}

// categoryExtensions maps each category to its recognised lowercase extensions.
var categoryExtensions = map[string][]string{
	"Executables": {"exe", "msi", "bat"},
	"Documents":   {"pdf", "docx", "xlsx", "txt"},
	"Archives":    {"zip", "rar", "7z"}, // tar.gz handled via HasSuffix below
	"Images":      {"jpg", "jpeg", "png", "gif", "svg", "webp"},
}

// extToCategory is a flat lookup map built from categoryExtensions.
var extToCategory map[string]string

func init() {
	extToCategory = make(map[string]string)
	for cat, exts := range categoryExtensions {
		for _, ext := range exts {
			extToCategory[ext] = cat
		}
	}
}

// Scan reads the top-level entries of downloadsDir, categorises each file,
// and returns the full list of FileEntry values sorted by category then name.
// ageThreshold is the number of days after which a file is considered "old".
func Scan(downloadsDir string, ageThreshold int) ([]FileEntry, error) {
	entries, err := os.ReadDir(downloadsDir)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	threshold := time.Duration(ageThreshold) * 24 * time.Hour

	var results []FileEntry

	for _, de := range entries {
		// Skip subdirectories (no recursion)
		if de.IsDir() {
			continue
		}

		info, err := de.Info()
		if err != nil {
			continue
		}

		name := de.Name()
		sourcePath := filepath.Join(downloadsDir, name)
		modTime := info.ModTime()
		isOld := now.Sub(modTime) > threshold

		ext := resolveExtension(name)
		category := resolveCategory(ext)

		fe := FileEntry{
			Name:       name,
			SourcePath: sourcePath,
			ModTime:    modTime,
			IsOld:      isOld,
			Category:   category,
		}

		if isOld {
			// Old files go to Archive/ regardless of type
			fe.DestDir = filepath.Join(downloadsDir, "Archive")
			fe.DestPath = filepath.Join(fe.DestDir, name)
			fe.IsSkipped = false
		} else if category == "" {
			fe.IsSkipped = true
		} else {
			fe.DestDir = filepath.Join(downloadsDir, category)
			fe.DestPath = filepath.Join(fe.DestDir, name)
			fe.IsSkipped = false
		}

		results = append(results, fe)
	}

	// Sort: Archive/old first, then by category, then by name
	sort.Slice(results, func(i, j int) bool {
		ci := sortKey(results[i])
		cj := sortKey(results[j])
		if ci != cj {
			return ci < cj
		}
		return results[i].Name < results[j].Name
	})

	return results, nil
}

func sortKey(fe FileEntry) string {
	if fe.IsSkipped {
		return "z_skipped"
	}
	if fe.IsOld {
		return "a_archive"
	}
	return fe.Category
}

// resolveExtension returns the effective lowercase extension for a filename,
// handling the special case of ".tar.gz" double extensions.
func resolveExtension(name string) string {
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, ".tar.gz") {
		return "tar.gz"
	}
	ext := filepath.Ext(lower)
	return strings.TrimPrefix(ext, ".")
}

// resolveCategory returns the category for a given extension, or "" if unknown.
func resolveCategory(ext string) string {
	if ext == "tar.gz" {
		return "Archives"
	}
	return extToCategory[ext]
}
