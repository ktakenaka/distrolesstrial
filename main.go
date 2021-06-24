package main

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"math/rand"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"testing"
	"time"
	"unsafe"
)

const (
	cpuProfFile = "cpu.pprof"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 256)
		return &b
	},
}

var hashPool = sync.Pool{
	New: func() interface{} {
		return sha256.New()
	},
}

func foo(n int) string {
	bufptr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufptr)

	buf := *bufptr
	buf = buf[:0]

	h := hashPool.Get().(hash.Hash)
	defer hashPool.Put(h)
	h.Reset()

	x := strconv.AppendInt(buf, int64(n), 10)
	for i := 0; i < 100000; i++ {
		h.Write(x)
	}
	buf = buf[:0]
	sum := h.Sum(buf)

	b := make([]byte, 0, 256)
	for i := 0; i < int(sum[0]); i++ {
		x := sum[(i*7+1)%len(sum)] ^ sum[(i*5+3)%len(sum)]
		c := "abcdefghijklmnopqrstuvwxyz"[x%26]
		b = append(b, c)
	}
	sum = sum[:0]
	sum = append(sum, b...)
	return *(*string)(unsafe.Pointer(&sum))
}

func main() {
	cpufile, err := os.Create(cpuProfFile)
	if err != nil {
		panic(err)
	}
	defer cpufile.Close()

	err = pprof.StartCPUProfile(cpufile)
	if err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	if foo(12345) != "aajmtxaattdzsxnukawxwhmfotnm" {
		panic("wrong value!")
	}

	start := time.Now()
	fmt.Println("Allocs:", int(testing.AllocsPerRun(100, func() {
		foo(rand.Int())
	})))
	end := time.Now()
	fmt.Printf("time taken: %d", end.Sub(start))
}
