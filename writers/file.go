package writers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/boz/kail"
)

var _1kb = 1024
var _1mb = 1024 * _1kb
var MaxBufferLen = 32 * _1mb
var FlushBufferIntervalSeconds = 1000

func NewFileWriter(ctx context.Context) Writer {
	w := &writerFile{
		files:   make(map[string]*os.File),
		writers: make(map[string]*bytes.Buffer),
		ctx:     ctx,
	}

	w.startTimer()
	return w
}

type writerFile struct {
	ctx     context.Context
	files   map[string]*os.File
	writers map[string]*bytes.Buffer
	mu      sync.RWMutex
}

func (w *writerFile) Print(ev kail.Event) error {
	return w.Fprint(nil, ev)
}

func (w *writerFile) Fprint(out io.Writer, ev kail.Event) error {
	fileName := fmt.Sprintf("%s_%s.log", ev.Source().Name(), ev.Source().Container())

	writer := w.createOrGetWriter(fileName)
	writer.Write(ev.Log())
	// TODO: write line break by OS
	writer.Write([]byte("\n"))

	return nil
}

func (w *writerFile) startTimer() {
	go func() {
		defer func() {
			fmt.Printf("FileWriter timer stopped")
		}()
		for {
			if w.ctx.Err() == nil {
				for fileName, writer := range w.writers {
					w.flushWriter(fileName, writer)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (w *writerFile) flush(fileName string, writer *bytes.Buffer) {

	if writer.Len() >= MaxBufferLen {
		w.flushWriter(fileName, writer)

	}

}

func (w *writerFile) flushWriter(fileName string, writer *bytes.Buffer) {

	w.mu.Lock()
	defer w.mu.Unlock()

	file := w.createOrGetFile(fileName)

	_, err := file.Write(writer.Bytes())
	writer.Reset()
	if err != nil {
		panic(fmt.Sprintf("unable to write file %s: %s", fileName, err.Error()))
	}

}

func (w *writerFile) createOrGetFile(fileName string) *os.File {
	file, ok := w.files[fileName]
	if ok {
		return file
	}

	file, ok = w.files[fileName]
	if !ok {
		_file, err := os.Create(fileName)
		if err != nil {
			panic(fmt.Sprintf("unable to write file %s: %s", fileName, err.Error()))
		}
		file = _file
		w.files[fileName] = file
	}

	return file
}

func (w *writerFile) createOrGetWriter(fileName string) *bytes.Buffer {
	writer, ok := w.writers[fileName]
	if ok {
		return writer
	}

	writer, ok = w.writers[fileName]
	if !ok {
		writer = &bytes.Buffer{}
		w.writers[fileName] = writer
	}

	return writer
}
