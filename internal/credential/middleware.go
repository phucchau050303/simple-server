package credential

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type requestLoggerKey struct{}

var requestCounter uint64

func newRequestID() string {
	id := make([]byte, 8)
	if _, err := rand.Read(id); err != nil {
		return fmt.Sprintf("fallback-%d", atomic.AddUint64(&requestCounter, 1))
	}
	return hex.EncodeToString(id)
}

func RequestLogger(ctx context.Context) *log.Logger {
	requestID, _ := ctx.Value(requestLoggerKey{}).(string)
	return log.New(log.Writer(), fmt.Sprintf("[request_id=%s] ", requestID), log.Flags())
}

func RequestIDInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	requestID := newRequestID()
	ctx = context.WithValue(ctx, requestLoggerKey{}, requestID)
	if err := grpc.SetHeader(ctx, metadata.Pairs("x-request-id", requestID)); err != nil {
		return nil, err
	}
	return handler(ctx, req)
}
