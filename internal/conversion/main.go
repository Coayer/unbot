package conversion

import (
	"fmt"
	"github.com/Coayer/unbot/internal/pkg"
	"math"
	"strconv"
)

var length = map[string]float64{"km": 1000, "miles": 1609, "m": 1, "meters": 1, "feet": 0.3048}

var mass = map[string]float64{"kg": 1, "kilos": 1, "stone": 6.35, "lbs": 0.45359, "pounds": 0.45359, "g": 0.001,
	"grams": 0.001, "ounces": 0.02835, "oz": 0.02835}

var temperature = map[string]bool{"Celsius": true, "c": true, "Fahrenheit": false, "f": false}

func Convert(query string) string {
	value, unit1, unit2 := parseConversion(query)

	if isLength(unit1) && isLength(unit2) {
		return fmt.Sprintf("%.2f %s", value*length[unit1]/length[unit2], unit2)
	} else if isMass(unit1) && isMass(unit2) {
		return fmt.Sprintf("%.2f %s", value*mass[unit1]/mass[unit2], unit2)
	} else if isTemperature(unit1) && isTemperature(unit2) {
		var result int

		if temperature[unit1] {
			result = int(math.Round(value*1.8 + 32))
		} else {
			result = int(math.Round((value - 32) * 0.555556))
		}

		return strconv.Itoa(result) + " degrees " + unit2
	}

	return "Invalid conversion"
}

func parseConversion(query string) (float64, string, string) {
	tokens := pkg.BaseTokenize(pkg.RemoveStopWords(query))

	var value float64
	var unit1, unit2 string

	for _, token := range tokens {
		if pkg.IsNumeric(token) {
			value, _ = strconv.ParseFloat(token, 64)
		}

		if isUnit(token) && value != 0 && unit1 == "" {
			unit1 = token
		} else if isUnit(token) && value != 0 && unit2 == "" {
			unit2 = token
		}
	}

	return value, unit1, unit2
}

func isUnit(token string) bool {
	return isLength(token) || isMass(token) || isTemperature(token)
}

func isMass(token string) bool {
	for unit := range mass {
		if unit == token {
			return true
		}
	}

	return false
}

func isLength(token string) bool {
	for unit := range length {
		if unit == token {
			return true
		}
	}

	return false
}

func isTemperature(token string) bool {
	for unit := range temperature {
		if unit == token {
			return true
		}
	}

	return false
}
