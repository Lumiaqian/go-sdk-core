package rest

import (
	"context"
	"io"
	"net/http"
)

type RequestPayload struct {
	Body       io.Reader            // 用于普通的请求体，如 JSON 或 XML
	FormFields map[string]string    // 普通的表单字段
	FileFields map[string]FileField // 文件字段
}

type FileField struct {
	Filename string    // 文件名
	Content  io.Reader // 文件内容
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
