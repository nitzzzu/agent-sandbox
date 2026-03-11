package logutil

import (
	"context"
	"strings"

	"github.com/agent-sandbox/agent-sandbox/pkg/requestid"
	"k8s.io/klog/v2"
)

func prefix(ctx context.Context) string {
	rid := requestid.FromContext(ctx)
	if rid == "" {
		return ""
	}
	return "[rid=" + rid + "] "
}

// Infof logs with classic klog formatting (no quoted key/value pairs) and
// automatically prefixes request id (if present in ctx).
func Infof(ctx context.Context, v klog.Level, format string, args ...any) {
	p := prefix(ctx)
	if p != "" && !strings.HasPrefix(format, p) {
		format = p + format
	}
	klog.V(v).Infof(format, args...)
}

// Errorf logs with classic klog formatting (no quoted key/value pairs) and
// automatically prefixes request id (if present in ctx).
func Errorf(ctx context.Context, format string, args ...any) {
	p := prefix(ctx)
	if p != "" && !strings.HasPrefix(format, p) {
		format = p + format
	}
	klog.Errorf(format, args...)
}


