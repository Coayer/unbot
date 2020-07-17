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

var airlines = loadCsv("data/airlines.csv")

var apiURL = fmt.Sprintf("https://opensky-network.org/api/states/all?lamin=%f&lomin=%f&lamax=%f&lomax=%f",
	utils.LAT-0.3, utils.LON-0.3, utils.LAT+0.3, utils.LON+0.3)

type OpenSkyFetch struct {
	States [][17]interface{}
}

//GetPlane is used by calling code to run the package
func GetPlane() string {
	log.Println(apiURL)
	stateVectors := parsePlanes(utils.HttpGet(apiURL))
	plane := closestPlane(stateVectors)
	return formatCallsign(plane[1].(string))
}

func closestPlane(stateVectors OpenSkyFetch) [17]interface{} {
	minDistance := math.Inf(1)
	var plane [17]interface{}
	for _, vector := range stateVectors.States {
		log.Println(vector[1])

		distance := math.Pow((vector[5].(float64)-utils.LON)*math.Cos(utils.LAT), 2) + math.Pow(vector[6].(float64)-utils.LAT, 2)
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

//loadCsv is used to load the ICAO airline codes to their full names
func loadCsv(path string) map[string]string {
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

	hashTable := make(map[string]string)

	for _, line := range records {
		hashTable[line[0]] = line[1]
	}

	return hashTable
}
