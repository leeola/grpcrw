package grpcrw

import "io"

const defaultChunkSize = 1024

// ChunkRead implements limiting of a Reader []byte slice
// size to MaxChunkSize.
//
// In basic testing, gRPC seems to struggle
type ChunkRead struct {
	io.Reader
	MaxChunkSize int
}

// Copy the reader to the writer, wrapping the reader with a
// ChunkReader automatically.
//
// See ChunkRead for further docs.
func Copy(w io.Writer, r io.Reader) (int64, error) {
	return io.Copy(w, ChunkReader(r))
}

// ChunkReader returns an io.Reader compatible chunker
// which ensures that reads are limited to MaxChunkSize
// per read.
//
// See ChunkRead for further docs.
func ChunkReader(r io.Reader) ChunkRead {
	return ChunkRead{
		Reader:       r,
		MaxChunkSize: defaultChunkSize,
	}
}

func (m ChunkRead) Read(p []byte) (int, error) {
	if len(p) > m.MaxChunkSize {
		p = p[:m.MaxChunkSize]
	}

	return m.Reader.Read(p)
}
