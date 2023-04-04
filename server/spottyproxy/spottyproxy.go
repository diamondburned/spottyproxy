package spottyproxy

import librespot "github.com/librespot-org/librespot-golang/Spotify"

// Timestamp is a timestamp in unix epoch seconds.
type Timestamp int64

type ImageSizeType uint8

const (
	ImageSizeDefault = librespot.Image_DEFAULT
	ImageSizeSmall   = librespot.Image_SMALL
	ImageSizeLarge   = librespot.Image_LARGE
	ImageSizeXLarge  = librespot.Image_XLARGE
)

type Image struct {
	FileID   []byte        `json:"fileID"`
	SizeType ImageSizeType `json:"sizeType"`
	Width    uint32        `json:"width"`
	Height   uint32        `json:"height"`
}

type Playlist struct {
	URI    string  `json:"uri"`
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

type Track struct {
	GID        []byte   `json:"gid"`
	Name       string   `json:"name"`
	Album      Album    `json:"album"`
	Artists    []Artist `json:"artists"`
	Duration   uint32   `json:"duration"`
	Popularity int32    `json:"popularity"`
}

type AlbumType uint8

const (
	AlbumTypeAlbum       = librespot.Album_ALBUM
	AlbumTypeSingle      = librespot.Album_SINGLE
	AlbumTypeCompilation = librespot.Album_COMPILATION
	AlbumTypeEP          = librespot.Album_EP
)

type Album struct {
	GID        []byte    `json:"gid"`
	Name       string    `json:"name"`
	Type       AlbumType `json:"type"`
	Date       Timestamp `json:"date"`
	Artists    []Artist  `json:"artists"`
	Covers     []Image   `json:"covers"`
	Popularity int32     `json:"popularity"`
}

type Artist struct {
	GID          []byte             `json:"gid"`
	Name         string             `json:"name"`
	TopTracks    map[string][]Track `json:"topTracks"`
	Albums       []Album            `json:"albums"`
	Singles      []Album            `json:"singles"`
	Compilations []Album            `json:"compilations"`
	AppearsOn    []Album            `json:"appearsOn"`
	Portraits    []Image            `json:"portraits"`
	Related      []Artist           `json:"related"`
}
