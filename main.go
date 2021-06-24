package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	cpuProfFile = "cpu.pprof"
)

func foo(n int) string {
	var buf bytes.Buffer
	for i := 0; i < 100000; i++ {
		buf.WriteString(strconv.Itoa(n))
	}
	sum := sha256.Sum256(buf.Bytes())

	var b []byte
	for i := 0; i < int(sum[0]); i++ {
		x := sum[(i*7+1)%len(sum)] ^ sum[(i*5+3)%len(sum)]
		c := strings.Repeat("abcdefghijklmnopqrstuvwxyz", 10)[x]
		b = append(b, c)
	}
	return string(b)
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
