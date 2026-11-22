package registration

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
)

func newRequestID() string {
	id := make([]byte, 8)
	if _, err := rand.Read(id); err != nil {
		return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(id)
}

// RequestIDMiddleware adds a request ID to the response and request headers.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := newRequestID()
		w.Header().Set("X-Request-ID", requestID)
		r.Header.Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

// RequestLogger creates a logger that includes the request ID in each message.
func RequestLogger(r *http.Request) *log.Logger {
	return log.New(log.Writer(), fmt.Sprintf("[request_id=%s] ", r.Header.Get("X-Request-ID")), log.Flags())
}
