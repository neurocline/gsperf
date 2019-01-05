// gsperf/disk.go
// - Disk performance measurement

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// ----------------------------------------------------------------------------------------------

// diskPhysicalPerf writes a large file and randomly seeks in it, avoiding
// re-reading blocks. On Windows, we actually force the disk cache to be flushed
// before we read, whereas on Mac we have to call an external tool to do that.
func diskPhysicalPerf(enabled bool) {
	if !enabled {
		return
	}
	testsSelected += 1
	fmt.Printf("Testing physical disk performance\n")

	// Create large test file
	// Because this file is so large, we don't delete it at the end, and
	// so we don't create it if it already exists
	path := "templargetestfile"
	const bigFileSize = 8 * 1024 * 1024 * 1024
	fh, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Creating large file (will not be deleted)...")
		fh = createTestFile(path, bigFileSize)
		fmt.Fprintf(os.Stderr, "\n")
	}

	// Read in various chunk sizes
	chunkSizes := []int{1024, 4096, 16384, 65536, 4 * 65536, 1024 * 1024, 4 * 1024 * 1024}
	chunkTags := []string{"1K", "4K", "16K", "64K", "256K", "1M", "4M"}
	bytespersec := make([]int, len(chunkSizes))

	for cs := 0; cs < len(chunkSizes); cs++ {
		fmt.Printf("Doing %s reads...", chunkTags[cs])
		chunkSize := chunkSizes[cs]
		var deltaTime float64

		// Cause the file in question to be flushed from the disk cache
		// (we may end up clearing the entire disk cache, on some operating systems)
		fmt.Fprintf(os.Stderr, "flushing...")
		fh.Close()
		openNoBuffering(path)
		fh, err = os.Open(path)
		if err != nil {
			log.Fatalf("Couldn't open %s: %s\n", path, err)
		}
		var probe int
		readblocks := make([]int, bigFileSize/chunkSize)
		fmt.Fprintf(os.Stderr, "reading...")

		// Read random chunks from the file until our desired
		// time has elapsed. This has the drawback of timing
		// the seek as well as the read, but the read should
		// dominate, presumably all seek does is set a value
		buf := make([]byte, chunkSize)
		startTime := time.Now()
		numReads := 0
		for {
			deltaTime = elapsed(startTime)
			if deltaTime >= 5.0 {
				break
			}

			// seek somewhere in the file, aligned by the chunk size
			seekBlock := rand.Intn(bigFileSize / chunkSize)
			for readblocks[seekBlock] != 0 {
				seekBlock += 5
				if seekBlock >= bigFileSize/chunkSize {
					seekBlock = 0
				}
				probe += 1
			}
			readblocks[seekBlock] = 1
			startPos := int64(seekBlock) * int64(chunkSize)
			_, err := fh.Seek(startPos, 0)
			if err != nil {
				log.Fatalf("Couldn't seek to pos %d in %s: %s\n", startPos, path, err)
			}

			// read a chunk
			amt, err := fh.Read(buf)
			if err != nil || amt != len(buf) {
				log.Fatalf("Couldn't read %s: %s\n", path, err)
			}
			numReads += 1
		}

		bytes_per_sec := int(float64(numReads*chunkSize) / deltaTime)
		bytespersec[cs] = bytes_per_sec
		fmt.Printf("%d in %.2f sec (%d)\n", numReads, deltaTime, probe)
	}

	// Clean up (leave large file behind for further runs)
	fh.Close()

	// Show results
	for i := 0; i < len(chunkSizes); i++ {
		fmt.Printf("%s reads: %.2f MB/sec\n", chunkTags[i], float64(bytespersec[i])/(1024.0*1024.0))
	}
}

// ----------------------------------------------------------------------------------------------

