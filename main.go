package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"stage/api"
	"stage/models"
	"stage/utils"
	"strings"
	"strconv"
)

var funcMap = template.FuncMap{
	"formatLocation": utils.FormattedDatesAndLocation,
}

var tmpl = template.Must(
	template.New("").Funcs(funcMap).ParseFiles(
		"templates/index.html",
		"templates/artist.html",
		"templates/home.html",
		"templates/search.html",
	),
)

// var tmpl = template.Must(template.ParseFiles("templates/index.html", "templates/artist.html", "templates/home.html"))

var artistsCache []models.FullArtist

func HomeHandler(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		http.NotFound(w, r)
		return
	}
	err := tmpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil{
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request){

	err := tmpl.ExecuteTemplate(w, "index.html", artistsCache)
	if err != nil{
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func SingleArtistHandler(w http.ResponseWriter, r *http.Request){
	path := r.URL.Path
	// pathSlice := strings.Split(path, "/")
	// val := pathSlice[len(pathSlice)-1]
	newPath := strings.Trim(path, "/")
	pathSlice := strings.Split(newPath, "/")
	if len(pathSlice) != 2 {
		http.NotFound(w, r)
		return
	}
	
	val := pathSlice[1]

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

func HandleSearch(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query().Get("query")
	query = strings.TrimSpace(query)
	if query == ""{
		http.Redirect(w, r, "/artists", http.StatusSeeOther)
		return
	}
	query = strings.ToLower(query)

	var result []models.FullArtist
	for _, artist := range artistsCache{
		matched := false
		if strings.Contains(strings.ToLower(artist.Name), query){
			matched = true
		}

		if !matched {
			for _, member := range artist.Members{
				if strings.Contains(strings.ToLower(member), query){
					matched = true
					break
				}
			}
		}
		if matched {
			result = append(result, artist)
		}
	}

	searchResult := models.SearchResult{
		Search: query,
		Artists: result,
	}

	err := tmpl.ExecuteTemplate(w, "search.html", searchResult)
	if err != nil{
		http.Error(w, "Template Error", http.StatusInternalServerError)
	}
}

func main(){
	var err error
	artistsCache, err = api.GetFullArtist()
	if err != nil{
		fmt.Println("Error fetching artists", err)
		return
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/artists", ArtistsHandler)
	http.HandleFunc("/artists/", SingleArtistHandler)
	http.HandleFunc("/search", HandleSearch)
	log.Println("Server currently running on port:http://localhost:8000")
	err = http.ListenAndServe(":8000", nil)
	if err != nil {
    	log.Fatal(err)
	}
}