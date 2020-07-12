package calculator

import (
	"fmt"
	"github.com/Coayer/unbot/internal/utils"
	"log"
	"math"
	"strconv"
	"strings"
)

func Evaluate(query string) string {
	expression := parseExpression(query)
	if len(expression) == 0 {
		return "Invalid expression"
	} else {
		return evaluateExpression(expression)
	}
}

func evaluateExpression(expression []string) string {
	operatorPrecedence := map[string]uint8{"": 0, "-": 1, "+": 2, "x": 3, "/": 4, "^": 5}
	result := ""

	for {
		log.Println(expression)

		if len(expression) == 1 {
			result = expression[0]
			break
		} else if len(expression) == 3 {
			result = calculate(expression)
			break
		}

		operator := ""
		operatorIndex := -1

		for i, unit := range expression {
			if isOperator(unit) {
				if operatorPrecedence[unit] > operatorPrecedence[operator] {
					operator = unit
					operatorIndex = i
				}
			}
		}

		tempSlice := expression[operatorIndex+2:]
		expression = append(expression[:operatorIndex-1], calculate(expression[operatorIndex-1:operatorIndex+2]))
		expression = append(expression, tempSlice...)
	}

	return formatResult(result)
}

func formatResult(result string) string {
	decimalSplit := strings.Split(result, ".")
	fmt.Println(decimalSplit)

	if decimalSplit[1][0:2] == "00" {
		return decimalSplit[0]
	} else {
		return result
	}
}

func calculate(expression []string) string {
	x1, _ := strconv.ParseFloat(expression[0], 64)
	x2, _ := strconv.ParseFloat(expression[2], 64)
	operator := expression[1]

	switch operator {
	case "-":
		return fmt.Sprintf("%f", x1-x2)
	case "+":
		return fmt.Sprintf("%f", x1+x2)
	case "x":
		return fmt.Sprintf("%f", x1*x2)
	case "/":
		return fmt.Sprintf("%f", x1/x2)
	case "^":
		return fmt.Sprintf("%f", math.Pow(x1, x2))
	default:
		return ""
	}
}

func parseExpression(query string) []string {
	tokens := utils.BaseTokenize(query)
	startToken := 0

	for i, token := range tokens {
		if isNumeric(token) {
			startToken = i
			break
		}
	}

	for i, token := range tokens[startToken:] {
		if !((i%2 == 0 && isNumeric(token)) || (i%2 != 0 && isOperator(token))) {
			return []string{}
		}
	}

	expression := tokens[startToken:]

	if len(expression)%2 == 0 {
		return []string{}
	} else {
		return expression
	}
}

func isNumeric(token string) bool {
	for _, char := range token {
		if (char < '0' || char > '9') && char != '.' {
			return false
		}
	}
	return true
}

func isOperator(token string) bool {
	if len(token) == 1 && (token == "+" || token == "-" || token == "x" || token == "/" || token == "^") {
		return true
	} else {
		return false
	}
}
