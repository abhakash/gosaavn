package api

import (
	"fmt"
	"github.com/abhakash/gosaavn/internal/fasthttp"
	"github.com/abhakash/gosaavn/internal/logging"
	"github.com/abhakash/gosaavn/internal/util"
	"github.com/abhakash/gosaavn/model"
	"github.com/tidwall/gjson"
	"net/url"

	"strconv"
)

var headers = map[string]string{
	"Accept":          "application/json, text/plain, */*",
	"Accept-Language": "en-US,en;q=0.9",
	"Connection":      "keep-alive",
}

func SearchSongs(params model.SearchParams) ([]model.Song, error) {
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
		logging.Log.Errorf("Failed parsing Saavn URL %s : \n %+v", saavnApiUrl, err)
		return nil, err
	}
	searchUrl.RawQuery = searchParams.Encode()

	searchResultJson, err := fasthttp.GetJson(searchUrl.String(), headers)
	if err != nil {
		logging.Log.Errorf("API request for search songs %s FAILED \n %+v", params.Query, err)
		return nil, err
	}

	var songs []model.Song = []model.Song{}

	if !searchResultJson.Get("results").Exists() || len(searchResultJson.Get("results").Array()) == 0 {
		logging.Log.Errorf("No result found while searching songs for query: %s", params.Query)
		return nil, fmt.Errorf("no song results for query %s", params.Query)
	}

	for _, song := range searchResultJson.Get("results").Array() {
		songs = append(songs, *parseSongJson(song))
	}
	return songs, nil
}

func GetSong(id string) (*model.Song, error) {
	searchParams := url.Values{}
	searchParams.Add("__call", "song.getDetails")
	searchParams.Add("pids", id)
	searchParams.Add("_format", "json")
	searchParams.Add("api_version", strconv.Itoa(saavnApiVersion))

	songDetailsUrl, err := url.Parse(saavnApiUrl)
	if err != nil {
		logging.Log.Error("Failed parsing Saavn URL {}", saavnApiUrl, err)
		return nil, err
	}

	songDetailsUrl.RawQuery = searchParams.Encode()
	songResultJson, err := fasthttp.GetJson(songDetailsUrl.String(), headers)

	if err != nil {
		logging.Log.Errorf("API request for retrieving song details for id %s FAILED: \n %+v", id, err)
		return nil, err
	}
	if !songResultJson.Get(id).Exists() {
		logging.Log.Errorf("No result found while querying song id %s", id)
		return nil, fmt.Errorf("no song found for id %s", id)
	}
	songResult := songResultJson.Get(id)
	return parseSongJson(songResult), nil

}

func DownloadSong(song *model.Song, locationDirectory string) error {

	err := util.CreateDirectory(locationDirectory)
	if err != nil {
		logging.Log.Errorf("Could not create directory at %s", locationDirectory)
		return err
	}

	authMediaUrl, mediaType, err := generateAuthUrlForSong(song)
	if err != nil {
		return err
	}
	decodedAuthMediaUrl, err := url.QueryUnescape(authMediaUrl)
	logging.Log.Info("Sending request to" + decodedAuthMediaUrl + " \n")
	if err != nil {
		logging.Log.Errorf("Error decoding Auth media URL string: %s", err)
		return err
	}
	statusCode, songResponse, err := fasthttp.GetBytesWithRedirects(decodedAuthMediaUrl, headers)
	if err != nil {
		logging.Log.Errorf("API request for downloading song for id %s FAILED: \n %+v", song.Id, err)
		return err
	}
	if statusCode != 200 {
		logging.Log.Errorf("Failed to download file for songId [%s]. Auth URL : [%s]. \n StatusCode : [%d]",
			song.Id, decodedAuthMediaUrl, statusCode)
		return err
	}

	songFileName := fmt.Sprintf("%s_%s.%s", song.Title, song.Subtitle, mediaType)

	file, err := util.CreateFile(locationDirectory, songFileName)
	if err != nil {
		logging.Log.Errorf("Error creating file: %s", err)
		return err
	}
	defer file.Close()

	// Write the downloaded content to the file
	_, err = file.Write(songResponse)
	if err != nil {
		logging.Log.Errorf("Error writing to file: %s", err)
		return err
	}
	return nil
}

