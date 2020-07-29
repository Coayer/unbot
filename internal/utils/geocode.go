package utils

import (
	"encoding/json"
	"log"
	"strconv"
)

type Place struct {
	Lat string
	Lon string
}

var CurrentPlace string

func UpdatePlace(query string) {
	CurrentPlace = GetEntities(query + "?")
}

func GetLocation() (float64, float64) {
	bytes := HttpGet("https://nominatim.openstreetmap.org/search?countrycodes=gb&format=json&q=" +
		FormatHTTPQuery(CurrentPlace))

	var response []Place

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Println(err)
	}

	latitude, _ := strconv.ParseFloat(response[0].Lat, 64)
	longitude, _ := strconv.ParseFloat(response[0].Lon, 64)

	return latitude, longitude
}
