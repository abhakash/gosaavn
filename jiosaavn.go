package gosaavn

import (
	"github.com/abhakash/gosaavn/internal/fasthttp"
	"github.com/abhakash/gosaavn/internal/logging"
	"net/url"
	"strconv"
)

type SearchSongsParams struct {
	Query string
	Limit int
	Page  int
}

func SearchSongs(params SearchSongsParams) {
	if params.Limit == 0 {
		params.Limit = 20
	}

	if params.Page == 0 {
		params.Page = 1
	}

	searchParams := url.Values{}
	searchParams.Add("__call", "search.getResults")
	searchParams.Add("q", params.Query)
	searchParams.Add("p", strconv.Itoa(params.Page))
	searchParams.Add("n", strconv.Itoa(params.Limit))
	searchParams.Add("_format", "json")
	searchParams.Add("api_version", strconv.Itoa(saavnApiVersion))

	searchUrl, err := url.Parse(saavnApiUrl)
	if err != nil {
		logging.Log.Error("Failed parsing Saavn URL {}", saavnApiUrl, err)
	}
	searchUrl.RawQuery = searchParams.Encode()
	headers := map[string]string{
		"Accept":          "application/json, text/plain, */*",
		"Accept-Language": "en-US,en;q=0.9",
		"Connection":      "keep-alive",
	}

	searchResultJson, err := fasthttp.GetRequest(searchUrl.String(), headers)
	if err != nil {
		logging.Log.Error("Search FAILED")
	}
	logging.Log.Println(searchResultJson)
}
