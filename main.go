// gsperf/main.go
// Copyright 2019 Brian Fitzgerald <neurocline@gmail.com>
//
// Use of this source code is governed by an MIT-style license that can be found in the LICENSE file.
//
// gsperf is a simple cross-platform performance benchmarking tool

package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	args := parseArgs()

	cpuIntPerf(args.cpuInt)
	cpuFloatPerf(args.cpuFloat)

	diskPhysicalPerf(args.diskPhysical)
	diskCachePerf(args.diskCache)

	if testsSelected == 0 {
		fmt.Printf("No tests selected\n")
		usage(1)
	}
}

var testsSelected int

// ----------------------------------------------------------------------------------------------

type cmdArgs struct {
	all bool

	cpu      bool
	cpuInt   bool
	cpuFloat bool

	disk         bool
	diskPhysical bool
	diskCache    bool

	help    bool
	verbose bool
}

func usage(fail int) {
	fmt.Printf("Usage: perf [--cpu] [--cpu-int] [--cpu-float]\n" +
		"            [--disk] [--disk-physical] [--disk-cache]\n" +
		"            [--all] [-v|--verbose] [-h|--help]\n")
	os.Exit(fail)
}

// Process command-line
func parseArgs() *cmdArgs {
	p := cmdArgs{}

	for _, arg := range os.Args[1:] {

		parsebool := func(opt string, val *bool) bool {
			if arg != opt {
				return false
			}
			*val = true
			return true
		}

		if !parsebool("--cpu-int", &p.cpuInt) &&
			!parsebool("--cpu-float", &p.cpuFloat) &&
			!parsebool("--cpu", &p.cpu) &&
			!parsebool("--disk-physical", &p.diskPhysical) &&
			!parsebool("--disk-cache", &p.diskCache) &&
			!parsebool("--disk", &p.disk) &&
			!parsebool("--all", &p.all) &&
			!parsebool("-v", &p.verbose) &&
			!parsebool("--verbose", &p.verbose) &&
			!parsebool("-h", &p.help) &&
			!parsebool("--help", &p.help) {
			usage(1)
		}
	}

	if p.help {
		usage(0)
	}

	// --all is shorthand for all tests
	if p.all {
		p.cpu = true
		p.disk = true
	}

	// --cpu is shorthand for all cpu tests
	if p.cpu {
		p.cpuInt = true
		p.cpuFloat = true
	}

	// --disk is shorthand for all disk tests
	if p.disk {
		p.diskPhysical = true
		p.diskCache = true
	}

	return &p
}

// elapsed is intended to be used as a defer in a function, to
// measure the time spent in that function
//   defer elapsed(time.Now())
func elapsed(start time.Time) float64 {
	return time.Since(start).Seconds()
}
