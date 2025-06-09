package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/rahadianir/dealls/internal/config"
)

var reqID uint64

type TracerMiddleware struct {
}

func (tracer *TracerMiddleware) Tracer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// setup request id
		requestID := r.Header.Get("request-id")
		if requestID == "" {
			prefix := uuid.NewString()
			myid := atomic.AddUint64(&reqID, 1)
			requestID = fmt.Sprintf("%s-%06d", prefix, myid)
		}
		ctx = context.WithValue(ctx, config.RequestIDKey, requestID)

		// setup real IP
		ip := getRealIP(r)
		if ip == "" {
			ip = r.RemoteAddr
		}
		ctx = context.WithValue(ctx, config.IPKey, ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

var trueClientIP = http.CanonicalHeaderKey("True-Client-IP")
var xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
var xRealIP = http.CanonicalHeaderKey("X-Real-IP")

func getRealIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}
