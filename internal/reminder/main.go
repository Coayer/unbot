package reminder

import (
	"encoding/json"
	"github.com/Coayer/unbot/internal/pkg"
	"github.com/Coayer/unbot/internal/pkg/bert"
	"io/ioutil"
	"log"
	"math"
	"strings"
	"time"
)

type Reminder struct {
	Value string
	Place string
	Time  int64
}

const REMINDERPATH = "data/reminders.json"

func GetReminders(condition string) string {
	if condition == "time" {
		return getTimeReminders()
	} else {
		return getPlaceReminders(condition)
	}
}

func SetReminder(query string) string {
	parsedTime := pkg.ParseTime(query)

	var place string
	var triggerTime int64

	for _, otherPlace := range pkg.Config.Places.Names {
		if strings.Contains(query, otherPlace) {
			place = otherPlace
			triggerTime = parsedTime
		}
	}

	//sets time as condition if no place was found
	if parsedTime != math.MaxInt64 {
		triggerTime = parsedTime
	}

	//if no place or time was found
	if place == "" && triggerTime == math.MaxInt64 {
		return "Please specify reminder conditions"
	} else {
		conditions := getConditions(query)
		value := bert.AskBert("What should I do at "+conditions+"?", query)
		writeReminders(append(readReminders(), Reminder{Value: value, Place: place, Time: triggerTime}))
		return "Reminder set for " + conditions + ": " + value
	}
}

func getConditions(query string) string {
	query = strings.ToLower(query)
	var naturalString strings.Builder

	if timeString := pkg.TimeRegex.FindString(query); timeString != "" {
		naturalString.WriteString(timeString + " ")
	}

	for _, condition := range append([]string{"today", "now", "tomorrow", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}, pkg.Config.Places.Names...) {
		for _, token := range pkg.BaseTokenize(query) {
			if token == condition {
				naturalString.WriteString(token + " ")
			}
		}
	}

	return naturalString.String()
}

func getTimeReminders() string {
	var builder strings.Builder
	var currentReminders []Reminder

	for _, reminder := range readReminders() {
		if reminder.Time <= time.Now().Unix() {
			builder.WriteString(reminder.Value + ". ")
		} else {
			currentReminders = append(currentReminders, reminder)
		}
	}

	writeReminders(currentReminders)
	return builder.String()
}

func getPlaceReminders(place string) string {
	var builder strings.Builder
	var currentReminders []Reminder

	for _, reminder := range readReminders() {
		if reminder.Place == place {
			builder.WriteString(reminder.Value + ". ")
		} else {
			currentReminders = append(currentReminders, reminder)
		}
	}

	writeReminders(currentReminders)
	return builder.String()
}

func writeReminders(reminders []Reminder) {
	data, err := json.Marshal(reminders)
	if err != nil {
		log.Fatal(err)
	}

	if ioutil.WriteFile(REMINDERPATH, data, 600) != nil {
		log.Fatal(err)
	}
}

func readReminders() []Reminder {
	var reminders []Reminder

	data, err := ioutil.ReadFile(REMINDERPATH)
	if err != nil {
		log.Fatal(err)
	}

	if json.Unmarshal(data, &reminders) != nil {
		return reminders
	}

	return reminders
}
