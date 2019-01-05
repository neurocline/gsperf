// file code for windows
// -- Windows-specific file code

package main

import (
	"golang.org/x/sys/windows"
)

// openNoBuffering on Windows opens the file with FILE_FLAG_NO_BUFFERING,
// which has the side-effect of removing all cached data relevant to this
// file.
func openNoBuffering(path string) {

	// For some reason, the Go team didn't think this flag was worth exposing
	const (
		FILE_FLAG_NO_BUFFERING = 0x20000000
	)

	// CreateFile parameters
	pathp, _ := windows.UTF16PtrFromString(path)
	var access uint32 = windows.GENERIC_READ
	sharemode := uint32(windows.FILE_SHARE_READ | windows.FILE_SHARE_WRITE)
	//var sa *windows.SecurityAttributes
	createmode := uint32(windows.OPEN_EXISTING)
	flags := uint32(FILE_FLAG_NO_BUFFERING)

	// Open the file, which wipes the cache, and close it, because we'll
	// use stock Go routines after this
	h, _ := windows.CreateFile(pathp, access, sharemode, nil, createmode, flags, 0)
	if h != windows.InvalidHandle {
		windows.CloseHandle(h)
	}
}
