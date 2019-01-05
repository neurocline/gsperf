# gsperf - System performance measurement

gsperf is a system performance tool, measuring CPU and disk, with more tests
planned. This tool is aimed primarily at developers who want to correlate
their own performance tests against the hardware being tested on.

## Getting started

Getting started is as easy as having a Go compiler installed on
your system and:

```
$ git clone git@github.com:neurocline/gsgit.git
$ cd gsgit
$ go build .
$ ./gsperf [options...]
```

With Go 1.11 or later, modules means you no longer need to have code
in the GOPATH. Hurray!

Run some performance tests:

```
$ ./gsperf --cpu-int
Testing CPU integer performance
ccitt_crc16 435B: 4782/second

$ ./gsperf --disk-cache
Testing disk cache performance
Doing 1K reads...686028 in 2.05 sec
Doing 4K reads...611042 in 2.17 sec
Doing 16K reads...381888 in 1.78 sec
Doing 64K reads...186190 in 1.69 sec
Doing 256K reads...73096 in 1.99 sec
1K reads: 328.91 MB/sec
4K reads: 1106.17 MB/sec
16K reads: 3388.23 MB/sec
64K reads: 6946.75 MB/sec
256K reads: 9218.50 MB/sec
```

## Mac version requires sudo at the moment

Physical disk performance tests are challenging, because modern operating systems
are good at caching, because that's one of the keys to being fast, or even usable.

However, for the Mac version, at the moment, you need to run gsperf with sudo if you want
to do the `--disk-physical` test, e.g. 

```
$ sudo ./gsperf --disk-physical
```

This is because I use the `usr/bin/purge` tool to reset the disk cache, and that
requires admin access. This is temporary, but the fix is a fair amount of code.

## Why this project?

There are a million of these, but this is slightly unique. It's written
in Go, and the mainstream Go compiler targets many operating systems with
the same compiler. Currently, most Linux, Mac, and Windows machines use
Intel processors. This means that gsperf will be compiled to the same code
(more or less) on all these machines, and thus gsperf won't introduce its
own differences into the measurements. As long as you use the same version
of gsperf on all the machines you want to measure, you'll get a set of
results that can be directly compared with each other.

Note that Go is not unique in this regard; you could probably accomplish
the same with Rust or D, and now that Clang can build Windows console programs,
even with C++. But Go is the easist of all of those languages to work with.

## What's next?

This is the bare-bones version that I needed when I was working on a cross-platform
Go project and worrying about performance on different machines. It was enough
to be able to compare performance for a body of code that does integer computation
and reads/writes files.

A list of what's queued up to be added:

- [ ] Linux version
- [ ] Mac version that doesn't need to run sudo
- [ ] Prettier output
- [ ] Memory bandwidth
- [ ] CPU floating-point performance
- [ ] Simple GPU performance
