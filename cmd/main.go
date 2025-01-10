package main

import (
	"github.com/abhakash/gosaavn"
	"github.com/abhakash/gosaavn/internal/logging"
	"github.com/sirupsen/logrus"
)

func main(){
	logging.Init("app.log", logrus.DebugLevel)
	params := gosaavn.SearchSongsParams{Query: "eminem"}
	gosaavn.SearchSongs(params)
}