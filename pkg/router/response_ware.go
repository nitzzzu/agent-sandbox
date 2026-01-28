package router

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
)

var (
	_ http.Flusher        = (*ResponseWare)(nil)
	_ http.ResponseWriter = (*ResponseWare)(nil)
)

// ResponseWare is an implementation of http.ResponseWriter and http.Flusher
// that captures the response code and size.
type ResponseWare struct {
	ResponseCode int
	ResponseSize int

	writer      http.ResponseWriter
	wroteHeader bool
	// hijacked is whether this connection has been hijacked
	// by a Handler with the Hijacker interface.
	// This is guarded by a mutex in the default implementation.
	hijacked atomic.Bool
}

// NewResponseWare creates an http.ResponseWriter that captures the response code and size.
func NewResponseWare(w http.ResponseWriter, responseCode int) *ResponseWare {
	return &ResponseWare{
		writer:       w,
		ResponseCode: responseCode,
	}
}

// Unwrap returns the underlying writer
func (rr *ResponseWare) Unwrap() http.ResponseWriter {
	return rr.writer
}

// Flush flushes the buffer to the client.
func (rr *ResponseWare) Flush() {
	rr.writer.(http.Flusher).Flush()
}

// Hijack calls Hijack() on the wrapped http.ResponseWriter if it implements
// http.Hijacker interface, which is required for net/http/httputil/reverseproxy
// to handle connection upgrade/switching protocol.  Otherwise returns an error.
func (rr *ResponseWare) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	c, rw, err := HijackIfPossible(rr.writer)
	if err != nil {
		rr.hijacked.Store(true)
	}
	return c, rw, err
}

// Header returns the header map that will be sent by WriteHeader.
func (rr *ResponseWare) Header() http.Header {
	return rr.writer.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
func (rr *ResponseWare) Write(p []byte) (int, error) {
	rr.ResponseSize += len(p)
	return rr.writer.Write(p)
}

// WriteHeader sends an HTTP response header with the provided status code.
func (rr *ResponseWare) WriteHeader(code int) {
	if rr.wroteHeader || rr.hijacked.Load() {
		return
	}

	rr.writer.WriteHeader(code)
	rr.wroteHeader = true
	rr.ResponseCode = code
}

// HijackIfPossible calls Hijack() on the given http.ResponseWriter if it implements
// http.Hijacker interface, which is required for net/http/httputil/reverseproxy
// to handle connection upgrade/switching protocol.  Otherwise returns an error.
func HijackIfPossible(w http.ResponseWriter) (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("wrapped writer of type %T can't be hijacked", w)
	}
	return hj.Hijack()
}
