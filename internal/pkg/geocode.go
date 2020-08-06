package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type Place struct {
	Lat string
	Lon string
}

var CurrentPlace string

func UpdatePlace(query string) {
	CurrentPlace = GetEntities(query + "?") //needed for prose NER to work
}

func GetLocation() (float64, float64) {
	bytes := HttpGet(fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&countrycodes=%s&q=%s",
		Config.Location.Country, FormatHTTPQuery(CurrentPlace)))

	var response []Place

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Println(err)
	}

	latitude, _ := strconv.ParseFloat(response[0].Lat, 64)
	longitude, _ := strconv.ParseFloat(response[0].Lon, 64)

	return latitude, longitude
}
