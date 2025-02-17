package model

type Artist struct {
	Id string
	Name string
	Bio string
	DOB string
	Url string
	ImageUrl string
	Wiki string
	FanCount uint32
	FollowerCount uint32
	TopSongs []Song
	TopAlbums []Album
	SimilarArtists []Artist
	IsVerified bool
}
