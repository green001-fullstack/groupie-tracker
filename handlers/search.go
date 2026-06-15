package handlers

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"stage/api"
	"stage/models"
	"stage/utils"
)

type Handler struct{
	Cache *api.ArtistCacheb
}


func HandleSearch(w http.ResponseWriter, r *http.Request)([]models.FullArtist,error) {
	artistsCache := api.GetFullArtist()
	if r.Method != http.MethodGet {
		utils.RenderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return []models.FullArtist{}, fmt.Errorf("method not allowed")
	}

	query := r.URL.Query().Get("query")
	query = strings.TrimSpace(query)
	
	var result []models.FullArtist
	if query != "" {
	query = strings.ToLower(query)

	for _, artist := range artistsCache {
		matched := false
		if strings.Contains(strings.ToLower(artist.Name), query) {
			matched = true
		}

		if !matched {
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), query) {
					matched = true
					break
				}
			}
		}

		if !matched {
			for location := range artist.DatesLocations {
				if strings.Contains(strings.ToLower(location), query) {
					matched = true
					break
				}
			}
		}

		if !matched {
			creationDate := strconv.Itoa(artist.CreationDate)
			if strings.Contains(creationDate, query) {
				matched = true
				break
			}
		}

		if !matched {
			if strings.Contains(artist.FirstAlbum, query) {
				matched = true
				break
			}
		}
		if matched {
			result = append(result, artist)
		}
	}
}else{
	result = artistsCache
}



	return result, nil
}