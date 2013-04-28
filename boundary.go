package boundary

import (
	"bytes"
	"io"
)

// Writer allows to split passed writes on arbitraty boundaries.
// e.g. Write([]byte("foo\nbar\baz")) actually become THREE writes
// to the underlying writer.
type Writer struct {
	w            io.Writer
	Boundary     []byte // Boundary we split the passed write on
	SkipBoundary bool   // Should we ignore (true) or pass on (false) the Boundary?
}

func NewWriter(w io.Writer, boundary string, skip bool) *Writer {
	return &Writer{w: w, Boundary: []byte(boundary), SkipBoundary: skip}
}

func (bw *Writer) Write(s []byte) (n int, err error) {
	if len(s) == 0 || len(bw.Boundary) == 0 {
		return bw.w.Write(s)
	}
	sep := bw.Boundary
	save := len(sep)
	if bw.SkipBoundary {
		save = 0
	}
	c := sep[0]
	start := 0
	na := 0
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i] == c && (len(sep) == 1 || bytes.Equal(s[i:i+len(sep)], sep)) {
			na, err = bw.w.Write(s[start : i+save])
			n += na
			if err != nil {
				return n, err
			}
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	na, err = bw.w.Write(s[start:])
	n += na
	return n, err
}
