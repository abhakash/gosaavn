package fasthttp

import (
	"github.com/abhakash/gosaavn/internal/logging"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"time"
)

var DEFAULT_TIMEOUT_MILLIS = 60 * 1000
var DEFAULT_MAX_REDIRECTS = 3

func GetJson(url string, headers map[string]string) (gjson.Result, error) {

	response, err := get(url, headers)
	defer fasthttp.ReleaseResponse(response)
	if err != nil {
		logging.Log.Error("Error while sending request: ", err)
		return gjson.Result{}, err
	}
	responseBytes, err := response.BodyUncompressed()
	if err != nil {
		logging.Log.Error("Error while parsing response: ", err)
		return gjson.Result{}, err
	}

	// Print the response body (usually a JSON response)
	logging.Log.Debug("Response Code: ", response.StatusCode())

	return gjson.ParseBytes(responseBytes), nil
}

func GetBytes(url string, headers map[string]string) (int, []byte, error) {

	response, err := get(url, headers)
	defer fasthttp.ReleaseResponse(response)

	if err != nil {
		logging.Log.Error("Error while sending request: ", err)
		return -1, nil, err
	}
	responseBytes, err := response.BodyUncompressed()
	if err != nil {
		logging.Log.Error("Error while parsing response: ", err)
		return -1, nil, err
	}
	responseBody := append([]byte{}, responseBytes...)

	// Print the response body (usually a JSON response)
	logging.Log.Debug("Response Code: ", response.StatusCode())

	return response.StatusCode(), responseBody, nil
}

func GetBytesWithRedirects(url string, headers map[string]string) (int, []byte, error) {
	// Create a request object
	req := fasthttp.AcquireRequest()

	// Create a response object
	resp := fasthttp.AcquireResponse()

	// Release the request and response object back to the pool
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	// Set method and URL
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	for headerKey, headerVal := range headers {
		req.Header.Set(headerKey, headerVal)
	}
	logging.Log.Debug("Sending req to URL: ", url, "\n", req)

	// Perform the request
	err := fasthttp.DoRedirects(req, resp, DEFAULT_MAX_REDIRECTS)
	if err != nil {
		logging.Log.Error("Error while sending request: ", err)
		return -1, nil, err
	}
	responseBytes, err := resp.BodyUncompressed()
	if err != nil {
		logging.Log.Error("Error while parsing response: ", err)
		return -1, nil, err
	}
	responseBody := append([]byte{}, responseBytes...)

	return resp.StatusCode(), responseBody, nil
}

func get(url string, headers map[string]string) (*fasthttp.Response, error) {

	// Create a request object
	req := fasthttp.AcquireRequest()

	// Create a response object
	resp := fasthttp.AcquireResponse()

	// Only Release the request object back to the pool, since response object will be released by the caller
	defer fasthttp.ReleaseRequest(req)

	// Set method and URL
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")

	for headerKey, headerVal := range headers {
		req.Header.Set(headerKey, headerVal)
	}
	logging.Log.Debug("Sending req to URL: ", url, "\n", req)

	// Perform the request
	err := fasthttp.DoTimeout(req, resp, time.Duration(DEFAULT_TIMEOUT_MILLIS)*time.Millisecond)
	if err != nil {
		logging.Log.Error("Error while sending request: ", err)
		return nil, err
	}
	return resp, nil
}
