package main

import (
	"github.com/abhakash/gosaavn/api"
	"github.com/abhakash/gosaavn/model"
	"github.com/abhakash/gosaavn/internal/logging"
	"github.com/sirupsen/logrus"
)

func main() {
	logging.Init("app.log", logrus.DebugLevel)
	params := model.SearchParams{Query: "eminem"}
	api.SearchSongs(params)

	// APT song
	song, _ := api.GetSong("uNWdki50")
	logging.Log.Info(song)
	api.DownloadSong(song, "downloads")
}
