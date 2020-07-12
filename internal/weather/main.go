package weather

import (
	"encoding/json"
	"fmt"
	"github.com/Coayer/unbot/internal/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

func GetWeather() string {
	log.Println(apiURL)
	weather := utils.HttpGet(apiURL)
	return parseCurrentWeather(weather)
}

func parseCurrentWeather(bytes []byte) string {
	var response OWMFetch

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Println(err)
		return "Could not get weather"
	}

	var description strings.Builder

	for _, condition := range response.Current.Weather {
		description.WriteString(condition.Description + " ")
	}

	return fmt.Sprintf("%s, %d degrees, %d percent humidity",
		description.String(), int(response.Current.Temp), response.Current.Humidity)
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