// diskCachePerf does repeated reads of data that should stay in the OS
// filesystem cache
func diskCachePerf(enabled bool) {
	if !enabled {
		return
	}
	testsSelected += 1
	fmt.Printf("Testing disk cache performance\n")

	// Create a small-ish temp file
	path := "temptestfile"
	fh := createTestFile(path, 4*1024*1024)

	// Start by reading the entire file, just to make sure it's in the cache
	readfromWrapped(1024*1024, 4, fh, path)

	// Read in various chunk sizes
	chunkSizes := []int{1024, 4096, 16384, 65536, 4 * 65536}
	chunkTags := []string{"1K", "4K", "16K", "64K", "256K"}
	bytespersec := make([]int, len(chunkSizes))

	for cs := 0; cs < len(chunkSizes); cs++ {
		beginTime := time.Now()
		fmt.Printf("Doing %s reads...", chunkTags[cs])
		chunkSize := chunkSizes[cs]
		_, err := fh.Seek(0, 0)
		if err != nil {
			log.Fatalf("Couldn't seek to start %s: %s\n", path, err)
		}

		// Do a few reads to see how many make a second
		// (also to put data in cache)
		smallIter := 128
		var delta1 float64
		for {
			delta1 = readfromWrapped(chunkSize, smallIter, fh, path)
			if delta1 > 0.01 {
				break
			}
			smallIter = smallIter * 2
		}

		// Now do a bunch of reads
		iter_1sec := int(float64(smallIter) / delta1)
		niter := 2 * iter_1sec

		delta2 := readfromWrapped(chunkSize, niter-smallIter, fh, path)

		per_sec := int(float64(niter) / (delta1 + delta2))
		bytespersec[cs] = chunkSize * per_sec
		fmt.Printf("%d in %.2f sec\n", niter, elapsed(beginTime))
	}

	fh.Close()
	err := os.Remove(path)
	if err != nil {
		log.Fatalf("Failed to delete %s: %s\n", path, err)
	}

	// Show results
	for i := 0; i < len(chunkSizes); i++ {
		fmt.Printf("%s reads: %.2f MB/sec\n", chunkTags[i], float64(bytespersec[i])/(1024.0*1024.0))
	}
}

// ----------------------------------------------------------------------------------------------

// readfromWrapped repeatedly reads from the same file, wrapping to beginning as needed.
// This assumes the buffer is modulo chunksize. To fix this, pass
// in file size so we can track how close we are.
func readfromWrapped(chunksize int, iterations int, fh *os.File, path string) float64 {
	_, err := fh.Seek(0, 0)
	if err != nil {
		log.Fatalf("Couldn't seek to start %s: %s\n", path, err)
	}

	buf := make([]byte, chunksize)
	startTime := time.Now()
	for i := 0; i < iterations; i++ {
		amt, err := fh.Read(buf)
		if amt == 0 {
			_, err = fh.Seek(0, 0)
			if err != nil {
				log.Fatalf("Couldn't seek to start %s: %s\n", path, err)
			}
			i -= 1
		}
		if err != nil {
			log.Fatalf("Couldn't read %s: %s\n", path, err)
		}
	}
	return elapsed(startTime)
}

// createTestFile creates a test file of the indicated size
func createTestFile(path string, filesize int) *os.File {
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Couldn't create %s: %s\n", path, err)
	}

	// Write in chunks of 1 MB - show status every 256 MB
	// TODO: make contents non-zero
	buf := make([]byte, 1024*1024)
	remain := filesize
	written := 0
	for pos := 0; pos < filesize; {
		if len(buf) > remain {
			buf = buf[:remain]
		}
		wrote, err := f.Write(buf)
		if err != nil || wrote != len(buf) {
			log.Fatalf("Write err %d/%d to %s: %s\n", wrote, len(buf), path, err)
		}

		pos += wrote
		remain -= wrote

		written += wrote
		if written >= 256*1024*1024 {
			fmt.Fprintf(os.Stderr, ".")
			written -= 256 * 1024 * 1024
		}
	}

	f.Close()

	// Now re-open the file as a courtesy for the caller
	f, err = os.Open(path)
	if err != nil {
		log.Fatalf("Couldn't open %s: %s\n", path, err)
	}
	return f
}
