package main

import (
	"net/http"
	"stage/api"
	"fmt"
	"log"
	"html/template"
)

var tmpl = template.Must(template.ParseFiles("templates/index.html", "templates/artist.html"))

func ArtistHandler(w http.ResponseWriter, r *http.Request){
	artist, err := api.GetArtists()
	if err != nil{
		fmt.Println(err)
		return
	}
	for _, a := range artist{
		fmt.Fprintln(w, a.Name)
	}
}

func main(){
	http.HandleFunc("/artists", ArtistHandler)
	log.Println("Server currently running on port:http://localhost:8000")
	http.ListenAndServe(":8000", nil)
}