package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"math/rand"
)

func intHash(i int) uint64 {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(i))
	hashBytes := sha1.Sum(buf)
	ret := binary.BigEndian.Uint64(hashBytes[:])
	return ret
}

func fingerPrint(in []float64) uint64 {
	fpv := [64]float64{}
	for i, v := range in {
		h := intHash(i)
		for p := uint(0); p < 64; p++ {
			on := (h >> p) & 1
			if on > 0 {
				fpv[p] += v
			} else {
				fpv[p] -= v
			}
		}
	}
	out := uint64(0)

	for i := 0; i < 64; i++ {
		if fpv[i] > 0.0 {
			out = out | uint64(1<<uint(i))
		}
	}
	return out
}

func popcount(x uint64) (n byte) {
	// bit population count, see
	// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
	x -= (x >> 1) & 0x5555555555555555
	x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
	x += x >> 4
	x &= 0x0f0f0f0f0f0f0f0f
	x *= 0x0101010101010101
	return byte(x >> 56)
}

func gen(i int) []float64 {
	ret := []float64{}
	for ; i > 0; i-- {
		ret = append(ret, rand.Float64())
	}
	return ret
}

func main() {
	r := gen(10000)
	f := fingerPrint(r)
	fmt.Printf("%064b\n", f)
	for i := 0; i < 10; i++ {
		for j := 0; j < 100; j++ {
			r[rand.Int()%len(r)] = rand.Float64()
		}
		d := popcount(f ^ fingerPrint(r))
		f = fingerPrint(r)
		fmt.Printf("%064b: %d\n", f, d)
	}
}
