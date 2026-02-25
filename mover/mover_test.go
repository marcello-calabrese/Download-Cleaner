package mover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConflictReturnsOriginalWhenFree(t *testing.T) {
	tempDir := t.TempDir()
	dest := filepath.Join(tempDir, "report.pdf")

	got, err := resolveConflict(dest)
	if err != nil {
		t.Fatalf("resolveConflict returned error: %v", err)
	}
	if got != dest {
		t.Fatalf("expected %q, got %q", dest, got)
	}
}

func TestResolveConflictFindsNextAvailableSuffix(t *testing.T) {
	tempDir := t.TempDir()
	dest := filepath.Join(tempDir, "report.pdf")
	if err := os.WriteFile(dest, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed creating %q: %v", dest, err)
	}
	firstCandidate := filepath.Join(tempDir, "report_1.pdf")
	if err := os.WriteFile(firstCandidate, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed creating %q: %v", firstCandidate, err)
	}

	got, err := resolveConflict(dest)
	if err != nil {
		t.Fatalf("resolveConflict returned error: %v", err)
	}
	want := filepath.Join(tempDir, "report_2.pdf")
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestResolveConflictPreservesTarGzExtension(t *testing.T) {
	tempDir := t.TempDir()
	dest := filepath.Join(tempDir, "backup.tar.gz")
	if err := os.WriteFile(dest, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed creating %q: %v", dest, err)
	}

	got, err := resolveConflict(dest)
	if err != nil {
		t.Fatalf("resolveConflict returned error: %v", err)
	}
	want := filepath.Join(tempDir, "backup_1.tar.gz")
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
