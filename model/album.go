package model

type Album struct {

	Id string
	Title string
	Subtitle string
	Songs []Song
	PrimaryArtist Artist
	Artists []Artist
	Url string
	ImageUrl string
	Year string
	Language string
	SongCount uint8
}