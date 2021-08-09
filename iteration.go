package main

import(
	"io"
	"fmt"
)

type action interface {
	Run(buf[] byte, abspos, tcnt, n uint64, lastbit bool) (uint64, error)
}

type readevent struct {ibuf, l int; err error}

func readFileIntoAlternatingBuffers(buf [2][]byte, evt chan<- readevent,
	from, n uint64, stop *bool) {
	var err error
	totalbytecnt := uint64(0)
	ibuf := 0
	_, err = g_file.Seek(int64(from), 0)
	if err != nil {
		evt <- readevent{0, 0, err}
		return
	}
	var l int
	for err == nil && totalbytecnt < n && !(*stop) {
		l, err = g_file.Read(buf[ibuf])
		if l > int(n-totalbytecnt) {
			l = int(n-totalbytecnt)
			err = io.EOF
		}
		totalbytecnt += uint64(l)
		evt <- readevent{ibuf, l, err}
		ibuf = (ibuf+1) % 2
	}
	close(evt)
	return
}

func iterateOverFile(from, n uint64, ac action) error {
	var err error
	extrasize := int(g_bufsize * 100) / 10
	bufsize := int(g_bufsize) + extrasize
	buf := [2][]byte{
		make([]byte, bufsize),
		make([]byte, bufsize),
	}
	evt := make(chan readevent)
	stop := false
	// Double buffered reading in a separate goroutine.
	go readFileIntoAlternatingBuffers(
		[2][]byte{
			buf[0][extrasize:], 
			buf[1][extrasize:],
		},
		evt,
		from, n, &stop)
	extraspaceused := int(0)
	totalbytecnt := uint64(0)
	// Handle the bytes in chunks of g_buflen.
	for e := range evt {
// 		fmt.Printf("n=%d l=%d offs=%d  \n", n, l, offs)
		ibuf, l, err := e.ibuf, e.l, e.err
		if err != nil && err != io.EOF {
			stop = true
			return err
		}
		lastbit := (err == io.EOF)
		ibufNext := (ibuf+1) % 2
		var bytesused uint64
		runbuf := buf[ibuf][extrasize-extraspaceused:extrasize+l]
		bytesused, err = ac.Run(runbuf, from+totalbytecnt, totalbytecnt, n, lastbit)
		if bytesused > uint64(l) {
			panic(fmt.Errorf("bytesused > l"))
		}
		totalbytecnt += bytesused
		// Run() might not have used all bytes.
		// Copy the left over bytes to the extra space at the begining of the other buffer
		carry := buf[ibuf][extrasize+int(bytesused):extrasize+l]
		extraspaceused = len(carry)
		if extraspaceused > extrasize {
			stop = true
			return fmt.Errorf("Not enough extra buffer space: Please use a larger buffer size (-b flag) or shorten your search pattern(s).")
		}
		copy(buf[ibufNext][extrasize-extraspaceused:extrasize], carry)
	}
	stop = true
	// In case the file ended before we have read n bytes.
	if err != io.EOF {
		return err
	}
	return nil
}
