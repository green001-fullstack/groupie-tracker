package api

import (
	"stage/models"
	"sync"
)

type geoResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type ArtistCache struct{
	artistsCache []models.FullArtist
	mu sync.RWMutex
	artistByID map[int]models.FullArtist
	geoCache *GeoCache
}

func NewArtistCache(geoCache *GeoCache) *ArtistCache{
	return &ArtistCache{
		geoCache: geoCache,
		artistByID: make(map[int]models.FullArtist),
	}
}
