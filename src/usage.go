package src

import (
	"fmt"
	"os"
	"runtime"
)

const VERSION = 1.0

func Usage() {
	fmt.Fprintf(os.Stderr, `substr
A commandline utility for replacing text in files

Argument            Shortened	Description
--------            ---------  --------------------------------
 --help              -h         display help and exit		
 --version                      display version and exit
 --original-text     -t         text to replace
 --subsitute-text    -s         final text
 --output-file       -o         file path to redirect the output
 --output-dir        -d         dir path to redirect output from
                                all files

Example: (substr -t foo -s bar example.txt) replaces all 
occurences of foo with bar
`)
}

func Version() {
	fmt.Fprintf(os.Stdout, "substr: %.1f, golang: %s\n", VERSION, runtime.Version())
}
