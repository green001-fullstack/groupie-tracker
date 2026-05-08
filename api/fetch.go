package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"stage/models"
)

func GetFullArtist()([]models.FullArtist, error){
	artists, err := GetArtists()
	if err != nil{
		return nil, err
	}
	relation, err := GetRelations()
	if err != nil{
		return nil, err
	}

    var ArrayFullArtist []models.FullArtist

	for _, artist := range artists{

		info := models.FullArtist{
			 Artist:artist,
             DatesLocations:relation[artist.Id],
		}

		ArrayFullArtist = append(ArrayFullArtist, info)
	}

    return ArrayFullArtist, nil


}

func GetArtists()([]models.Artist, error){
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil{
		fmt.Println("Error fetching api", err)
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Println("Status code is:", resp.Status)

	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("Bad Status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}
	
	var artists []models.Artist

	err = json.Unmarshal(body, &artists)
	if err != nil{
		return nil, err
	}
	return artists, nil
}

func GetRelations()(map[int]map[string][]string, error){
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("Bad Status %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil{
		return nil, err
	}

	var relations models.Relations

	err = json.Unmarshal(data, &relations)
	if err != nil{
		return nil, err
	}

	relationMap := make(map[int]map[string][]string)

	for _, item := range relations.Index{
		relationMap[item.Id] = item.DatesLocations
	}
	return relationMap, nil
}

