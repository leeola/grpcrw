package grpcrw

import "io"

const defaultMessageSize = 1024

type MessageRead struct {
	io.Reader
	MessageSize int
}

func MessageReader(r io.Reader) MessageRead {
	return MessageRead{
		Reader:      r,
		MessageSize: defaultMessageSize,
	}
}

func (m MessageRead) Read(p []byte) (int, error) {
	if len(p) > 1024 {
		p = p[:1024]
	}

	return m.Reader.Read(p)
}
