package main

import (
	"html/template"
	"log"
	"net/http"
	"slices"
	"sort"
	"stage/api"
	"stage/models"
	"stage/utils"
	"strconv"
	"strings"
	// "fmt"
	// "sync"
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
		"templates/error.html",
	),
)

// var tmpl = template.Must(template.ParseFiles("templates/index.html", "templates/artist.html", "templates/home.html"))

var artistsCache []models.FullArtist

func RenderError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(code)
    data := models.ErrorPage{
        Code:    code,
        Message: message,
    }
    if err := tmpl.ExecuteTemplate(w, "error.html", data); err != nil {
        log.Println("template execution error:", err)
    }
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	if r.URL.Path != "/" {
		RenderError(w, http.StatusNotFound, "Page Not found")
		return
	}
	err := tmpl.ExecuteTemplate(w, "home.html", nil)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

func ArtistsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filter := r.URL.Query().Get("sort")
	filter = strings.TrimSpace(filter)

	filteredArtist := slices.Clone(artistsCache)

	switch filter {
	case "ascending":
		sort.Slice(filteredArtist, func(i, j int) bool {
			return filteredArtist[i].Name < filteredArtist[j].Name
		})
	case "descending":
		sort.Slice(filteredArtist, func(i, j int) bool {
			return filteredArtist[i].Name > filteredArtist[j].Name
		})
	case "oldest":
		sort.Slice(filteredArtist, func(i, j int) bool {
			return filteredArtist[i].CreationDate < filteredArtist[j].CreationDate
		})
	case "newest":
		sort.Slice(filteredArtist, func(i, j int) bool {
			return filteredArtist[i].CreationDate > filteredArtist[j].CreationDate
		})
	case "default":
		// No sorting, keep the original order
		filteredArtist = artistsCache
	}

	err := tmpl.ExecuteTemplate(w, "index.html", filteredArtist)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Template Error")
		return
	}
}

func SingleArtistHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	newPath := strings.Trim(path, "/")
	pathSlice := strings.Split(newPath, "/")

	if len(pathSlice) != 2 {
		RenderError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	val := pathSlice[1]

	id, err := strconv.Atoi(val)
	if err != nil {
		RenderError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	for _, artist := range artistsCache {
		if artist.Id == id {

			// var locations []models.LocationInfo
			// locations := make([]models.LocationInfo, len(artist.DatesLocations))

			// i := 0

			// var wg sync.WaitGroup
			// // var mu sync.Mutex

			// for location, dates := range artist.DatesLocations {

			// 	wg.Add(1)

			// 	index := i

				// go func(location string, dates []string) {
				// 	defer wg.Done()

				// 	coordinates, err := api.GeocodeLocation(location)
				// 	if err != nil {
				// 		coordinates = models.Geolocation{}
				// 	}

				// 	mu.Lock()
				// 	locations = append(locations, models.LocationInfo{
				// 		Name:  location,
				// 		Lat:   coordinates.Lat,
				// 		Lon:   coordinates.Lon,
				// 		Dates: dates,
				// 	})
				// 	mu.Unlock()

				// }(location, dates)

			// 	go func(index int, location string, dates []string){

			// 		defer wg.Done()

			// 		coordinates, err := api.GeocodeLocation(location)
			// 		if err != nil{
			// 			coordinates = models.Geolocation{}
			// 		}

			// 		locations[index] = models.LocationInfo{
			// 			Name: location,
			// 			Lat: coordinates.Lat,
			// 			Lon: coordinates.Lon,
			// 			Dates: dates,
			// 		}
			// 	}(index, location, dates)
			// 	i++
			// }

			// wg.Wait()

			pageData := models.ArtistPageData{
				Artist:         artist.Artist,
				DatesLocations: artist.DatesLocations,
				Locations:      artist.Locations,
			}

			err := tmpl.ExecuteTemplate(w, "artist.html", pageData)
			if err != nil {
				RenderError(w, http.StatusInternalServerError, "Template Error")
				return
			}
			return
		}
	}

	RenderError(w, http.StatusNotFound, "Artist not found")
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		RenderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	query := r.URL.Query().Get("query")
	query = strings.TrimSpace(query)
	if query == "" {
		http.Redirect(w, r, "/artists", http.StatusSeeOther)
		return
	}
	query = strings.ToLower(query)

	var result []models.FullArtist
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
			for location := range artist.DatesLocations{
				if strings.Contains(strings.ToLower(location), query){
					matched = true
					break
				}
			}
		}

		if !matched {
			creationDate := strconv.Itoa(artist.CreationDate)
			if strings.Contains(creationDate, query){
				matched = true
				break
			}
		}

		if  !matched {
			if strings.Contains(artist.FirstAlbum, query){
				matched = true
				break
			}
		}
		if matched {
			result = append(result, artist)
		}
	}

	searchResult := models.SearchResult{
		Search:  query,
		Artists: result,
	}

	err := tmpl.ExecuteTemplate(w, "search.html", searchResult)
	if err != nil {
		RenderError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

func main() {
	var err error
	err = api.LoadCacheFromFile()
	if err != nil{
		log.Println("No cache file found, starting fresh")
	}
	artistsCache, err = api.GetFullArtist()
	// fmt.Println(artistsCache[0])
	if err != nil {
		log.Fatal("Error fetching artists: ", err)
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
