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

func F(x, y, z uint32) uint32 {
	return (x & y) | ((^x) & z)
}

func G(x, y, z uint32) uint32 {
	return (x & z) | (y & (^z))
}

func H(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func I(x, y, z uint32) uint32 {
	return y ^ (x | (^z))
}

func FF(a, b, c, d, x, s, ac uint32) uint32 {
	a += F(b, c, d) + x + ac
	a = leftrotate(a, s)
	a += b
	return a
}
func GG(a, b, c, d, x, s, ac uint32) uint32 {
	a += G(b, c, d) + x + ac
	a = leftrotate((a), (s))
	a += b
	return a
}

func HH(a, b, c, d, x, s, ac uint32) uint32 {
	a += H (b, c, d) + x + ac
	a = leftrotate (a, s)
	a += b
	return a
}

func II(a, b, c, d, x, s, ac uint32) uint32 {
	a += I(b, c, d) + x + ac
	a = leftrotate (a, s)
	a += b
	return a
}

func (ac *md5action) iteratemd5(buf []byte) uint64 {
	bytesused := uint64(0)
	for ch := 0; ch < len(buf)/64 ; ch++ {
		// break chunk into sixteen 32-bit words x[i], 0 <= i <= 15
		var x [16]uint32
		for i := 0; i< len(x) ; i++ {
			x[i]  = uint32(buf[ch*64+i*4+0]) << 0
			x[i] |= uint32(buf[ch*64+i*4+1]) << 8
			x[i] |= uint32(buf[ch*64+i*4+2]) << 16
			x[i] |= uint32(buf[ch*64+i*4+3]) << 24
		}
		a := ac.a0
		b := ac.b0
		c := ac.c0
		d := ac.d0
		
		/* Round 1 */
		a = FF (a, b, c, d, x[ 0],  7, 0xd76aa478); /* 1 */
		d = FF (d, a, b, c, x[ 1], 12, 0xe8c7b756); /* 2 */
		c = FF (c, d, a, b, x[ 2], 17, 0x242070db); /* 3 */
		b = FF (b, c, d, a, x[ 3], 22, 0xc1bdceee); /* 4 */
		a = FF (a, b, c, d, x[ 4],  7, 0xf57c0faf); /* 5 */
		d = FF (d, a, b, c, x[ 5], 12, 0x4787c62a); /* 6 */
		c = FF (c, d, a, b, x[ 6], 17, 0xa8304613); /* 7 */
		b = FF (b, c, d, a, x[ 7], 22, 0xfd469501); /* 8 */
		a = FF (a, b, c, d, x[ 8],  7, 0x698098d8); /* 9 */
		d = FF (d, a, b, c, x[ 9], 12, 0x8b44f7af); /* 10 */
		c = FF (c, d, a, b, x[10], 17, 0xffff5bb1); /* 11 */
		b = FF (b, c, d, a, x[11], 22, 0x895cd7be); /* 12 */
		a = FF (a, b, c, d, x[12],  7, 0x6b901122); /* 13 */
		d = FF (d, a, b, c, x[13], 12, 0xfd987193); /* 14 */
		c = FF (c, d, a, b, x[14], 17, 0xa679438e); /* 15 */
		b = FF (b, c, d, a, x[15], 22, 0x49b40821); /* 16 */

		/* Round 2 */
		a = GG (a, b, c, d, x[ 1],  5, 0xf61e2562); /* 17 */
		d = GG (d, a, b, c, x[ 6],  9, 0xc040b340); /* 18 */
		c = GG (c, d, a, b, x[11], 14, 0x265e5a51); /* 19 */
		b = GG (b, c, d, a, x[ 0], 20, 0xe9b6c7aa); /* 20 */
		a = GG (a, b, c, d, x[ 5],  5, 0xd62f105d); /* 21 */
		d = GG (d, a, b, c, x[10],  9,  0x2441453); /* 22 */
		c = GG (c, d, a, b, x[15], 14, 0xd8a1e681); /* 23 */
		b = GG (b, c, d, a, x[ 4], 20, 0xe7d3fbc8); /* 24 */
		a = GG (a, b, c, d, x[ 9],  5, 0x21e1cde6); /* 25 */
		d = GG (d, a, b, c, x[14],  9, 0xc33707d6); /* 26 */
		c = GG (c, d, a, b, x[ 3], 14, 0xf4d50d87); /* 27 */
		b = GG (b, c, d, a, x[ 8], 20, 0x455a14ed); /* 28 */
		a = GG (a, b, c, d, x[13],  5, 0xa9e3e905); /* 29 */
		d = GG (d, a, b, c, x[ 2],  9, 0xfcefa3f8); /* 30 */
		c = GG (c, d, a, b, x[ 7], 14, 0x676f02d9); /* 31 */
		b = GG (b, c, d, a, x[12], 20, 0x8d2a4c8a); /* 32 */

		/* Round 3 */
		a = HH (a, b, c, d, x[ 5],  4, 0xfffa3942); /* 33 */
		d = HH (d, a, b, c, x[ 8], 11, 0x8771f681); /* 34 */
		c = HH (c, d, a, b, x[11], 16, 0x6d9d6122); /* 35 */
		b = HH (b, c, d, a, x[14], 23, 0xfde5380c); /* 36 */
		a = HH (a, b, c, d, x[ 1],  4, 0xa4beea44); /* 37 */
		d = HH (d, a, b, c, x[ 4], 11, 0x4bdecfa9); /* 38 */
		c = HH (c, d, a, b, x[ 7], 16, 0xf6bb4b60); /* 39 */
		b = HH (b, c, d, a, x[10], 23, 0xbebfbc70); /* 40 */
		a = HH (a, b, c, d, x[13],  4, 0x289b7ec6); /* 41 */
		d = HH (d, a, b, c, x[ 0], 11, 0xeaa127fa); /* 42 */
		c = HH (c, d, a, b, x[ 3], 16, 0xd4ef3085); /* 43 */
		b = HH (b, c, d, a, x[ 6], 23,  0x4881d05); /* 44 */
		a = HH (a, b, c, d, x[ 9],  4, 0xd9d4d039); /* 45 */
		d = HH (d, a, b, c, x[12], 11, 0xe6db99e5); /* 46 */
		c = HH (c, d, a, b, x[15], 16, 0x1fa27cf8); /* 47 */
		b = HH (b, c, d, a, x[ 2], 23, 0xc4ac5665); /* 48 */

		/* Round 4 */
		a = II (a, b, c, d, x[ 0],  6, 0xf4292244); /* 49 */
		d = II (d, a, b, c, x[ 7], 10, 0x432aff97); /* 50 */
		c = II (c, d, a, b, x[14], 15, 0xab9423a7); /* 51 */
		b = II (b, c, d, a, x[ 5], 21, 0xfc93a039); /* 52 */
		a = II (a, b, c, d, x[12],  6, 0x655b59c3); /* 53 */
		d = II (d, a, b, c, x[ 3], 10, 0x8f0ccc92); /* 54 */
		c = II (c, d, a, b, x[10], 15, 0xffeff47d); /* 55 */
		b = II (b, c, d, a, x[ 1], 21, 0x85845dd1); /* 56 */
		a = II (a, b, c, d, x[ 8],  6, 0x6fa87e4f); /* 57 */
		d = II (d, a, b, c, x[15], 10, 0xfe2ce6e0); /* 58 */
		c = II (c, d, a, b, x[ 6], 15, 0xa3014314); /* 59 */
		b = II (b, c, d, a, x[13], 21, 0x4e0811a1); /* 60 */
		a = II (a, b, c, d, x[ 4],  6, 0xf7537e82); /* 61 */
		d = II (d, a, b, c, x[11], 10, 0xbd3af235); /* 62 */
		c = II (c, d, a, b, x[ 2], 15, 0x2ad7d2bb); /* 63 */
		b = II (b, c, d, a, x[ 9], 21, 0xeb86d391); /* 64 */
		
		// Add this chunk's hash to result so far:
		ac.a0 += a
		ac.b0 += b
		ac.c0 += c
		ac.d0 += d
		bytesused += 64
	}
	return bytesused
}

func (a *md5action) report() {
	md5 := [16]byte{
		byte((a.a0 >>  0) & 0xff),
		byte((a.a0 >>  8) & 0xff),
		byte((a.a0 >> 16) & 0xff),
		byte((a.a0 >> 24) & 0xff),
		byte((a.b0 >>  0) & 0xff),
		byte((a.b0 >>  8) & 0xff),
		byte((a.b0 >> 16) & 0xff),
		byte((a.b0 >> 24) & 0xff),
		byte((a.c0 >>  0) & 0xff),
		byte((a.c0 >>  8) & 0xff),
		byte((a.c0 >> 16) & 0xff),
		byte((a.c0 >> 24) & 0xff),
		byte((a.d0 >>  0) & 0xff),
		byte((a.d0 >>  8) & 0xff),
		byte((a.d0 >> 16) & 0xff),
		byte((a.d0 >> 24) & 0xff),
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
	//fmt.Println("buf", buf)
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
		for i := 0 ; i < 8 ; i++ {
			b := byte((bitlen >> (i*8)) & 0xff)
			buf2 = append(buf2, b)
			cnt++
		}
		//fmt.Println("buf2", buf2)
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
