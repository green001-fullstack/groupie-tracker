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
	geoCache *GeoCache
}

func NewArtistCache(geoCache *GeoCache) *ArtistCache{
	return &ArtistCache{
		geoCache: geoCache,
	}
}
