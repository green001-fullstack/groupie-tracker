package handlers

import (
	// "strings"
	// "strconv"
	// "stage/api"
	// "stage/models"
)

// type Handler struct{
// 	Cache *api.ArtistCache
// }

// func (h *Handler) HandleSearch(query string)([]models.FullArtist, error) {

// 	artistsCache := h.Cache.GetAllArtists()
	
// 	query = strings.TrimSpace(query)
	
// 	var result []models.FullArtist
// 	if query != "" {
// 	query = strings.ToLower(query)

// 	for _, artist := range artistsCache {
// 		matched := false
// 		if strings.Contains(strings.ToLower(artist.Name), query) {
// 			matched = true
// 		}

// 		if !matched {
// 			for _, member := range artist.Members {
// 				if strings.Contains(strings.ToLower(member), query) {
// 					matched = true
// 					break
// 				}
// 			}
// 		}

// 		if !matched {
// 			for location := range artist.DatesLocations {
// 				if strings.Contains(strings.ToLower(location), query) {
// 					matched = true
// 					break
// 				}
// 			}
// 		}

// 		if !matched {
// 			creationDate := strconv.Itoa(artist.CreationDate)
// 			if strings.Contains(creationDate, query) {
// 				matched = true
// 				// break
// 			}
// 		}

// 		if !matched {
// 			if strings.Contains(strings.ToLower(artist.FirstAlbum), query) {
// 				matched = true
// 				// break
// 			}
// 		}
// 		if matched {
// 			result = append(result, artist)
// 		}
// 	}
// }else{
// 	result = artistsCache
// }
// 	return result, nil
// }