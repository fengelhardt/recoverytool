package main

import(
	"fmt"
	"time"
)

func siValue(v float64, unit string) string {
	bigsiunits  := []string{"", "k", "M", "G", "T", "E"}
	smallsiunits  := []string{"", "m", "µ", "n", "p", "f"}
	maxorder := 5
	order := 0
	ret := ""
	if v < 0 {
		v = -v
		ret = "-"
	}
	if v < 1.0 {
		for v <= 1.0 && order < maxorder {
			order++
			v *= 1000.0
		}
		ret = fmt.Sprintf("%.1f%s%s", v, smallsiunits[order], unit)
	} else {
		for v >= 1000.0 && order < maxorder {
			order++
			v /= 1000.0
		}
		ret = fmt.Sprintf("%.1f%s%s", v, bigsiunits[order], unit)
	}
	return ret
}

func printHex(b byte) {
	if b >= ' ' && b <='~' {
		fmt.Printf("%c", b)
	} else if b == '\t' || b == '\r' || b == '\n' || b == '\f' {
		fmt.Print(" ")
	} else {
		fmt.Print(".")
	}
}

func printHexLine(offs uint64, buf []byte) {
	// same as hexdump -e '1/128  "%_ax "' -e '128/ "%_p"' -e '1 "\n"' <file>
	fmt.Printf("%x ", offs)
	for _, b := range buf {
		printHex(b)
	}
	fmt.Println()
}

func calcETA(amount, throughput float64) string {
	return time.Duration(amount/throughput*1000000000.0).String()
}
