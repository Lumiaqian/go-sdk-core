package middleware

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/Lumiaqian/go-sdk-core/log"
	"github.com/Lumiaqian/go-sdk-core/rest"
)

type LogMiddleware struct {
	logger log.Logger
}

func NewLogMiddleware(logger log.Logger) *LogMiddleware {
	return &LogMiddleware{
		logger: logger,
	}
}

func (m *LogMiddleware) Handle(ctx context.Context, req *http.Request, next rest.MiddlewareHandler) (*http.Response, error) {
	start := time.Now()

	// log before request
	m.logger.Log(log.INFO, "request_start", "", "method", req.Method, "url", req.URL.String())

	// call next
	resp, err := next(ctx, req)

	if resp != nil {
		// copy response body
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()                                    // close response body
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // re-open response body

		// log response
		m.logger.Log(log.INFO, "response_body", bodyBytes)
	}
	// log after request
	m.logger.Log(log.INFO, "request_end", "method", req.Method, "url", req.URL.String(), "duration", time.Since(start).String())

	return resp, err
}