func generateAuthUrlForSong(song *model.Song) (string, string, error) {
	encrypted_media_url := url.QueryEscape(song.EncryptedMediaUrl)

	searchParams := url.Values{}
	searchParams.Add("__call", "song.generateAuthToken")
	searchParams.Add("url", encrypted_media_url)
	searchParams.Add("bitrate", "320")
	searchParams.Add("_format", "json")
	searchParams.Add("api_version", strconv.Itoa(saavnApiVersion))

	authUrl, err := url.Parse(saavnApiUrl)
	if err != nil {
		logging.Log.Error("Failed parsing Saavn URL {}", saavnApiUrl, err)
		return "", "", err
	}

	authUrl.RawQuery = searchParams.Encode()
	authResultJson, err := fasthttp.GetJson(authUrl.String(), map[string]string{})
	if err != nil {
		logging.Log.Errorf("API request for retrieving song auth url for id %s FAILED: \n %+v", song.Id, err)
		return "", "", err
	}
	if authResultJson.Get("status").Exists() && authResultJson.Get("status").String() != "success" {
		logging.Log.Errorf("No Auth URL found for song id %s", song.Id)
		return "", "", fmt.Errorf("could not retrieve Auth URL for song id %s", song.Id)
	}
	if !authResultJson.Get("auth_url").Exists() {
		logging.Log.Errorf("No Auth URL found for song id %s", song.Id)
		return "", "", fmt.Errorf("no Auth URL found for song id %s", song.Id)
	}
	mediaType := authResultJson.Get("type").String()
	return authResultJson.Get("auth_url").String(), mediaType, nil
}

func parseSongJson(songJson gjson.Result) *model.Song {

	song := new(model.Song)
	song.Id = songJson.Get("id").String()
	song.Title = songJson.Get("title").String()
	song.Subtitle = songJson.Get("subtitle").String()
	song.Url = songJson.Get("perma_url").String()
	song.ImageUrl = songJson.Get("image").String()
	if playCount, err := util.StringToUInt(songJson.Get("play_count").String(), 32); err == nil {
		song.PlayCount = playCount.(uint32)
	}

	if year, err := util.StringToUInt(songJson.Get("year").String(), 16); err == nil {
		song.Year = year.(uint16)
	}

	song.Language = songJson.Get("language").String()
	song.HasLyrics = songJson.Get("has_lyrics").Bool()
	if songJson.Get("more_info").Exists() {
		moreInfo := songJson.Get("more_info")
		song.Supports320Kbps = moreInfo.Get("320kbps").Bool()
		song.EncryptedMediaUrl = moreInfo.Get("encrypted_media_url").String()
		if duration, err := util.StringToUInt(moreInfo.Get("duration").String(), 32); err == nil {
			song.DurationMillis = duration.(uint32)
		}
		if moreInfo.Get("artistMap").Exists() {
			artists := []model.Artist{}
			for _, artistJson := range moreInfo.Get("artistMap").Get("artists").Array() {
				artists = append(artists, parseArtistJson(artistJson))
			}
			song.Artists = artists
		}
		album := new(model.Album)
		album.Id = moreInfo.Get("album_id").String()
		album.Title = moreInfo.Get("album").String()
		album.Url = moreInfo.Get("album_url").String()
		song.Album = *album
	}
	return song
}

func parseArtistJson(artistJson gjson.Result) model.Artist {
	artist := new(model.Artist)
	artist.Id = artistJson.Get("id").String()
	artist.Name = artistJson.Get("name").String()
	artist.Url = artistJson.Get("perma_url").String()
	artist.ImageUrl = artistJson.Get("image").String()
	return *artist
}
