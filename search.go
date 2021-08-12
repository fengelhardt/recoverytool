package main

import (
	"time"
	"bytes"
)

type searchaction struct {
	tp float64
	t1 time.Time
	needles [][]byte
	maxlen int
}

func (a *searchaction) Init(needles [][]byte) {
	a.t1 = time.Now()
	a.tp = 0.0
	a.needles = needles
	a.maxlen = 0
	for _ , n := range needles {
		if len(n) > a.maxlen {a.maxlen = len(n)}
	}
	uiPrintf1(2, "Searching for the following patterns: %s\n", a.needles)
	return
}

func (a *searchaction) report(needle []byte, addr uint64) {
	uiPrintf1(0, "Match for %s at adress %x\n", needle, addr)
	rstartaddr := addr - addr%g_linelen
	rstartline := rstartaddr / g_linelen
	rnlines := g_linesbefore + g_linesafter + 1
	if rstartaddr > g_linesbefore*g_linelen {
		rstartaddr -= g_linesbefore*g_linelen
	} else { 
		rnlines = rstartline + g_linesafter + 1
		rstartaddr = uint64(0)
	}
	printHexLines(rstartaddr, rnlines)
	return
}

func (a *searchaction) Run(buf[] byte, abspos, tcnt, n uint64, lastbit bool) (uint64, error) {
	if len(buf) < 0 {
		return uint64(0), nil
	}
	for _, needle := range a.needles {
		bufcnt := 0
		idx := bytes.Index(buf, needle)
		for idx != -1 {
			bufcnt += idx
			a.report(needle, abspos+uint64(bufcnt))
			bufcnt += len(needle)
			idx = bytes.Index(buf[bufcnt:], needle)
		}
	}
	t2 := time.Now()
	itp := float64(len(buf)) / float64((t2.Sub(a.t1)).Milliseconds()) * 1000.0
	a.tp = 0.9*a.tp + 0.1*itp
	eta := calcETA(float64(n-tcnt), a.tp)
	a.t1 = t2
	percentcomplete := float64(tcnt)/float64(n)*100.0
	uiPrintf1(1, "%5.1f%%, %s ETA %s          \r", percentcomplete, siValue(a.tp, "B/s"), eta)
	if lastbit {
		uiPrintf2(1, "\n")
		return uint64(len(buf)), nil
	}
	return uint64(len(buf)-a.maxlen), nil
}
