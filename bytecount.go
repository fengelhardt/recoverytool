package main

import (
	"fmt"
	"time"
)

type bytecountaction struct {
	cnt uint64
	tp float64
	t1 time.Time
}

func (a *bytecountaction) Init() {
	a.t1 = time.Now()
	a.tp = 0.0
	return
}

func (a *bytecountaction) Run(buf[] byte, abspos, tcnt, n uint64, lastbit bool) (uint64, error) {
	for i := 0; i < len(buf) ; i++ {
		a.cnt++
		tcnt++
		abspos++
	}
	t2 := time.Now()
	itp := float64(len(buf)) / float64((t2.Sub(a.t1)).Milliseconds()) * 1000.0
	a.tp = 0.9*a.tp + 0.1*itp
	eta := calcETA(float64(n-tcnt), a.tp)
	a.t1 = t2
	percentcomplete := float64(tcnt)/float64(n)*100.0
	fmt.Printf("%5.1f%%, %s ETA %s          \r", percentcomplete, siValue(a.tp, "B/s"), eta)
	if lastbit {
		fmt.Println()
		fmt.Println("Counted ", a.cnt)
	}
	return uint64(len(buf)), nil
}
