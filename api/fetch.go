package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stage/models"
	"stage/utils"
	"strings"
	"sync"
	"log"
	"os"
)

var geoMutex sync.RWMutex
var geoCache = make(map[string]models.Geolocation)

type geoResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type ArtistCache struct{
	artistsCache []models.FullArtist
	mu sync.RWMutex
}

func NewArtistCache() *ArtistCache{
	return &ArtistCache{}
}

func (a *ArtistCache) GetAllArtists()[]models.FullArtist{
	a.mu.RLock()
	defer a.mu.RUnlock()

	out := make([]models.FullArtist, len(a.artistsCache))
	copy(out, a.artistsCache)
	return out
}

func (a *ArtistCache) LoadCacheFromFile()error{
	data, err := os.ReadFile("artistsCache.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &a.artistsCache)
}

func (a *ArtistCache) SaveArtistsToCache() error{
	data, err := json.MarshalIndent(a.artistsCache, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile("artistsCache.json", data, 0644)
}

func (a *ArtistCache) Refresh() error{
	artists, err := GetArtists()
	if err != nil {
		return  err
	}

	relation, err := GetRelations()
	if err != nil {
		return err
	}

	var newCache []models.FullArtist

	semaphore := make(chan struct{}, 5)

	for _, artist := range artists {

		locationMap := relation[artist.Id]

		// create slice with correct size
		locations := make([]models.LocationInfo, len(locationMap))

		var wg sync.WaitGroup

		// channel to safely collect results
		type result struct {
			location models.LocationInfo
		}

		results := make(chan result, len(locationMap))

		// launch goroutines
		for location, dates := range locationMap {

			wg.Add(1)

			go func(location string, dates []string) {

				defer wg.Done()

				semaphore <- struct{}{}

				defer func() {
					<-semaphore
				}()

				coordinates, err := GeocodeLocation(location)
				if err != nil {
					log.Println("geocode failed:", err)
					coordinates = models.Geolocation{}
				}

				results <- result{
					location: models.LocationInfo{
						Name:  location,
						Lat:   coordinates.Lat,
						Lon:   coordinates.Lon,
						Dates: dates,
					},
				}

			}(location, dates)
		}

		// wait for all goroutines
		wg.Wait()

		// close results channel
		close(results)

		// safely fill locations slice
		i := 0
		for res := range results {
			locations[i] = res.location
			i++
		}

		info := models.FullArtist{
			Artist:         artist,
			DatesLocations: relation[artist.Id],
			Locations:      locations,
		}

		newCache = append(newCache, info)
	}
	a.mu.Lock()
	a.artistsCache = newCache
	a.mu.Unlock()

	// SaveArtistsToCache()
	err = a.SaveArtistsToCache()
	if err != nil{
		return err
	}
	
	SaveCacheToFile()
	return nil
}

// var artistsCache []models.FullArtist

func GeocodeLocation(location string) (models.Geolocation, error) {
	formattedLocation := utils.FormatForGeocoding(location)

	cacheKey := strings.ToLower(formattedLocation)
	geoMutex.RLock()
	result, found := geoCache[cacheKey]
	geoMutex.RUnlock()

	if found {
		return result, nil
	}

	url := "https://nominatim.openstreetmap.org/search?q=" +
		strings.ReplaceAll(formattedLocation, " ", "+") +
		"&format=json&limit=1"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request", err)
		return models.Geolocation{}, err
	}
	req.Header.Set("User-Agent", "groupie-tracker-app")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return models.Geolocation{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Geolocation{}, err
	}

	var data []geoResponse

	err = json.Unmarshal(body, &data)
	if err != nil {
		return models.Geolocation{}, err
	}

	if len(data) == 0 {
		fmt.Println("No location found")
		return models.Geolocation{}, fmt.Errorf("no location found")
	}

	result = models.Geolocation{
		Lat: data[0].Lat,
		Lon: data[0].Lon,
	}

	geoMutex.Lock()
	geoCache[cacheKey] = result
	geoMutex.Unlock()

	return result, nil
}

func GetArtists() ([]models.Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("Error fetching api", err)
		return []models.Artist{}, err
	}
	defer resp.Body.Close()

	fmt.Println("Status code is:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad Status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []models.Artist{}, err
	}

	var artists []models.Artist

	err = json.Unmarshal(body, &artists)
	if err != nil {
		return []models.Artist{}, err
	}
	return artists, nil
}

func GetRelations() (map[int]map[string][]string, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad Status %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var relations models.Relations

	err = json.Unmarshal(data, &relations)
	if err != nil {
		return nil, err
	}

	relationMap := make(map[int]map[string][]string)

	for _, item := range relations.Index {
		relationMap[item.Id] = item.DatesLocations
	}
	return relationMap, nil
}
