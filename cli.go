package main

import (
	"flag"
	"fmt"
	"github.com/oyamo/substr/src"
	"os"
)

func main() {

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
