package api

import(
	"encoding/json"
	"os"
)

func SaveCacheToFile() error {
	data, err := json.MarshalIndent(geoCache, "", " ")
	if err != nil{
		return err
	}
	return os.WriteFile("geocache.json", data, 0644)
}

func LoadCacheFromFile() error{
	data, err := os.ReadFile("geocache.json")
	if err != nil{
		return err
	}

	return json.Unmarshal(data, &geoCache)
}

