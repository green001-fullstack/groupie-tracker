package api

import(
	"net/http"
	"fmt"
	"io"
	"stage/models"
	"encoding/json"
)

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