package butch

import (
	"bytes"
	"strings"
)

type workerByteBuffer struct {
	full   *bytes.Buffer
	chunk  *bytes.Buffer
	target *string
}

func (w *workerByteBuffer) Write(p []byte) (int, error) {
	// Writing
	w.full.Write(p)
	w.chunk.Write(p)

	idx := bytes.IndexRune(p, '\n')
	if idx >= 0 {
		full := string(w.chunk.Bytes())
		idx := strings.Index(full, "\n")
		line := strings.TrimSpace(full[0:idx])
		if len(line) > 0 {
			*w.target = line
		}
		w.chunk = bytes.NewBufferString(full[idx+1:])
	}

	return len(p), nil
}

func (w *workerByteBuffer) done() {
	if w.chunk.Len() > 0 {
		*w.target = strings.TrimSpace(string(w.chunk.Bytes()))
	}
}
