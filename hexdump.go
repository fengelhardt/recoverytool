package main

import (
// 	"fmt"
)

type hexdumpaction struct {
	startaddr uint64
}

func (a *hexdumpaction) Init(startaddr uint64) {
	a.startaddr = startaddr
	return
}

func (a *hexdumpaction) Run(buf[] byte, tcnt uint64, n uint64, lastbit bool) (uint64, error) {
// 	fmt.Printf("n=%d, tcnt=%d, buflen=%d, linelen=%d\n", n, tcnt, len(buf), g_linelen)
	lines := uint64(len(buf)) / g_linelen
	bufbytecnt := uint64(0)
	if lines > 0 {
		// handle complete lines
		lines := uint64(len(buf)) / g_linelen
		for i := uint64(0); i < lines ; i++ {
			offs1 := bufbytecnt
			offs2 := bufbytecnt + g_linelen
			line := buf[offs1:offs2]
			printHexLine(a.startaddr + tcnt, line)
			bufbytecnt += offs2-offs1
			tcnt += offs2-offs1
		}
	} 
	if lastbit && bufbytecnt < uint64(len(buf)) {
		printHexLine(a.startaddr + tcnt, buf[bufbytecnt:])
		return uint64(len(buf)), nil
	} else {
		// leave the rest for the next call
		return bufbytecnt, nil
	}
}
