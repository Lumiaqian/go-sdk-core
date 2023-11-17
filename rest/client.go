package rest

import (
	"context"
	"io"
	"net/http"
)

type RequestPayload struct {
	Body       io.Reader            // request body
	FormFields map[string]string    // form fields
	FileFields map[string]FileField // file fields
}

type FileField struct {
	Filename string    // file name
	Content  io.Reader // file content
}

type HttpResponse struct {
	StatusCode  int
	Body        []byte
	Headers     map[string][]string
	ContentType string
}

type Client interface {
	DoRequest(ctx context.Context, method, url string, headers map[string]string, payload *RequestPayload) (*HttpResponse, error)
	Use(middleware Middleware)
}

type Middleware interface {
	Handle(ctx context.Context, req *http.Request, next MiddlewareHandler) (*http.Response, error)
}

type MiddlewareHandler func(ctx context.Context, req *http.Request) (*http.Response, error)
