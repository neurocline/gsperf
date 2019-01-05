// gsperf/file_darwin.go

package main

import (
	"log"
	"os/exec"
)

// openNoBuffering on Mac OS X flushes the entire disk cache,
// because I don't know how to flush for a single file.
func openNoBuffering(path string) {
	path = ""
	purgeCmd := exec.Command("purge")
	results, err := purgeCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running purge: %s\n%s\n", err, results)
		fmt.Printf("You might need to run as admin with sudo\n")
	}
}
