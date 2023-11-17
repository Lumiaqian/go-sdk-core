package rest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

type DefaultHttpClient struct {
	client      *http.Client
	middlewares []Middleware
}

func NewDefaultHttpClient() *DefaultHttpClient {
	return &DefaultHttpClient{
		client:      &http.Client{},
		middlewares: []Middleware{},
	}
}

func (c *DefaultHttpClient) Use(middleware Middleware) {
	c.middlewares = append(c.middlewares, middleware)
}

func (c *DefaultHttpClient) DoRequest(ctx context.Context, method, url string, headers map[string]string, payload *RequestPayload) (*HttpResponse, error) {
	if payload == nil {
		return nil, errors.New("payload cannot be nil")
	}
	// 创建一个新的 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, method, url, payload.Body)
	if err != nil {
		return nil, err
	}

	// 添加 headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 如果存在表单字段或文件字段，则处理 multipart/form-data
	if len(payload.FormFields) > 0 || len(payload.FileFields) > 0 {
		buffer := &bytes.Buffer{}
		writer := multipart.NewWriter(buffer)

		for field, value := range payload.FormFields {
			writer.WriteField(field, value)
		}

		for field, file := range payload.FileFields {
			fileWriter, err := writer.CreateFormFile(field, file.Filename)
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(fileWriter, file.Content)
			if err != nil {
				return nil, err
			}
		}

		writer.Close()
		req.Body = io.NopCloser(buffer)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	// 发送请求
	resp, err := c.executeMiddlewares(ctx, req, 0)
	if err != nil {
		return nil, err
	}

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	// 返回自定义的 HttpResponse 结构
	return &HttpResponse{
		StatusCode:  resp.StatusCode,
		Body:        body,
		Headers:     resp.Header,
		ContentType: resp.Header.Get("Content-Type"),
	}, nil
}

func (c *DefaultHttpClient) executeMiddlewares(ctx context.Context, req *http.Request, index int) (*http.Response, error) {
	if index >= len(c.middlewares) {
		return c.client.Do(req) // 如果没有其他中间件要执行，则直接发送请求
	}

	middleware := c.middlewares[index]
	return middleware.Handle(ctx, req, func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return c.executeMiddlewares(ctx, req, index+1) // 调用下一个中间件
	})
}
