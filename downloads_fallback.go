package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func fallbackDownloadsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine home directory: %w", err)
	}
	return filepath.Join(home, "Downloads"), nil
}
