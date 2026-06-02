package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// HTTPMethod defines a custom type for HTTP methods.
type HTTPMethod string

// Constants for valid HTTP methods.
const (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	PUT    HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
	// PATCH  HTTPMethod = "PATCH"
)

// RequestParams holds parameters for the HttpRequest function.
type RequestParams struct {
	Method      HTTPMethod
	URL         string
	Params      map[string]string
	Body        interface{}
	Headers     map[string]string
	ContentType string
}

// HttpRequest performs an HTTP request with the specified method and URL.
// It optionally takes query parameters, a request body, headers, and content type,
// returning the response body and an error, if any.
func HttpRequest[T any](params RequestParams) (T, error) {
	var result T

	// Validate the HTTP method
	if !isValidMethod(params.Method) {
		return result, fmt.Errorf("invalid HTTP method")
	}

	// Prepare the URL with query parameters if the method is GET
	reqURL, err := prepareURL(params.URL, params.Method, params.Params)
	if err != nil {
		return result, fmt.Errorf("preparing URL: %w", err)
	}

	// Prepare the request body
	bodyReader, err := prepareRequestBody(params.Body, params.Method)
	if err != nil {
		return result, fmt.Errorf("preparing request body: %w", err)
	}

	// Create the request
	req, err := http.NewRequest(string(params.Method), reqURL, bodyReader)
	if err != nil {
		return result, fmt.Errorf("creating request: %w", err)
	}

	// Set headers
	setHeaders(req, params.Headers, params.ContentType)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	// Handle the response
	if err := handleResponse(resp, &result); err != nil {
		return result, err
	}

	return result, nil
}

// isValidMethod checks if the provided HTTP method is valid.
func isValidMethod(method HTTPMethod) bool {
	switch method {
	case GET, POST, PUT, DELETE:
		return true
	}
	return false
}

// prepareURL prepares the URL with query parameters if the method is GET.
func prepareURL(baseURL string, method HTTPMethod, params map[string]string) (string, error) {
	if method != GET || params == nil {
		return baseURL, nil
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("parsing URL: %w", err)
	}

	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()
	return u.String(), nil
}

// prepareRequestBody prepares the request body by marshalling it into JSON if necessary.
func prepareRequestBody(body interface{}, method HTTPMethod) (io.Reader, error) {
	if body == nil || method == GET {
		return nil, nil
	}

	switch b := body.(type) {
	case string:
		return bytes.NewBufferString(b), nil
	case []byte:
		return bytes.NewBuffer(b), nil
	default:
		bodyBytes, err := json.Marshal(b)
		if err != nil {
			return nil, fmt.Errorf("marshalling request body: %w", err)
		}
		return bytes.NewBuffer(bodyBytes), nil
	}
}

// setHeaders sets the headers for the HTTP request.
func setHeaders(req *http.Request, headers map[string]string, contentType string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if contentType == "" {
		contentType = "application/json"
	}
	req.Header.Set("Content-Type", contentType)
}

// handleResponse reads and unmarshals the response body, checking the HTTP response status code.
func handleResponse[T any](resp *http.Response, result *T) error {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	_ = json.Unmarshal(responseBody, result)

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("received non-success status code: %d %s [%+v]", resp.StatusCode, resp.Status, result)
	}

	return nil
}
