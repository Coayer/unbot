package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type Place struct {
	Latitude  float64
	Longitude float64
}

func GetLocation(query string) Place {
	unknownPlace := GetEntities(query + "?")

	if unknownPlace == "" {
		return Config.Places.Default
	} else {
		bytes := HttpGet(fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&countrycodes=%s&q=%s",
			Config.Country, FormatHTTPQuery(unknownPlace)))

		var response []struct {
			Lat string
			Lon string
		}

		if err := json.Unmarshal(bytes, &response); err != nil {
			log.Println(err)
		}

		latitude, _ := strconv.ParseFloat(response[0].Lat, 64)
		longitude, _ := strconv.ParseFloat(response[0].Lon, 64)

		return Place{Latitude: latitude, Longitude: longitude}
	}
}
