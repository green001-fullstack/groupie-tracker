package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"stage/api"
	"stage/models"
	"strings"
	"strconv"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html", "templates/artist.html"))

var artistsCache []models.Artist

func ArtistsHandler(w http.ResponseWriter, r *http.Request){

	tmpl.ExecuteTemplate(w, "index.html", artistsCache)
}

func SingleArtistHandler(w http.ResponseWriter, r *http.Request){
	fmt.Println(r.URL.Path)
	path := r.URL.Path
	pathSlice := strings.Split(path, "/")
	val := pathSlice[len(pathSlice)-1]

	id, err := strconv.Atoi(val)
	if err != nil {
    	http.Error(w, "Invalid ID", http.StatusBadRequest)
    	return
	}

	for _, artist := range artistsCache{
		if artist.Id == id{
			err:= tmpl.ExecuteTemplate(w, "artist.html", artist)
			if err != nil{
				http.Error(w, "Template Error", http.StatusInternalServerError)
			}
			return
		}
	}
	http.NotFound(w, r)
}

func main(){
	var err error
	artistsCache, err = api.GetArtists()
	if err != nil{
		fmt.Println("Error fetching artists", err)
		return
	}
	
	http.HandleFunc("/artists", ArtistsHandler)
	http.HandleFunc("/artists/", SingleArtistHandler)
	log.Println("Server currently running on port:http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}