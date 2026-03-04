package requestid

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// HeaderName is the canonical HTTP header name used across this repo.
// We also accept common variants for compatibility.
const HeaderName = "X-Request-Id"

const headerNameAlt = "X-Request-ID"

type ctxKey struct{}

// New returns a best-effort unique request id.
// (We intentionally keep it dependency-free; it only needs to be unique enough for logs.)
func New() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// FromHeader reads the request id from headers (canonical first, then common variant).
func FromHeader(h http.Header) string {
	if h == nil {
		return ""
	}
	if v := h.Get(HeaderName); v != "" {
		return v
	}
	return h.Get(headerNameAlt)
}

// SetHeader sets the canonical header name.
func SetHeader(h http.Header, id string) {
	if h == nil || id == "" {
		return
	}
	h.Set(HeaderName, id)
}

// EnsureHeader ensures the canonical header exists and returns the id.
func EnsureHeader(h http.Header) string {
	id := FromHeader(h)
	if id == "" {
		id = New()
	}
	SetHeader(h, id)
	return id
}

// With stores request id into context.
func With(ctx context.Context, id string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if id == "" {
		return ctx
	}
	return context.WithValue(ctx, ctxKey{}, id)
}

// FromContext extracts request id from context.
func FromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	v := ctx.Value(ctxKey{})
	id, _ := v.(string)
	return id
}


