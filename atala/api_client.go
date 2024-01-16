package atala

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

// Create a new Client, pass an http.Client or nil to create a new one
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{httpClient: httpClient}
	return c
}

// Convenience wrappers to make GET, POST, PUT & PATCH requests
func GetRequest[O any](c *Client, path string, successCodes ...int) (*http.Response, *O, *ApiError, error) {
	return makeRequest[O](c, "GET", path, nil, successCodes...)
}
func PostRequest[O any](c *Client, path string, body []byte, successCodes ...int) (*http.Response, *O, *ApiError, error) {
	return makeRequest[O](c, "POST", path, body, successCodes...)
}
func PutRequest[O any](c *Client, path string, body []byte, successCodes ...int) (*http.Response, *O, *ApiError, error) {
	return makeRequest[O](c, "PUT", path, body, successCodes...)
}
func PatchRequest[O any](c *Client, path string, body []byte, successCodes ...int) (*http.Response, *O, *ApiError, error) {
	return makeRequest[O](c, "PATCH", path, body, successCodes...)
}

// Create and excute an HTTP request
func makeRequest[O any](c *Client, method string, path string, body []byte, successCodes ...int) (*http.Response, *O, *ApiError, error) {
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, nil, nil, err
	}
	var obj O
	var apiErr ApiError
	resp, err := c.do(req, &obj, &apiErr, successCodes...)
	if apiErr.Status > 0 {
		return resp, nil, &apiErr, err
	}
	return resp, &obj, nil, err
}

// Create an HTTP request
// *TODO - decide how we want to pass/parse the body object
func (c *Client) newRequest(method, path string, body []byte) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = bytes.NewBuffer(body)
		// buf = new(bytes.Buffer)
		// err := json.NewEncoder(buf).Encode(body)
		// if err != nil {
		// 	return nil, err
		// }
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	// req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

// Execute an HTTP request
func (c *Client) do(req *http.Request, v interface{}, apiErr *ApiError, successCodes ...int) (*http.Response, error) {
	// Execute request
	resp, err := c.httpClient.Do(req)
	req.Close = true
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Determine which status codes indicate success, in addition to 200 & 201
	var apiSuccessResponse bool
	for _, c := range append(successCodes, 200, 201) {
		if c == resp.StatusCode {
			apiSuccessResponse = true
		}
	}
	if apiSuccessResponse {
		// Decode success response JSON
		err = json.NewDecoder(resp.Body).Decode(v)
	} else if resp.StatusCode == 404 {
		// Respond to 404 with an APIError
		*apiErr = ApiError{Status: resp.StatusCode}
	} else {
		// Decode failure response APIError
		err = json.NewDecoder(resp.Body).Decode(&apiErr)
	}
	return resp, err
}
