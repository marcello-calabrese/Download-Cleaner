//go:build windows

package main

import (
	"fmt"
	"unsafe"

	"syscall"
)

type knownFolderID struct {
	Data1 uint32
	Data2 uint16
	Data3 uint16
	Data4 [8]byte
}

var folderIDDownloads = knownFolderID{
	Data1: 0x374DE290,
	Data2: 0x123F,
	Data3: 0x4565,
	Data4: [8]byte{0x91, 0x64, 0x39, 0xC4, 0x92, 0x5E, 0x46, 0x7B},
}

var (
	shell32                  = syscall.NewLazyDLL("shell32.dll")
	ole32                    = syscall.NewLazyDLL("ole32.dll")
	procSHGetKnownFolderPath = shell32.NewProc("SHGetKnownFolderPath")
	procCoTaskMemFree        = ole32.NewProc("CoTaskMemFree")
)

func defaultDownloadsDir() (string, error) {
	var outPath uintptr
	hr, _, callErr := procSHGetKnownFolderPath.Call(
		uintptr(unsafe.Pointer(&folderIDDownloads)),
		0,
		0,
		uintptr(unsafe.Pointer(&outPath)),
	)

	if hr != 0 || outPath == 0 {
		fallback, fallbackErr := fallbackDownloadsDir()
		if fallbackErr != nil {
			if callErr != nil && callErr != syscall.Errno(0) {
				return "", fmt.Errorf("known folder lookup failed (hr=0x%x, err=%v) and fallback failed: %w", hr, callErr, fallbackErr)
			}
			return "", fmt.Errorf("known folder lookup failed (hr=0x%x) and fallback failed: %w", hr, fallbackErr)
		}
		return fallback, nil
	}

	defer procCoTaskMemFree.Call(outPath)

	path := utf16PtrToString((*uint16)(unsafe.Pointer(outPath)))
	if path == "" {
		return fallbackDownloadsDir()
	}

	return path, nil
}

func utf16PtrToString(ptr *uint16) string {
	if ptr == nil {
		return ""
	}

	chars := make([]uint16, 0, 260)
	for p := uintptr(unsafe.Pointer(ptr)); ; p += 2 {
		ch := *(*uint16)(unsafe.Pointer(p))
		if ch == 0 {
			break
		}
		chars = append(chars, ch)
	}

	return syscall.UTF16ToString(chars)
}
