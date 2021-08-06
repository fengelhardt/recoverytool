package main

import(
	"io"
// 	"fmt"
)

type action interface {
	Run(buf[] byte, abspos, tcnt, n uint64, lastbit bool) (uint64, error)
}

func iterateOverFile(from uint64, n uint64, ac action) error {
	var err error
	buf := make([]byte, g_bufsize)
	_, err = g_file.Seek(int64(from), 0)
	if err != nil {
		return err
	}
	offs := int(0)
	totalbytecnt := uint64(0)
	var l int
	// Read the file in chunks of g_buflen.
	// The action does not need to use all data and signals that by its return value.
	// Unused data is put to the beginning of the next chunk.
	for l, err = g_file.Read(buf); err == nil && totalbytecnt < n; l, err = g_file.Read(buf[offs:]) {
		l += offs
		// Respect n if we reached the last chunk.
		maxcnt := min(uint64(l), n - totalbytecnt)
		lastbit := (maxcnt != uint64(l))
// 		fmt.Printf("n=%d l=%d offs=%d  \n", n, l, offs)
		var bytesused uint64
		bytesused, err = ac.Run(buf[:maxcnt], from+totalbytecnt, totalbytecnt, n, lastbit)
		totalbytecnt += bytesused
		if err != nil {
			return err
		}
		// Put the remaining buffer content that did not fit into a complete line
		// to the beginning and load new data.
		carry := buf[bytesused:maxcnt]
		offs = len(carry)
		copy(buf[:offs], carry)
	}
	if err != io.EOF {
		return err
	}
	// Handle the last incomplete chunk.
	// The action should use every byte now.
	if totalbytecnt < n {
		ac.Run(buf[totalbytecnt:n-totalbytecnt], from+totalbytecnt, totalbytecnt, n, true)
	}
	return nil
}
