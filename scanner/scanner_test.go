package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestResolveExtension(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     string
	}{
		{name: "regular extension", fileName: "document.PDF", want: "pdf"},
		{name: "double extension tar.gz", fileName: "archive.TAR.GZ", want: "tar.gz"},
		{name: "no extension", fileName: "README", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := resolveExtension(tc.fileName)
			if got != tc.want {
				t.Fatalf("resolveExtension(%q): want %q, got %q", tc.fileName, tc.want, got)
			}
		})
	}
}

func TestResolveCategory(t *testing.T) {
	tests := []struct {
		ext  string
		want string
	}{
		{ext: "pdf", want: "Documents"},
		{ext: "tar.gz", want: "Archives"},
		{ext: "unknown", want: ""},
	}

	for _, tc := range tests {
		got := resolveCategory(tc.ext)
		if got != tc.want {
			t.Fatalf("resolveCategory(%q): want %q, got %q", tc.ext, tc.want, got)
		}
	}
}

func TestScanOldAndUnknownClassification(t *testing.T) {
	tempDir := t.TempDir()

	oldFile := filepath.Join(tempDir, "old.zip")
	if err := os.WriteFile(oldFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed creating old file: %v", err)
	}
	oldTime := time.Now().Add(-48 * time.Hour)
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatalf("failed setting old file timestamp: %v", err)
	}

	unknownFile := filepath.Join(tempDir, "notes.xyz")
	if err := os.WriteFile(unknownFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed creating unknown file: %v", err)
	}

	entries, err := Scan(tempDir, 1)
	if err != nil {
		t.Fatalf("Scan returned error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	byName := make(map[string]FileEntry, len(entries))
	for _, entry := range entries {
		byName[entry.Name] = entry
	}

	oldEntry, ok := byName["old.zip"]
	if !ok {
		t.Fatalf("missing entry for old.zip")
	}
	if !oldEntry.IsOld {
		t.Fatalf("expected old.zip to be marked old")
	}
	if oldEntry.DestDir != filepath.Join(tempDir, "Archive") {
		t.Fatalf("expected old.zip destination dir %q, got %q", filepath.Join(tempDir, "Archive"), oldEntry.DestDir)
	}

	unknownEntry, ok := byName["notes.xyz"]
	if !ok {
		t.Fatalf("missing entry for notes.xyz")
	}
	if !unknownEntry.IsSkipped {
		t.Fatalf("expected notes.xyz to be skipped")
	}
}
