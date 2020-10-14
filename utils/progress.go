package utils

import (
	"fmt"
	"io"

	"github.com/schollz/progressbar/v3"
	// "strings"
	"errors"
	"os"
)

type ProgressReader struct {
	bar *progressbar.ProgressBar
	io.Reader
}

func (p *ProgressReader) Wrap(r io.Reader, size int64) io.Reader {
	bar := progressbar.DefaultBytes(size)
	bar.RenderBlank() // show the progress at 0%
	return &ProgressReader{
		Reader: r,
		bar:    bar,
	}
}

func (p *ProgressReader) Read(b []byte) (int, error) {
	n, err := p.Reader.Read(b)
	p.bar.Add(n)
	return n, err
}

type ProgressWriter struct {
	bar    *progressbar.ProgressBar
	Writer io.WriteCloser
	tot    int
}

func (p *ProgressWriter) Wrap(w io.WriteCloser) io.Writer {
	p.Writer = w
	return p
}

func (p *ProgressWriter) Reset(size int64, filename string, payloadNumber int) {
	if len(filename) >= 20 {
		filename = fmt.Sprintf("...%s", filename[len(filename)-20:])
	}
	filename = fmt.Sprintf("%d: %s", payloadNumber, filename)
	p.bar = progressbar.DefaultBytes(size, filename)
}

func (p *ProgressWriter) Write(b []byte) (int, error) {
	n, err := p.Writer.Write(b)
	if err != nil {
		fmt.Println("ERRRR!!!")
		os.Exit(2)
	}
	if errors.Is(err, io.EOF) {
		fmt.Println("Finished!")
		os.Exit(1)
		p.bar.Finish()
		return n, err
	}
	p.bar.Add(n)
	return n, err
}
