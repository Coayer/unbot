package plane

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Coayer/unbot/internal/utils"
	"log"
	"os"
	"strings"
)

var airlines = loadCsv("internal/plane/data/airlines.csv")

var apiURL = fmt.Sprintf("https://opensky-network.org/api/states/all?lamin=%f&lomin=%f&lamax=%f&lomax=%f",
	utils.LAT-0.3, utils.LON-0.3, utils.LAT+0.3, utils.LON+0.3)

type OpenSkyFetch struct {
	States [][]string
}

func GetPlane() string {
	log.Println(apiURL)

	stateVectors := parsePlanes(utils.HttpGet(apiURL))
	var result strings.Builder

	for _, vector := range stateVectors.States {
		result.WriteString(formatCallsign(vector[1]))
	}
	return result.String()
}

func formatCallsign(callsign string) string {
	if callsign == "" {
		return "No callsign recieved"
	}

	icaoAirline, flightNumber := callsign[:3], callsign[3:]

	for _, entry := range airlines {
		if icaoAirline == entry[1] {
			return entry[2] + " " + flightNumber
		}
	}

	return "Can't find airline"
}

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
