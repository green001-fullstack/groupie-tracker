package main

import (
	"log"
	"net/http"
	"stage/api"
	"stage/handlers"
	"time"
)

func main() {
	var err error

	cache := api.NewArtistCache()
	err = cache.LoadCacheFromFile()
	if err != nil{
		log.Println("No cache file found, starting fresh")
	}
	

	if len(cache.GetAllArtists()) == 0{
		err := cache.Refresh()
		if err != nil{
			log.Fatal("Failed to refresh cache:", err)
		}
	}

		go func() {
		for {
			log.Println("Refreshing artist cache in background...")

			err := cache.Refresh()
			if err != nil {
				log.Println("Background cache refresh failed:", err)
			} else {
				log.Println("Background cache updated successfully")
			}

			time.Sleep(24 * time.Hour)
		}
	}()
	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/artist", handlers.ArtistsHandler)
	http.HandleFunc("/singleArtists/", handlers.SingleArtistHandler)
	log.Println("Server currently running on port:http://localhost:8001")
	err = http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal(err)
	}
}
