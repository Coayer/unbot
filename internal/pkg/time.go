package pkg

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var TimeRegex = regexp.MustCompile("\\d\\d:\\d\\d")

func ParseTime(query string) int64 {
	day := ParseDay(query)

	var hours, minutes int
	if timeString := TimeRegex.FindString(query); timeString != "" {
		splitTime := strings.Split(timeString, ":")
		hours, _ = strconv.Atoi(splitTime[0])
		minutes, _ = strconv.Atoi(splitTime[1])
	}

	currentHours, currentMinutes, currentSeconds := time.Now().Clock()
	startEpoch := time.Now().Unix() - int64(currentHours*3600+currentMinutes*60+currentSeconds)

	future := int64(day*86400 + hours*3600 + minutes*60)

	if future == 0 {
		return math.MaxInt64
	} else {
		return startEpoch + future
	}
}

func ParseDay(query string) int {
	query = strings.ToLower(query)
	currentDay := int(time.Now().Weekday())

	if strings.Contains(query, "now") || strings.Contains(query, "today") {
		return 0
	} else if strings.Contains(query, "tomorrow") || strings.Contains(query, time.Weekday(currentDay+1).String()) {
		return 1
	} else {
		for i := 1; i <= 7; i++ {
			//need currentDay relative to current currentDay
			if strings.Contains(query, strings.ToLower(time.Weekday((currentDay+i)%7).String())) {
				return i
			}
		}
	}

	return 0
}
