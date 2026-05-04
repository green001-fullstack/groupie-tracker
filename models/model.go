package models


type Artist struct {
    Id           int	`json:"id"`
    Image        string	`json:"image"`
    Name         string	`json:"name"`
    Members      []string	`json:"members"`
    CreationDate int	`json:"creationDate"`
    FirstAlbum   string	`json:"firstAlbum"`
}

type RelationItem struct{
	Id int	`json:"id"`
	DatesLocations map[string][]string	`json:"datesLocations"`
}

type Relations struct{
	Index []RelationItem `json:"index"`
}

type LocationItem struct{
	Id int `json:"id"`
	Locations []string `json:"locations"`
	DatesUrl string `json:"dates"`
}

type Locations struct{
	Index []LocationItem `json:"index"`
}

type DateItem struct{
	Id int	`json:"id"`
	Dates []string	`json:"dates"`
}

type Dates struct{
	Index []DateItem	`json:"index"`
}

