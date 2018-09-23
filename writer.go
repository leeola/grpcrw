package grpcrw

import (
	"bytes"
	"fmt"
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
	if err := w.SendBytes(p); err != nil {
		return 0, fmt.Errorf("grpcrw sendbytes: %v", err)
	}

	return len(p), nil
}
