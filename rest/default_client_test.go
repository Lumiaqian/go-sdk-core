package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the Middleware interface
type MockMiddleware struct {
	mock.Mock
}

func (m *MockMiddleware) Handle(ctx context.Context, req *http.Request, next MiddlewareHandler) (*http.Response, error) {
	args := m.Called(ctx, req, next)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestNewDefaultHttpClient(t *testing.T) {
	client := NewDefaultHttpClient()
	assert.NotNil(t, client, "NewDefaultHttpClient should create a new client instance")
	assert.NotNil(t, client.client, "NewDefaultHttpClient should have an http.Client")
	assert.Equal(t, 0, len(client.middlewares), "NewDefaultHttpClient should initialize with no middlewares")
}

func TestDefaultHttpClient_Use(t *testing.T) {
	client := NewDefaultHttpClient()
	middleware := new(MockMiddleware)
	client.Use(middleware)
	assert.Equal(t, 1, len(client.middlewares), "Middleware should be added to the middlewares slice")
}

func TestDefaultHttpClient_DoRequest_GET(t *testing.T) {
	// Setup a test HTTP server
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer testServer.Close()

	client := NewDefaultHttpClient()
	headers := map[string]string{"Content-Type": "application/json"}
	payload := &RequestPayload{ // Ensure payload is initialized
		Body:       nil,
		FormFields: map[string]string{},
		FileFields: map[string]FileField{},
	}
	resp, err := client.DoRequest(context.Background(), http.MethodGet, testServer.URL, headers, payload)
	assert.NoError(t, err, "DoRequest should not return an error on a GET request")
	assert.NotNil(t, resp, "Response should not be nil")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Response status code should be StatusOK")
	assert.Equal(t, "OK", string(resp.Body), "Response body should be 'OK'")
	_, err = client.DoRequest(context.Background(), http.MethodGet, testServer.URL, headers, nil)
	assert.Error(t, err, "DoRequest should return an error when payload is nil")
}

func TestDefaultHttpClient_DoRequest_POST(t *testing.T) {
	// Setup a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// POST
		assert.Equal(t, "POST", r.Method)

		// read body and check
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()
		assert.Equal(t, "post data", string(body))

		// send response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	}))
	defer server.Close()

	// init client
	client := NewDefaultHttpClient()

	// init payload
	payload := &RequestPayload{
		Body: strings.NewReader("post data"),
	}

	// send request
	response, err := client.DoRequest(context.Background(), "POST", server.URL, nil, payload)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "response", string(response.Body))
}

func TestDefaultHttpClient_DoRequest_PUT(t *testing.T) {
	// Setup a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)

		// read body and check
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()
		assert.Equal(t, "put data", string(body))

		// send response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("put response"))
	}))
	defer server.Close()

	client := NewDefaultHttpClient()
	payload := &RequestPayload{
		Body: strings.NewReader("put data"),
	}

	response, err := client.DoRequest(context.Background(), "PUT", server.URL, nil, payload)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "put response", string(response.Body))
}

func TestDefaultHttpClient_DoRequest_DELETE(t *testing.T) {
	// Setup a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)

		// send response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("delete response"))
	}))
	defer server.Close()
	// init client
	client := NewDefaultHttpClient()
	payload := &RequestPayload{
		Body: nil,
	}
	// send request
	response, err := client.DoRequest(context.Background(), "DELETE", server.URL, nil, payload)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "delete response", string(response.Body))
}

func TestDefaultHttpClient_DoRequest_PATCH(t *testing.T) {
	// Setup a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)

		// read body and check
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		defer r.Body.Close()
		assert.Equal(t, "patch data", string(body))

		// send response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("patch response"))
	}))
	defer server.Close()
	// init client
	client := NewDefaultHttpClient()
	payload := &RequestPayload{
		Body: strings.NewReader("patch data"),
	}
	// send request
	response, err := client.DoRequest(context.Background(), "PATCH", server.URL, nil, payload)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "patch response", string(response.Body))
}

func TestDefaultHttpClient_DoRequest_JSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var data map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&data)
		assert.NoError(t, err)
		assert.Equal(t, "value", data["key"])

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("json response"))
	}))
	defer server.Close()

	client := NewDefaultHttpClient()
	jsonBody, _ := json.Marshal(map[string]string{"key": "value"})
	payload := &RequestPayload{
		Body: bytes.NewReader(jsonBody),
	}
	heaer := map[string]string{
		"Content-Type": "application/json",
	}

	response, err := client.DoRequest(context.Background(), "POST", server.URL, heaer, payload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "json response", string(response.Body))
}

type SimpleXMLData struct {
	Key string `xml:"key"`
}

func TestDefaultHttpClient_DoRequest_XML(t *testing.T) {
	// Setup a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Confirm that the content type of the request header is XML
		assert.Equal(t, "application/xml", r.Header.Get("Content-Type"))

		// Decode the request body using an XML decoder
		var data SimpleXMLData
		err := xml.NewDecoder(r.Body).Decode(&data)
		assert.NoError(t, err)

		// Confirm that the request body contains the correct data
		assert.Equal(t, "value", data.Key)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("xml response"))
	}))
	defer server.Close()

	client := NewDefaultHttpClient()

	// Create an XML payload
	xmlData := SimpleXMLData{Key: "value"}
	xmlBody, err := xml.Marshal(xmlData)
	assert.NoError(t, err)

	payload := &RequestPayload{
		Body: bytes.NewReader(xmlBody),
	}
	header := map[string]string{
		"Content-Type": "application/xml",
	}
	// Send the request
	response, err := client.DoRequest(context.Background(), "POST", server.URL, header, payload)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	// Confirm that the response body contains the correct data
	assert.NoError(t, err)
	assert.Equal(t, "xml response", string(response.Body))
}

func TestDefaultHttpClient_DoRequest_FormData(t *testing.T) {
	// Setup a test HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "value", r.Form.Get("key"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("form response"))
	}))
	defer server.Close()

	client := NewDefaultHttpClient()
	formData := url.Values{}
	formData.Set("key", "value")
	payload := &RequestPayload{
		Body: strings.NewReader(formData.Encode()),
	}
	header := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	response, err := client.DoRequest(context.Background(), "POST", server.URL, header, payload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "form response", string(response.Body))
}

func TestDefaultHttpClient_DoRequest_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // sleep for 2 seconds
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewDefaultHttpClient()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.DoRequest(ctx, "GET", server.URL, nil, &RequestPayload{})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.DeadlineExceeded))
}

func TestDefaultHttpClient_DoRequest_Cancel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second) // sleep for 1 second
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewDefaultHttpClient()

	ctx, cancel := context.WithCancel(context.Background())

	// cancel now
	cancel()

	_, err := client.DoRequest(ctx, "GET", server.URL, nil, &RequestPayload{})
	assert.Error(t, err)
	// check error type
	assert.True(t, errors.Is(err, context.Canceled))
}
