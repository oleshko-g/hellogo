package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	_ "strings"

	"oleshko-g/benchstring/align"
)

const (
	alignMe           string = "ultra-left"
	maxLen            int    = 15
	leftSymbol        rune   = ' '
	cpuProfileDefault string = "cpu.out"
	memProfileDefault string = "mem.out"
)

var cpuprofile = flag.String("cpuprofile", cpuProfileDefault, "write cpu profile to `file`")
var memprofile = flag.String("memprofile", memProfileDefault, "write memory profile to `file`")

func main() {
	flag.Parse()
	fCPU, err := os.Create(*cpuprofile)
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer fCPU.Close() // error handling omitted for example

	if err := pprof.StartCPUProfile(fCPU); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	align.AlignRight(alignMe, maxLen, leftSymbol)

	fMem, err := os.Create(*memprofile)
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	defer fMem.Close() // error handling omitted for example
	runtime.GC()       // get up-to-date statistics
	// Lookup("allocs") creates a profile similar to go test -memprofile.
	// Alternatively, use Lookup("heap") for a profile
	// that has inuse_space as the default index.
	if err := pprof.Lookup("allocs").WriteTo(fMem, 0); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
}
