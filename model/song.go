package model

type Song struct {
	Id string
	Title string
	Subtitle string
	Album Album
	Artists []Artist
	Year uint16
	Language string
	Url string
	ImageUrl string
	PlayCount uint32
	DurationMillis uint32
	HasLyrics bool
	Supports320Kbps bool
	EncryptedMediaUrl string
}