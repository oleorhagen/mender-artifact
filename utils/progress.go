package utils

import (
	// "errors"
	"fmt"
	"io"
	"os"

	"github.com/oleorhagen/progress/progressbar"
)

type ProgressReader struct {
	bar *progressbar.Bar
	io.Reader
}

func (p *ProgressReader) Wrap(r io.Reader, size int64) io.Reader {
	bar := progressbar.New()
	bar.Size = size
	// bar.RenderBlank() // show the progress at 0% Nice to have for later
	return &ProgressReader{
		Reader: r,
		bar:    bar,
	}
}

func (p *ProgressReader) Read(b []byte) (int, error) {
	n, err := p.Reader.Read(b)
	p.bar.Tick(int64(n))
	return n, err
}

type ProgressWriter struct {
	bar    *progressbar.Bar
	Writer io.WriteCloser
	tot    int
}

func (p *ProgressWriter) Wrap(w io.WriteCloser) io.Writer {
	fmt.Fprintln(os.Stderr)
	p.Writer = w
	return p
}

func (p *ProgressWriter) Reset(size int64, filename string, payloadNumber int) {
	if len(filename) >= 20 {
		filename = fmt.Sprintf("...%s", filename[len(filename)-20:])
	}
	filename = fmt.Sprintf("%d: %s", payloadNumber, filename)
	p.bar = progressbar.New()
	p.bar.Size = size
}

func (p *ProgressWriter) Finish() {
	if p.bar != nil {
		p.bar.Finish()
	}
}

func (p *ProgressWriter) Write(b []byte) (int, error) {
	n, err := p.Writer.Write(b)
	p.bar.Tick(int64(n))
	return n, err
}
