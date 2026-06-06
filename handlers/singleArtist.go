package handlers


import (
	"net/http"
    "stage/utils"
	"stage/models"
	"stage/api"
	"strings"
	"strconv"
	"encoding/json"
	"html/template"
)


func SingleArtistHandler(w http.ResponseWriter, r *http.Request) {
    artistsCache := api.GetFullArtist()

	path := r.URL.Path
	newPath := strings.Trim(path, "/")
	pathSlice := strings.Split(newPath, "/")

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
