package plane

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Coayer/unbot/internal/pkg"
	"log"
	"math"
	"os"
)

var airlines = loadAirlines(loadCsv("data/airlines.csv"))
var aircraft = loadAircraft(loadCsv("data/aircraft.csv"))

type OpenSkyFetch struct {
	States [][17]interface{}
}

func GetPlane(query string) string {
	distance := 0.3
	location := pkg.GetLocation(query)

	openSkyURL := fmt.Sprintf("https://opensky-network.org/api/states/all?lamin=%f&lomin=%f&lamax=%f&lomax=%f",
		location.Latitude-distance, location.Longitude-distance, location.Latitude+distance, location.Longitude+distance)
	log.Println(openSkyURL)
	stateVectors := parsePlanes(pkg.HttpGet(openSkyURL))

	if stateVectors.States == nil {
		return "No planes found"
	}

	plane := closestPlane(stateVectors, location)
	return formatPlane(plane, location)
}

func formatPlane(vector [17]interface{}, location pkg.Place) string {
	if vector[8] == "true" {
		return formatCallsign(vector[1].(string)) + " on ground"
	} else {
		return fmt.Sprintf("%s, %s, %s, at %d feet, heading %s", formatCallsign(vector[1].(string)),
			aircraft[vector[0].(string)], directionToPlane(vector, location), int(vector[7].(float64)*3.281),
			bearingCardinal(vector[10].(float64)))
	}
}

func directionToPlane(vector [17]interface{}, location pkg.Place) string {
	planeLong, planeLat := vector[5].(float64), vector[6].(float64)
	y := math.Sin(planeLong-location.Longitude) * math.Cos(planeLat)
	x := math.Cos(location.Latitude)*math.Sin(planeLat) -
		math.Sin(location.Latitude)*math.Cos(planeLat)*math.Cos(planeLong-location.Longitude)
	bearing := math.Mod(math.Atan2(y, x)*180/math.Pi+360, 360)

	return bearingCardinal(bearing)
}

func bearingCardinal(bearing float64) string {
	directions := []string{"north", "north-northeast", "north-east", "east-northeast", "east", "east-southeast",
		"south-east", "south-southeast", "south", "south-southwest", "south-west", "west-southwest", "west",
		"west-northwest", "north-west", "north-northwest"}

	return directions[int((bearing+11.25)/22.5)%16]
}

func closestPlane(stateVectors OpenSkyFetch, location pkg.Place) [17]interface{} {
	minDistance := math.Inf(1)
	var plane [17]interface{}
	for _, vector := range stateVectors.States {
		log.Println(vector[1])

		distance := math.Pow((vector[5].(float64)-location.Longitude)*math.Cos(location.Latitude), 2) +
			math.Pow(vector[6].(float64)-location.Latitude, 2)
		if distance < minDistance {
			plane = vector
			minDistance = distance
		}
	}
	return plane
}

//formatCallsign splits a callsign into the full airline name and flight number
func formatCallsign(callsign string) string {
	if callsign == "" {
		return "No callsign recieved"
	}

	icaoAirline, flightNumber := callsign[:3], callsign[3:]

	if airline, exists := airlines[icaoAirline]; exists {
		return airline + " " + flightNumber
	} else {
		return callsign
	}
}

//parsePlanes is used to unmarshal the OpenSky JSON response
func parsePlanes(bytes []byte) OpenSkyFetch {
	var response OpenSkyFetch

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Println(err)
	}

	return response
}

func loadCsv(path string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records
}

func loadAirlines(records [][]string) map[string]string {
	log.Println("Loading airlines")
	hashTable := make(map[string]string)

	for _, line := range records {
		hashTable[line[0]] = line[1]
	}

	return hashTable
}

func loadAircraft(records [][]string) map[string]string {
	log.Println("Loading aircraft")
	hashTable := make(map[string]string)

	for _, line := range records {
		hashTable[line[0]] = line[2]
	}

	return hashTable
}
