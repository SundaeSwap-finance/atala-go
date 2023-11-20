package atala

import (
	"log"
	"net/http"
	"net/url"
)

const BASE_PATH = "prism-agent/"

// Create a new Client, pass the full url for an atala PRISM agent instance
func CreateClient(baseUrl string) *Client {
	httpClient := &http.Client{}
	c := NewClient(httpClient)
	var urlErr error
	c.BaseURL, urlErr = url.Parse(baseUrl)
	if urlErr != nil {
		log.Panicf("ERROR parsing BaseURL: %v", urlErr)
	}
	return c
}

// API Health Check
func (c *Client) SystemHealth() (*Health, *ApiError, int, error) {
	resp, health, apiErr, err := GetRequest[Health](c, BASE_PATH+"_system/health")
	return health, apiErr, resp.StatusCode, err
}
