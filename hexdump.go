package main

import (
)

type hexdumpaction struct {
	// does not need specific data
}

func (a *hexdumpaction) Run(buf[] byte, abspos, tcnt, n uint64, lastbit bool) (uint64, error) {
	uiPrintf2(3, "n=%d, tcnt=%d, buflen=%d, linelen=%d\n", n, tcnt, len(buf), g_linelen)
	lines := uint64(len(buf)) / g_linelen
	bufbytecnt := uint64(0)
	if lines > 0 {
		// handle complete lines
		lines := uint64(len(buf)) / g_linelen
		for i := uint64(0); i < lines ; i++ {
			offs1 := bufbytecnt
			offs2 := bufbytecnt + g_linelen
			line := buf[offs1:offs2]
			printHexLine(abspos, line)
			bufbytecnt += offs2-offs1
			tcnt += offs2-offs1
			abspos += offs2-offs1
		}
	} 
	if lastbit && bufbytecnt < uint64(len(buf)) {
		printHexLine(abspos, buf[bufbytecnt:])
		return uint64(len(buf)), nil
	} else {
		// leave the rest for the next call
		return bufbytecnt, nil
	}
}
