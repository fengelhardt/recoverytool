package main

import (
	"fmt"
	"time"
	"os"
)

type md5action struct {
	tp float64
	t1 time.Time
	a0, b0, c0, d0 uint32
}

func (a *md5action) Init() {
	a.t1 = time.Now()
	a.tp = 0.0
	a.a0 = 0x67452301
	a.b0 = 0xefcdab89
	a.c0 = 0x98badcfe
	a.d0 = 0x10325476
	return
}


var s = [64]uint32{
	7, 12, 17, 22,  7, 12, 17, 22,  7, 12, 17, 22,  7, 12, 17, 22,
	5,  9, 14, 20,  5,  9, 14, 20,  5,  9, 14, 20,  5,  9, 14, 20,
	4, 11, 16, 23,  4, 11, 16, 23,  4, 11, 16, 23,  4, 11, 16, 23,
	6, 10, 15, 21,  6, 10, 15, 21,  6, 10, 15, 21,  6, 10, 15, 21,
}

var K = [64]uint32{
	0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee,
	0xf57c0faf, 0x4787c62a, 0xa8304613, 0xfd469501,
	0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be,
	0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821,
	0xf61e2562, 0xc040b340, 0x265e5a51, 0xe9b6c7aa,
	0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
	0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed,
	0xa9e3e905, 0xfcefa3f8, 0x676f02d9, 0x8d2a4c8a,
	0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c,
	0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70,
	0x289b7ec6, 0xeaa127fa, 0xd4ef3085, 0x04881d05,
	0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
	0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039,
	0x655b59c3, 0x8f0ccc92, 0xffeff47d, 0x85845dd1,
	0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1,
	0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
}

func (a *md5action) iteratemd5(buf []byte) uint64 {
	bytesused := uint64(0)
	for ch := 0; ch < len(buf)/64 ; ch++ {
		// break chunk into sixteen 32-bit words M[i], 0 <= i <= 15
		var M [16]uint32
		for i := 0; i< len(M) ; i++ {
			M[i]  = uint32(buf[ch*64+i*4+0]) << 24
			M[i] |= uint32(buf[ch*64+i*4+1]) << 16
			M[i] |= uint32(buf[ch*64+i*4+2]) << 8
			M[i] |= uint32(buf[ch*64+i*4+3]) << 0
		}
		// TODO: buggy section begins here
		// Initialize hash value for this chunk:
		A := a.a0
		B := a.b0
		C := a.c0
		D := a.d0
		// Main loop:
		for i := uint32(0) ; i < uint32(64) ; i++ {
			var F, g uint32
			if i >= uint32(0) && i < uint32(16) {
				F = (B & C) | ((^B) & D)
				g = i
			}
			if i >= uint32(16) && i < uint32(32) {
				F = (D & B) | ((^D) & C)
				g = (uint32(5)*i + uint32(1)) % uint32(16)
			}
			if i >= uint32(32) && i < uint32(48) {
				F = B ^ C ^ D
				g = (uint32(3)*i + uint32(5)) % uint32(16)
			}
			if i >= uint32(48) && i < uint32(64) {
				F = C ^ (B ^ (^ D))
				g = (uint32(7)*i) % uint32(16)
			}
			// Be wary of the below definitions of a,b,c,d
			F = F + A + K[i] + M[g]  // M[g] must be a 32-bits block
			A = D
			D = C
			C = B
			B = B + leftrotate(F, s[i])
		}
		// Add this chunk's hash to result so far:
		a.a0 += A
		a.b0 += B
		a.c0 += C
		a.d0 += D
		// TODO: buggy section ends here
		// Test case: echo "a" > bla.txt ; ./recoverytool -md5 bla.txt
		// should yield 0cc175b9c0f1b6a831c399e269772661
		bytesused += 64
	}
	return bytesused
}

func (a *md5action) report() {
	md5 := [16]byte{
		byte((a.a0 >> 24) & 0xff),
		byte((a.a0 >> 16) & 0xff),
		byte((a.a0 >>  8) & 0xff),
		byte((a.a0 >>  0) & 0xff),
		byte((a.b0 >> 24) & 0xff),
		byte((a.b0 >> 16) & 0xff),
		byte((a.b0 >>  8) & 0xff),
		byte((a.b0 >>  0) & 0xff),
		byte((a.c0 >> 24) & 0xff),
		byte((a.c0 >> 16) & 0xff),
		byte((a.c0 >>  8) & 0xff),
		byte((a.c0 >>  0) & 0xff),
		byte((a.d0 >> 24) & 0xff),
		byte((a.d0 >> 16) & 0xff),
		byte((a.d0 >>  8) & 0xff),
		byte((a.d0 >>  0) & 0xff),
	}
	for _ , c := range md5 {
		fmt.Printf("%02x", c)
	}
	fmt.Println()
}

func (a *md5action) Run(buf[] byte, abspos, tcnt, n uint64, lastbit bool) (uint64, error) {
	if len(buf) < 0 {
		return uint64(0), nil
	}
	fmt.Println("buf", buf)
	bytesused := a.iteratemd5(buf)
	tcnt += bytesused
	t2 := time.Now()
	itp := float64(len(buf)) / float64((t2.Sub(a.t1)).Milliseconds()) * 1000.0
	a.tp = 0.9*a.tp + 0.1*itp
	eta := calcETA(float64(n-tcnt), a.tp)
	a.t1 = t2
	percentcomplete := float64(tcnt)/float64(n)*100.0
	fmt.Fprintf(os.Stderr, "%5.1f%%, %s ETA %s          \r", percentcomplete, siValue(a.tp, "B/s"), eta)
	if lastbit {
		// do the final md5 stuff
		bytesleft := uint64(len(buf))-bytesused
		buf2 := make([]byte, 0)
		buf2 = append(buf2, buf[bytesused:]...)
		cnt := bytesleft
		buf2 = append(buf2, 0x80)
		cnt++
		for ; cnt % 64 != 56 ; cnt++ {
			buf2 = append(buf2, 0x00)
		}
		bitlen := n*8
		for i := 7 ; i >= 0 ; i-- {
			b := byte((bitlen >> (i*8)) & 0xff)
			buf2 = append(buf2, b)
			cnt++
		}
		fmt.Println("buf2", buf2)
		bytesused = a.iteratemd5(buf2)
		if bytesused != uint64(len(buf2)) {
			panic("bytesused != uint64(len(buf2))")
		}
		tcnt += bytesleft
		if tcnt != n {
			panic("tcnt != n")
		}
		fmt.Fprintf(os.Stderr, "                                                           \r")
		a.report()
		return uint64(len(buf)), nil
	}
	return bytesused, nil
}

func leftrotate (x, c uint32) uint32 {
	return (x << c) | (x >> (uint32(32)-c))
}
