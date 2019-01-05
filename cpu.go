// gsperf/cpu.go
// - CPU performance measurement

package main

import (
	"fmt"
	"time"
)

// ----------------------------------------------------------------------------------------------

// cpuIntPerf does some simple cpu integer performance measurement
// TBD multi-core version
func cpuIntPerf(enabled bool) {
	if !enabled {
		return
	}
	testsSelected += 1
	fmt.Printf("Testing CPU integer performance\n")

	buf := make([]byte, 32768)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i & 255)
	}
	var crcsum uint16 = 0xFFFF

	startTime := time.Now()
	for i := 0; i < 100; i++ {
		crcsum = ccitt_crc16(buf, crcsum)
	}
	delta1 := elapsed(startTime)

	iter_1sec := int(100 / delta1)
	niter := 5 * iter_1sec

	startTime = time.Now()
	for i := 0; i < niter-100; i++ {
		crcsum = ccitt_crc16(buf, crcsum)
	}
	delta2 := elapsed(startTime)

	crcsum ^= 0xFFFF
	crcsum = (crcsum << 8) | (crcsum >> 8)

	per_sec := int(float64(niter) / (delta1 + delta2))
	fmt.Printf("ccitt_crc16 %02X: %d/second\n", crcsum, per_sec)
}

// Go version of code in stress-ng
// https://github.com/ColinIanKing/stress-ng/blob/ca141b61aef72c627b471c285d61a52d13983864/stress-cpu.c
func ccitt_crc16(buf []byte, crc uint16) uint16 {
	const polynomial uint16 = 0x8408

	for i := 0; i < len(buf); i++ {
		var val uint16 = uint16(buf[i])
		for bit := 0; bit < 8; bit++ {
			do_xor := 1 & (val ^ crc)
			crc >>= 1
			if do_xor != 0 {
				crc ^= polynomial
			}
		}
	}

	return crc
}

// ----------------------------------------------------------------------------------------------

func cpuFloatPerf(enabled bool) {
	if !enabled {
		return
	}
	testsSelected += 1
	fmt.Printf("Testing CPU float performance\n")

	// TODO - write first float-intensive code
}
