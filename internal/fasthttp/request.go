package fasthttp

import (
	"github.com/abhakash/gosaavn/internal/logging"
	"github.com/tidwall/gjson"
	"github.com/valyala/fasthttp"
	"time"
)

var DEFAULT_TIMEOUT_MILLIS = 60 * 1000

func GetRequest(url string, headers map[string]string) (gjson.Result, error) {

	// Create a request object
	req := fasthttp.AcquireRequest()

	// Create a response object
	resp := fasthttp.AcquireResponse()

	// Release the request and response objects back to the pool
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
	err := fasthttp.DoTimeout(req, resp, time.Duration(DEFAULT_TIMEOUT_MILLIS)*time.Millisecond)
	if err != nil {
		logging.Log.Error("Error while sending request: ", err)
		return gjson.Result{}, err
	}
	responseText, err := resp.BodyUncompressed()
	if err != nil {
		logging.Log.Error("Error while parsing response: ", err)
		return gjson.Result{}, err
	}

	// Print the response body (usually a JSON response)
	logging.Log.Debug("Response Code: ", resp.StatusCode())

	return gjson.ParseBytes(responseText), nil
}
