package weather

import (
	"encoding/json"
	"fmt"
	"github.com/Coayer/unbot/internal/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var apiURL = fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?units=metric&lat=%f&lon=%f&exclude=minutely,hourly&appid=%s",
	utils.LAT, utils.LON, loadKey())

type OWMFetch struct {
	Current struct {
		Temp     float32
		Humidity uint8
		Weather  []struct {
			Description string
		}
	}
	Daily []struct {
		Temp struct {
			Day   float32
			Night float32
			Eve   float32
			Morn  float32
		}
		Weather []struct {
			Description string
		}
	}
}

func GetWeather(query string) string {
	log.Println(apiURL)
	query = strings.ToLower(query)
	weather := parseWeather(utils.HttpGet(apiURL))
	day := int(time.Now().Weekday())

	var description strings.Builder
	var result string

	if strings.Contains(query, "now") || strings.Contains(query, "today") {
		for _, condition := range weather.Current.Weather {
			description.WriteString(condition.Description + " ")
		}

		result = fmt.Sprintf("%s, %d degrees, %d percent humidity", description.String(), int(weather.Current.Temp),
			weather.Current.Humidity)
	} else if strings.Contains(query, "tomorrow") || strings.Contains(query, time.Weekday(day+1).String()) {
		for _, condition := range weather.Daily[1].Weather {
			description.WriteString(condition.Description + " ")
		}

		result = fmt.Sprintf("%s, %d degrees", description.String(), int(weather.Daily[1].Temp.Day))
	} else {
		for i := 1; i <= 7; i++ {
			//owm gives weather relative to current day, not to start of week
			if strings.Contains(query, strings.ToLower(time.Weekday((day+i)%7).String())) {
				for _, condition := range weather.Daily[i].Weather {
					description.WriteString(condition.Description + " ")
				}

				result = fmt.Sprintf("%s, %d degrees", description.String(), int(weather.Daily[i].Temp.Day))
				break
			}
		}
	}

	return result
}

func parseWeather(bytes []byte) OWMFetch {
	var response OWMFetch

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Println(err)
	}

	return response
}

func loadKey() string {
	file, err := os.Open("owm_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	key, _ := ioutil.ReadAll(file)
	return string(key)
}
