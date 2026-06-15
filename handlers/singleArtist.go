package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"stage/api"
	"stage/models"
	"stage/utils"
	"strconv"
	"strings"
)

type SingleHandler struct{
	Artists *api.ArtistCache
}

func NewSingleArtist(cache *api.ArtistCache) *SingleHandler{
	return &SingleHandler{
		Artists: cache,
	}
}


func (h *SingleHandler) SingleArtistHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	newPath := strings.Trim(path, "/")
	pathSlice := strings.Split(newPath, "/")

	artistsCache := h.Artists.GetAllArtists()

	if len(pathSlice) != 2 {
		utils.RenderError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	val := pathSlice[1]

	id, err := strconv.Atoi(val)
	if err != nil {
		utils.RenderError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	
	for _, artist := range artistsCache {
		if artist.Id == id {
			locJs, err := json.Marshal(artist.Locations)
			if err != nil{
				utils.RenderError(w, http.StatusInternalServerError, "JSON Marshal Error")
			}

			pageData := models.ArtistPageData{
				Artist:         artist.Artist,
				DatesLocations: artist.DatesLocations,
				Locations:      template.JS(string(locJs)),
			}

			err = utils.Tmpl.ExecuteTemplate(w, "artist.html", pageData)
			if err != nil {
				utils.RenderError(w, http.StatusInternalServerError, "Template Error")
				return
			}
			return
		}
	}

	utils.RenderError(w, http.StatusNotFound, "Artist not found")
}
