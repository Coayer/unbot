package weather

import (
	"encoding/json"
	"fmt"
	"github.com/Coayer/unbot/internal/utils"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

var apiURL = fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?units=metric&lat=%f&lon=%f&exclude=minutely,hourly,current&appid=%s",
	utils.LAT, utils.LON, loadKey())

type DailyWeather struct {
	Sunrise int
	Sunset  int
	Temp    struct {
		Day   float64
		Night float64
		Eve   float64
		Morn  float64
	}
	Humidity uint8
	Weather  []struct {
		Description string
	}
}

func GetWeather(query string) string {
	log.Println(apiURL)
	query = strings.ToLower(query)
	weather := parseWeather(utils.HttpGet(apiURL))
	day := int(time.Now().Weekday())

	if strings.Contains(query, "now") || strings.Contains(query, "today") {
		return generateDescription(weather[0], query)
	} else if strings.Contains(query, "tomorrow") || strings.Contains(query, time.Weekday(day+1).String()) {
		return generateDescription(weather[1], query)
	} else {
		for i := 1; i <= 7; i++ {
			//owm gives weather relative to current day, not to start of week
			if strings.Contains(query, strings.ToLower(time.Weekday((day+i)%7).String())) {
				return generateDescription(weather[i], query)
			}
		}
	}

	return "No weather found"
}

func generateDescription(weather DailyWeather, query string) string {
	if strings.Contains(query, "sunset") {
		return formatTime(weather.Sunset)
	} else if strings.Contains(query, "sunrise") {
		return formatTime(weather.Sunrise)
	} else {
		return weatherDescription(weather, query)
	}
}

func weatherDescription(weather DailyWeather, query string) string {
	var description strings.Builder

	for _, condition := range weather.Weather {
		description.WriteString(condition.Description + " ")
	}

	var temperature float64

	if strings.Contains(query, "morning") {
		temperature = weather.Temp.Morn
	} else if strings.Contains(query, "evening") {
		temperature = weather.Temp.Eve
	} else if strings.Contains(query, "night") {
		temperature = weather.Temp.Night
	} else {
		temperature = weather.Temp.Day
	}

	return fmt.Sprintf("%s, %d degrees, %d percent humidity", description.String(), int(math.Round(temperature)),
		weather.Humidity)
}

func formatTime(epoch int) string {
	return time.Unix(int64(epoch), 0).Format("15:04")
}

func parseWeather(bytes []byte) []DailyWeather {
	var response struct{ Daily []DailyWeather }

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Println(err)
	}

	return response.Daily
}

func loadKey() string {
	file, err := os.Open("data/owm_key.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	key, _ := ioutil.ReadAll(file)
	return string(key)
}
