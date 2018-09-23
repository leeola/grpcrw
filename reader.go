package grpcrw

import (
	"bytes"
	"fmt"
	"io"
)

type Recver interface {
	RecvBytes() ([]byte, error)
}

type RecverFunc func() ([]byte, error)

type Reader struct {
	Recver
	buf bytes.Buffer
}

func NewReader(r Recver) Reader {
	return Reader{
		Recver: r,
	}
}

func NewReaderF(f func() ([]byte, error)) Reader {
	return NewReader(RecverFunc(f))
}

func (r Reader) Read(p []byte) (int, error) {
	if r.Recver != nil {
		b, err := r.RecvBytes()
		if err == io.EOF {
			return 0, err // no wrap EOF
		}
		if err != nil {
			return 0, fmt.Errorf("grpcrw recvbytes: %v", err)
		}

		// TODO(leeola): determine how we'll implicitly determine EOF.
		// 0 bytes seems a decent starting point. Though, we may want to
		// support 0, since Go readers do iirc.
		isEOF := len(b) == 0

		if !isEOF {
			// write recv'd bytes to buf, as it may be more than we can read into
			// p
			if _, err := r.buf.Write(b); err != nil {
				return 0, fmt.Errorf("grpcrw buf write: %v", err)
			}
		} else {
			// "close" this Recver so the Read will finish due to the Recver being EOF
			r.Recver = nil
		}
	}

	// read whatever length requested into p from the buf
	return r.buf.Read(p)
}

func (f RecverFunc) RecvBytes() ([]byte, error) {
	return f()
}
