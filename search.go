package main

import (
	"fmt"
	"time"
	"os"
)

type searchaction struct {
	tp float64
	t1 time.Time
	needles []string
	npos []int
	n int
}

func (a *searchaction) Init(needles []string) {
	a.t1 = time.Now()
	a.tp = 0.0
	a.n = len(needles)
	a.needles = needles
	a.npos = make([]int, a.n)
	return
}

func (a *searchaction) match(c byte, needle int, addr uint64) {
	if a.needles[needle][a.npos[needle]] == c {
		a.npos[needle]++
		if a.npos[needle] == len(a.needles[needle]) {
			a.report(needle, addr)
			a.npos[needle] = 0
		}
	} else {
		a.npos[needle] = 0
	}
	return
}

func (a *searchaction) report(needle int, addr uint64) {
	fmt.Printf("Match for %s at adress %x\n", a.needles[needle], addr)
	rstartaddr := addr - addr%g_linelen
	rstartline := rstartaddr / g_linelen
// 	fmt.Println(rstartaddr, g_linelen, rstartline)
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
	for i := 0; i < len(buf) ; i++ {
		c := buf[i]
		for j, _ := range a.needles {
			a.match(c, j, abspos)
		}
		tcnt++
		abspos++
	}
	t2 := time.Now()
	itp := float64(len(buf)) / float64((t2.Sub(a.t1)).Milliseconds()) * 1000.0
	a.tp = 0.9*a.tp + 0.1*itp
	eta := calcETA(float64(n-tcnt), a.tp)
	a.t1 = t2
	percentcomplete := float64(tcnt)/float64(n)*100.0
	fmt.Fprintf(os.Stderr, "\r%5.1f%%, %s ETA %s          ", percentcomplete, siValue(a.tp, "B/s"), eta)
	if lastbit {
		fmt.Fprintf(os.Stderr, "\n")
	}
	return uint64(len(buf)), nil
}
