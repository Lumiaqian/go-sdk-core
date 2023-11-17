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

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, method, url, payload.Body)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Handle multipart/form-data if there are form fields or file fields
	if len(payload.FormFields) > 0 || len(payload.FileFields) > 0 {
		buffer := &bytes.Buffer{}
		writer := multipart.NewWriter(buffer)

		// Write form fields
		for field, value := range payload.FormFields {
			writer.WriteField(field, value)
		}

		// Write file fields
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

		// Close the writer
		err = writer.Close()
		if err != nil {
			return nil, err
		}

		req.Body = io.NopCloser(buffer)
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	// Send the request
	resp, err := c.executeMiddlewares(ctx, req, 0)
	if err != nil {
		return nil, err
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	// Return a custom HttpResponse structure
	return &HttpResponse{
		StatusCode:  resp.StatusCode,
		Body:        body,
		Headers:     resp.Header,
		ContentType: resp.Header.Get("Content-Type"),
	}, nil
}

func (c *DefaultHttpClient) executeMiddlewares(ctx context.Context, req *http.Request, index int) (*http.Response, error) {
	if index >= len(c.middlewares) {
		return c.client.Do(req) // If there are no other middlewares to execute, send the request directly
	}

	// Execute the current middleware and pass in a callback function to call the next middleware
	return c.middlewares[index].Handle(ctx, req, func(ctx context.Context, req *http.Request) (*http.Response, error) {
		return c.executeMiddlewares(ctx, req, index+1)
	})
}
