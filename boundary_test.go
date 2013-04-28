package boundary

import (
	"bytes"
	"testing"
)

// Simple writer backed by a slice of slices of bytes
type AppendWriter [][]byte

func (a *AppendWriter) Write(p []byte) (int, error) {
	*a = AppendWriter(append([][]byte(*a), p))
	return len(p), nil
}

func TestWriter(t *testing.T) {
	a := AppendWriter{}
	w := NewWriter(&a, "\n", false)
	sample := []byte("foo\nbar\nbaz")
	want := bytes.SplitAfter(sample, []byte("\n"))
	w.Write(sample)
	if len(a) != len(want) {
		t.Errorf("got %d writes, want %d", len(a), len(want))
	}

	for i, v := range want {
		if !bytes.Equal(v, a[i]) {
			t.Errorf("%d: got '%s', want '%s'", a[i], v)
		}
	}
}
