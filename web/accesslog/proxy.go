package accesslog

import (
	"net/http"
)

var _ http.Flusher = &flushWriter{}

func WrapWriter(w http.ResponseWriter) WriterProxy {
	bw := basicWriter{ResponseWriter: w}
	if _, ok := w.(http.Flusher); ok {
		return &flushWriter{bw}
	}
	return &bw
}

type WriterProxy interface {
	http.ResponseWriter
	Status() int
	BytesWritten() int
	Unwrap() http.ResponseWriter
}

type basicWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       int
}

func (b *basicWriter) WriteHeader(code int) {
	if !b.wroteHeader {
		b.code = code
		b.wroteHeader = true
		b.ResponseWriter.WriteHeader(code)
	}
}

func (b *basicWriter) Write(buf []byte) (int, error) {
	b.WriteHeader(http.StatusOK)
	n, err := b.ResponseWriter.Write(buf)
	b.bytes += n
	return n, err
}

func (b *basicWriter) Status() int {
	return b.code
}

func (b *basicWriter) BytesWritten() int {
	return b.bytes
}

func (b *basicWriter) Unwrap() http.ResponseWriter {
	return b.ResponseWriter
}

type flushWriter struct {
	basicWriter
}

func (f *flushWriter) Flush() {
	fl := f.basicWriter.ResponseWriter.(http.Flusher)
	fl.Flush()
}
