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
	Sunrise int64
	Sunset  int64
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
	var apiURL string

	location := pkg.GetLocation(query)
	apiURL = fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?units=metric&lat=%f&lon=%f&exclude=minutely,hourly,current&appid=%s",
		location.Latitude, location.Longitude, pkg.Config.OwmKey)

	log.Println(apiURL)
	query = strings.ToLower(query)
	weather := parseWeather(pkg.HttpGet(apiURL))

	return generateDescription(weather[pkg.ParseDay(query)], query)
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

	if weather.Uvi > 2 && weather.Sunrise < time.Now().Unix() && weather.Sunset > time.Now().Unix() {
		description.WriteString("UV index " + strconv.Itoa(int(math.Round(weather.Uvi))))
	}

	return description.String()
}

func formatTime(epoch int64) string {
	return time.Unix(epoch, 0).Format("15:04")
}

func parseWeather(bytes []byte) []DailyWeather {
	var response struct{ Daily []DailyWeather }

	if err := json.Unmarshal(bytes, &response); err != nil {
		log.Println(err)
	}

	return response.Daily
}
