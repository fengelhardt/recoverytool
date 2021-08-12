package main

import (
	"os"
	"fmt"
	"flag"
	"errors"
	"io"
)

var g_startfrom int64
var g_enduntil int64
var g_numlines uint64
var g_linelen uint64
var g_patterns strlist
var g_filename string
var g_bufsize uint
var g_doprint bool
var g_linesbefore uint64
var g_linesafter uint64
var g_domd5 bool
var g_verbose bool
var g_quiet bool
var g_debug bool

var g_verbosity int
var g_file *os.File

type strlist []string

func (l *strlist) String() string {
    return fmt.Sprintf("\"%s\"", *l)
}
 
func (l *strlist) Set(value string) error {
	*l = append(*l, value)
    return nil
}

func specifyFlags() {
	flag.Int64Var(&g_startfrom, "s", 0, "begin at this byte offset")
	flag.Int64Var(&g_enduntil, "e", 0, "read up to this byte offset (not including the particular byte at offset)")
	flag.Uint64Var(&g_numlines, "n", 0, "number of \"lines\" to process")
	flag.Uint64Var(&g_linelen, "l", 128, "length of a \"line\" for report")
	flag.Uint64Var(&g_linesbefore, "lb", 3, "\"lines\" of context to report before a match")
	flag.Uint64Var(&g_linesafter, "la", 3, "\"lines\" of context to report after a match")
	flag.UintVar(&g_bufsize, "b", 128*1024*1024, "buffer size used")
	flag.BoolVar(&g_doprint, "p", false, "print out data")
	flag.BoolVar(&g_domd5, "md5", false, "calculate the md5 checksum")
	flag.Var(&g_patterns, "m", "search for a pattern")
	flag.BoolVar(&g_verbose, "v", false, "be more verbose")
	flag.BoolVar(&g_quiet, "q", false, "do not print status updates")
	flag.BoolVar(&g_debug, "d", false, "print debug output")
}

func checkFlags() error {
	filesize, err := g_file.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	if g_startfrom < 0 {
		g_startfrom = filesize + g_startfrom
	}
	if g_enduntil <= 0 {
		g_enduntil = filesize + g_enduntil
	}
	if g_startfrom >= filesize {
		return errors.New("Start offset is too big")
	}
	if g_enduntil > filesize {
		return errors.New("End offset is too big")
	}
	if g_enduntil < g_startfrom {
		return errors.New("End offset is smaller than start offset")
	}
	if uint64(g_bufsize) < g_linelen {
		return errors.New("Buffer size must be bigger than the line size")
	}
	actions := 0
	if g_doprint { actions++ }
	if g_domd5 { actions++ }
	if len(g_patterns) != 0 { actions++ }
	if actions == 0 {
		return errors.New("No action specified. Must specify one of -p -m -md5")
	}
	if actions >= 2 {
		return errors.New("Can only have one flag out of -p -m -md5")
	}
	g_verbosity = 0
	if !g_quiet { g_verbosity = 1 }
	if g_verbose { g_verbosity = 2 }
	if g_debug { g_verbosity = 3 }
	return nil
}

func min(a, b uint64) uint64 {
    if a < b {
        return a
    }
    return b
}

func toBytesSlice(s []string) [][]byte {
	var ret [][]byte
	for _ , str := range s {
		ret = append(ret, []byte(str))
	}
	return ret
}

func main() {
	specifyFlags()
	flag.Parse()
	if len(flag.Args()) == 0 {
		uiPrintf2(0, "Missing file name")
		return
	}
	g_filename = flag.Args()[0]
	var err error
	g_file, err = os.OpenFile(g_filename, os.O_RDONLY, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer g_file.Close()
	err = checkFlags()
	if err != nil {
		fmt.Println(err)
		return
	}
	
	n := uint64(g_enduntil - g_startfrom)
	if g_numlines != 0 {
		n = min(g_linelen*g_numlines, n)
	}
	
	if g_doprint {
		uiPrintf1(1, "Dumping %s from %x to %x: %d Bytes (%s)\n", 
			g_filename, g_startfrom, uint64(g_startfrom)+n, 
			n, siValue(float64(n), "B"))
		ac := hexdumpaction{}
		iterateOverFile(uint64(g_startfrom), n, &ac)
		if err != nil {
			uiPrintf2(0, "%s\n", err)
			return
		}
	}
	
	if g_domd5 {
		uiPrintf1(1, "Calculating checksum of %s from %x to %x: %d Bytes (%s)\n", 
			g_filename, g_startfrom, uint64(g_startfrom)+n,
			n, siValue(float64(n), "B"))
		ac := md5action{}
		ac.Init()
		iterateOverFile(uint64(g_startfrom), n, &ac)
		if err != nil {
			uiPrintf2(0, "%s\n", err)
			return
		}
	}
	
	if len(g_patterns) != 0 {
		uiPrintf1(1, "Searching %s from %x to %x: %d Bytes (%s)\n",
			g_filename, g_startfrom, uint64(g_startfrom)+n,
			n, siValue(float64(n), "B"))
		ac := searchaction{}
		ac.Init(toBytesSlice(g_patterns))
		iterateOverFile(uint64(g_startfrom), n, &ac)
		if err != nil {
			uiPrintf2(0, "%s\n", err)
			return
		}
		
	}
}
