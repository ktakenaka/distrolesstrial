package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"strconv"
	"testing"
	"time"
	"unsafe"
)

const (
	cpuProfFile = "cpu.pprof"
)

func foo(n int) string {
	x := strconv.Itoa(n)
	buf := make([]byte, 0, 100000*len(x))
	for i := 0; i < 100000; i++ {
		buf = append(buf, x...)
	}
	sum := sha256.Sum256(buf)

	b := make([]byte, 0, int(sum[0]))
	for i := 0; i < int(sum[0]); i++ {
		x := sum[(i*7+1)%len(sum)] ^ sum[(i*5+3)%len(sum)]
		c := "abcdefghijklmnopqrstuvwxyz"[x%26]
		b = append(b, c)
	}
	return *(*string)(unsafe.Pointer(&b))
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
