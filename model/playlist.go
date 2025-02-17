package model

type Playlist struct {
	Id string
	Title string
	Subtitle string
	Description string
	Url string
	ImageUrl string
	SongCount uint16
	PlayCount uint32
	FollowerCount uint32
	FanCount uint32
	Songs []Song
}