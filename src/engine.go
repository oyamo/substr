package src

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync/atomic"
)

type Engine struct {
	Input      *Input
	Cwd        string
	NumWorkers uint
}

type subProgress struct {
	Index uint
	Error error
}

func NewEngine() (engine *Engine, err error) {
	engine = &Engine{}

	engine.Cwd, err = os.Getwd()

	input, err := NewInput()
	if err != nil {
		return nil, err
	}

	// process the input
	var filesTmp = input.Files

	input.Files = make([]string, 0)

	for _, tmpFile := range filesTmp {
		filesScanned, err := filepath.Glob(tmpFile)
		if err != nil {
			return nil, err
		}
		input.Files = append(input.Files, filesScanned...)
	}

	if len(input.Files) > 1 && input.OuputFile != "" {
		return nil, errors.New("multiple inputs provided yet a single output file specified; use --output-dir instead")
	}

	engine.Input = input

	// Set default numworkers
	engine.NumWorkers = uint(runtime.NumCPU())

	return
}

func (e Engine) substituteChunk(start, stop uint, chanSubProgress chan subProgress, ctx *context.Context) {
	select {
	case <-(*ctx).Done():
		return
	default:
		re := regexp.MustCompile(e.Input.OriginalText)
		for i := range e.Input.Files[start:stop] {
			file, err := os.Open(e.Input.Files[i])
			if err != nil {
				chanSubProgress <- subProgress{Index: uint(i), Error: err}
				close(chanSubProgress)
				return
			}

			tmpFileName := fmt.Sprintf("%s/%s", os.TempDir(), uuid.New().String())
			if e.Input.OutputDir != "" {
				tmpFileName = fmt.Sprintf("%s/%s", e.Input.OutputDir, filepath.Base(e.Input.Files[i]))
			}
			if e.Input.OuputFile != "" && len(e.Input.Files) == 1 {
				tmpFileName = e.Input.OuputFile
			}
			// open file for output
			tmpFile, err := os.OpenFile(tmpFileName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
			if err != nil {
				chanSubProgress <- subProgress{Index: uint(i), Error: err}
				close(chanSubProgress)
				return
			}
			writer := bufio.NewWriter(tmpFile)
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				newStr := re.ReplaceAllLiteral(scanner.Bytes(), []byte(e.Input.SubstituteText))
				newStr = append(newStr, '\n')
				_, err := writer.Write(newStr)
				if err != nil {
					chanSubProgress <- subProgress{uint(i), err}
					close(chanSubProgress)
					return
				}
			}

			_ = writer.Flush()

			_ = file.Close()
			// copy the files
			dest, err := os.Create(e.Input.Files[i])
			if err != nil {
				chanSubProgress <- subProgress{Index: uint(i), Error: err}
				err = os.Remove(tmpFileName)
				close(chanSubProgress)
				return
			}
			_ = tmpFile.Close()
			source, err := os.Open(tmpFileName)
			if err != nil {
				chanSubProgress <- subProgress{Index: uint(i), Error: err}
				err = os.Remove(tmpFileName)
				close(chanSubProgress)
				return
			}

			if e.Input.OutputDir == "" && e.Input.OuputFile == "" {
				_, err = io.Copy(dest, source)
				if err != nil {
					chanSubProgress <- subProgress{Index: uint(i), Error: err}
					err = os.Remove(tmpFileName)
					close(chanSubProgress)
					return
				}
				err = os.Remove(tmpFileName)
			}

			chanSubProgress <- subProgress{uint(i), nil}

		}

	}
}

func (e Engine) substituteFiles() error {
	var progressChan = make(chan subProgress)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	var next atomic.Uint32
	for i := 0; uint(i) < e.NumWorkers; i += int(e.NumWorkers) {
		if uint(next.Load())+e.NumWorkers > uint(len(e.Input.Files)) {
			next.Store(uint32(len(e.Input.Files)))
		} else {
			next.Add(uint32(e.NumWorkers))
		}

		go e.substituteChunk(uint(i), uint(next.Load()), progressChan, &ctx)
	}

	var counter atomic.Uint32
	for pItem := range progressChan {
		if pItem.Error != nil {
			return pItem.Error
		}
		counter.Add(1)

		if counter.Load() == uint32(len(e.Input.Files)) {
			close(progressChan)
		}
	}

	return nil
}

func (e Engine) Run() (err error) {
	if e.Input.HasFlag(FLAG_VERSION) {
		Version()
		return
	}
	err = e.substituteFiles()
	return
}
