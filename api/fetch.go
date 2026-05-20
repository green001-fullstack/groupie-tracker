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
)

var geoCache = make(map[string]models.Geolocation)

type geoResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func GeocodeLocation(location string) (models.Geolocation, error) {
	formattedLocation := utils.FormatForGeocoding(location)

	cacheKey := strings.ToLower(formattedLocation)
	if result, found := geoCache[cacheKey]; found {
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

	if len(data) == 0{
		fmt.Println("No location found")
		return models.Geolocation{}, fmt.Errorf("no location found")
	}

	result := models.Geolocation{
		Lat : data[0].Lat,
		Lon : data[0].Lon,
	}

	geoCache[cacheKey] = result

	return result, nil
}

func GetFullArtist() ([]models.FullArtist, error) {
	artists, err := GetArtists()
	if err != nil {
		return nil, err
	}
	relation, err := GetRelations()
	if err != nil {
		return nil, err
	}

	var ArrayFullArtist []models.FullArtist

	for _, artist := range artists {

		locationMap := relation[artist.Id]

		locations := make([]models.LocationInfo, len(locationMap))

		var wg sync.WaitGroup

		i := 0

		for location, dates := range locationMap{
			index := i
			i++

			wg.Add(1)

			go func(index int, location string, dates []string){
				defer wg.Done()

				coordinates, err := GeocodeLocation(location)
				if err != nil{
					coordinates = models.Geolocation{}
				}
				locations[index] = models.LocationInfo{
					Name: location,
					Lat: coordinates.Lat,
					Lon: coordinates.Lon,
					Dates: dates,
				}
			}(index, location, dates)
		}



		info := models.FullArtist{
			Artist:         artist,
			DatesLocations: relation[artist.Id],
			Locations: locations,
		}

		ArrayFullArtist = append(ArrayFullArtist, info)
	}

	return ArrayFullArtist, nil
}

func GetArtists() ([]models.Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		fmt.Println("Error fetching api", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("Status code is:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad Status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var artists []models.Artist

	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
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
