package weather

import (
	"encoding/json"
	"fmt"
	"github.com/Coayer/unbot/internal/pkg"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
)

type DailyWeather struct {
	Sunrise int
	Sunset  int
	Temp    struct {
		Day   float64
		Night float64
		Eve   float64
		Morn  float64
	}
	Humidity int
	Weather  []struct {
		Description string
	}
	Uvi float64
}

func GetWeather(query string) string {
	pkg.UpdatePlace(query)
	var apiURL string

	if pkg.CurrentPlace == "" {
		apiURL = fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?units=metric&lat=%f&lon=%f&exclude=minutely,hourly,current&appid=%s",
			pkg.Config.Location.Latitude, pkg.Config.Location.Longitude, pkg.Config.OwmKey)
	} else {
		latitude, longitude := pkg.GetLocation()
		apiURL = fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?units=metric&lat=%f&lon=%f&exclude=minutely,hourly,current&appid=%s",
			latitude, longitude, pkg.Config.OwmKey)
	}
	log.Println(apiURL)
	query = strings.ToLower(query)
	weather := parseWeather(pkg.HttpGet(apiURL))
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
		description.WriteString(condition.Description + ", ")
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

	description.WriteString(strconv.Itoa(int(math.Round(temperature))) + " degrees, ")

	if weather.Humidity > 45 {
		description.WriteString(strconv.Itoa(weather.Humidity) + "% humidity, ")
	}

	if weather.Uvi > 2 {
		description.WriteString("UV index " + strconv.Itoa(int(math.Round(weather.Uvi))))
	}

	return description.String()
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
