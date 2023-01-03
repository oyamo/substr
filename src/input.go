package src

import (
	"flag"
	"fmt"
	"os"
)

const (
	FLAG_HELP    = 1 << iota
	FLAG_D       // -d, --output-dir
	FLAG_V       // -v, --verbose
	FLAG_VERSION // -version
	FLAG_O       // -o, --output-file
	FLAG_T       // -t, --original-text
	FLAG_S       // -s, --substitute-text

)

type Input struct {
	Files          []string
	Flags          []uint8 // tiny for less mem
	OuputFile      string
	OriginalText   string
	SubstituteText string
	OutputDir      string
}

type InputParser interface {
	HasFlag(uint8) bool
}

// NewInput creates a new Input struct and parses the command line arguments.
// It returns the Input struct and an error if there is one.
func NewInput() (*Input, error) {
	args := os.Args[1:]

	input := &Input{
		Files: make([]string, 0),
		Flags: make([]uint8, 0),
	}

	// set flag usage
	flag.Usage = Usage

	var help = flag.Bool("help", false, "")
	var version = flag.Bool("version", false, "")
	var verbose = flag.Bool("verbose", false, "")
	var verboseShort = flag.Bool("v", false, "")
	var outputFile = flag.String("output-file", "", "")
	var outputFileShort = flag.String("o", "", "")
	var outputDir = flag.String("output-dir", "", "")
	var outputDirShort = flag.String("d", "", "")
	var originalText = flag.String("original-text", "", "")
	var originalTextShort = flag.String("t", "", "")
	var substituteText = flag.String("substitute-text", "", "")
	var substituteTextShort = flag.String("s", "", "")

	flag.Parse()

	if len(args) == 0 {
		return nil, fmt.Errorf("no arguments")
	}

	if *help {
		input.Flags = append(input.Flags, FLAG_HELP)
		return input, nil
	}

	if *version {
		input.Flags = append(input.Flags, FLAG_VERSION)
		return input, nil
	}

	if *verbose || *verboseShort {
		input.Flags = append(input.Flags, FLAG_V)
	}

	if *outputFile != "" || *outputFileShort != "" {
		input.OuputFile = *outputFile
	}

	if *outputDir != "" || *outputDirShort != "" {
		input.OutputDir = *outputDir
		if *outputDir == "" {
			input.OutputDir = *outputDirShort
		}
	}

	if *originalText != "" || *originalTextShort != "" {
		input.OriginalText = *originalText
		if input.OriginalText == "" {
			input.OriginalText = *originalTextShort
		}
	} else {
		return nil, fmt.Errorf("no original text")
	}

	if *substituteText != "" || *substituteTextShort != "" {
		input.SubstituteText = *substituteText
		if *substituteText == "" {
			input.SubstituteText = *substituteTextShort
		}
	} else {
		return nil, fmt.Errorf("no substitute text")
	}

	input.Files = flag.Args()

	return input, nil
}

// HasFlag returns true if the input has the flag.
func (i *Input) HasFlag(flag uint8) bool {
	for _, f := range i.Flags {
		if f == flag {
			return true
		}
	}
	return false
}
