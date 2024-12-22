package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func infixToPostfix(expression string) ([]string, error) {
	var postfix []string // результат в ОПН
	var operators []rune // стек для операторов

	for i := 0; i < len(expression); i++ {
		char := rune(expression[i])
		// Пропускаем пробелы из значения
		if char == ' ' {
			continue
		}

		if unicode.IsDigit(char) || char == '.' {
			num, n := parseNumber(expression[i:])
			postfix = append(postfix, num)

			// требутеся для корректной работы после парсигна значение
			// Указываем с какого идекса продолжать программу
			i += n - 1
			continue
		}

		if char == '(' {
			operators = append(operators, char)
		} else if char == ')' {
			for len(operators) > 0 && operators[len(operators)-1] != '(' {
				postfix = append(postfix, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 || operators[len(operators)-1] != '(' {
				return nil, errors.New("несоответствующие скобки")
			}
			operators = operators[:len(operators)-1]
		} else if isOperator(char) {
			for len(operators) > 0 && getOperatorPriority(operators[len(operators)-1]) >= getOperatorPriority(char) {
				postfix = append(postfix, string(operators[len(operators)-1]))
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, char)
		} else {
			return nil, fmt.Errorf("неверный символ: %c", char)
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == '(' {
			return nil, errors.New("несоответствующие скобки")
		}
		postfix = append(postfix, string(operators[len(operators)-1]))
		operators = operators[:len(operators)-1]
	}

	return postfix, nil
}

func evaluatePostfix(postfix []string) (float64, error) {
	var stack []float64

	for _, token := range postfix {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else if len(token) == 1 && isOperator(rune(token[0])) {
			if len(stack) < 2 {
				return 0, errors.New("недостаточно операндов")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			var result float64
			switch token {
			case "+":
				result = a + b
			case "-":
				result = a - b
			case "*":
				result = a * b
			case "/":
				if b == 0 {
					return 0, errors.New("деление на ноль")
				}
				result = a / b
			}
			stack = append(stack, result)
		} else {
			return 0, fmt.Errorf("неверный токен: %s", token)
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("ошибка: неверное выражение")
	}
	return stack[0], nil
}

func parseNumber(s string) (string, int) {
	var numStr strings.Builder
	i := 0
	for i < len(s) && (unicode.IsDigit(rune(s[i])) || s[i] == '.') {
		numStr.WriteByte(s[i])
		i++
	}
	return numStr.String(), i
}

func isOperator(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/'
}

func getOperatorPriority(operator rune) int {
	switch operator {
	case '+', '-':
		return 1
	case '*', '/':
		return 2
	default:
		return 0
	}
}

func Calc(expression string) (float64, error) {
	postfix, err := infixToPostfix(expression)

	if err != nil {
		return 0, err
	}
	return evaluatePostfix(postfix)
}

type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, errMsg string) {
	if err := writeJSONResponse(w, statusCode, ErrorResponse{Error: errMsg}); err != nil {
		// На случай, если даже при кодировании ответа произошла ошибка
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusBadRequest, "Only POST is allowed")
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusUnprocessableEntity, "Invalid JSON input")
		return
	}

	if req.Expression == "" {
		respondWithError(w, http.StatusUnprocessableEntity, "expression parameter is missing")
		return
	}

	result, errCalc := Calc(req.Expression)
	if errCalc != nil {
		respondWithError(w, http.StatusUnprocessableEntity, errCalc.Error())
		return
	}

	resp := Response{
		Result: result,
	}
	if err := writeJSONResponse(w, http.StatusOK, resp); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}
func main() {
	http.HandleFunc("/api/v1/calculate", calculateHandler)

	log.Println("[DEBUG] We start the server on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[CRIT] Server startup error: %v", err)
	}
}
