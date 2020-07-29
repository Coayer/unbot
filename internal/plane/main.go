package plane

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Coayer/unbot/internal/utils"
	"log"
	"math"
	"os"
)

var airlines = loadAirlines(loadCsv("data/airlines.csv"))
var aircraft = loadAircraft(loadCsv("data/aircraft.csv"))

var apiURL = fmt.Sprintf("https://opensky-network.org/api/states/all?lamin=%f&lomin=%f&lamax=%f&lomax=%f",
	utils.Config.Location.Latitude-0.3, utils.Config.Location.Longitude-0.3, utils.Config.Location.Latitude+0.3, utils.Config.Location.Longitude+0.3)

type OpenSkyFetch struct {
	States [][17]interface{}
}

//GetPlane is used by calling code to run the package
func GetPlane() string {
	log.Println(apiURL)
	stateVectors := parsePlanes(utils.HttpGet(apiURL))
	plane := closestPlane(stateVectors)
	return formatPlane(plane)
}

func formatPlane(vector [17]interface{}) string {
	if vector[8] == "true" {
		return formatCallsign(vector[1].(string)) + " on ground"
	} else {
		return fmt.Sprintf("%s, %s, bearing %d, at %d feet, heading %d", formatCallsign(vector[1].(string)),
			aircraft[vector[0].(string)], bearingToPlane(vector), int(vector[7].(float64)*3.281), int(vector[10].(float64)))
	}
}

func bearingToPlane(vector [17]interface{}) int {
	planeLong, planeLat := vector[5].(float64), vector[6].(float64)
	y := math.Sin(planeLong-utils.Config.Location.Longitude) * math.Cos(planeLat)
	x := math.Cos(utils.Config.Location.Latitude)*math.Sin(planeLat) -
		math.Sin(utils.Config.Location.Latitude)*math.Cos(planeLat)*math.Cos(planeLong-utils.Config.Location.Longitude)
	return int(math.Mod(math.Atan2(y, x)*180/math.Pi+360, 360))
}

func closestPlane(stateVectors OpenSkyFetch) [17]interface{} {
	minDistance := math.Inf(1)
	var plane [17]interface{}
	for _, vector := range stateVectors.States {
		log.Println(vector[1])

		distance := math.Pow((vector[5].(float64)-utils.Config.Location.Longitude)*math.Cos(utils.Config.Location.Latitude), 2) +
			math.Pow(vector[6].(float64)-utils.Config.Location.Latitude, 2)
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
