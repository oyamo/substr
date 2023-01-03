package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"substr/src"
)

func main() {

	f, err := os.Create("cpu.profile")
	if err != nil {
		return
	}

	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	engine, err := src.NewEngine()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	err = engine.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
