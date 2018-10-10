package grpcrw

import (
	"bytes"
	"fmt"
	"io"
)

type Sender interface {
	SendBytes([]byte) error
}

type WriterFunc func([]byte) error

type Writer struct {
	Sender
	buf bytes.Buffer
}

func NewWriter(w Sender) Writer {
	return Writer{
		Sender: w,
	}
}

func NewWriterF(f func([]byte) error) Writer {
	return NewWriter(WriterFunc(f))
}

func (f WriterFunc) SendBytes(p []byte) error {
	return f(p)
}

func (w Writer) Write(p []byte) (int, error) {
	// TODO(leeola): automatically chunk p if p > ChunkSize
	err := w.SendBytes(p)
	if err == io.EOF {
		return 0, io.EOF
	}
	if err != nil {
		return 0, fmt.Errorf("grpcrw sendbytes: %v", err)
	}

	return len(p), nil
}
