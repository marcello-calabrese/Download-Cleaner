//go:build !windows

package main

func defaultDownloadsDir() (string, error) {
	return fallbackDownloadsDir()
}
