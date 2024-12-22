package main

import ( 
    "errors" 
    "strconv"
	"net/http"
	"fmt"
)

func Calc(expression string) (float64, error, bool) { 
    result, err, h := evaluateExpression(expression) 
    if err != nil { 
        return 0, err, h
    } 
    return result, nil, h
}
//*************************************************************************************
//создаём bool переменную для ошибки ttt,где если что-то не так то false

func evaluateExpression(expression string) (float64, error, bool) { 
    var ttt bool
	// Проверяем, что выражение не пустое 
    if len(expression) == 0 { 
        return 0, errors.New("empty expression"), ttt
    }

    // Преобразуем строку в обратную польскую запись
    postfixExpression, err, bol := infixToPostfix(expression)  // тут тоже может быть ошибка. обрабатываем
    if err != nil {
	    return 0, err, bol
    }

    // Вычисляем значение выражения в обратной польской записи
    stack := make([]float64, 0)
    for _, token := range postfixExpression {
	    if isOperator(token) {
		    if len(stack) < 2 {
			    return 0, errors.New("invalid expression"), ttt
		    }
		    operand2 := stack[len(stack)-1]
		    operand1 := stack[len(stack)-2]
		    stack = stack[:len(stack)-2]
		    result := performOperation(operand1, operand2, token)
		    stack = append(stack, result)
	    } else {
		    value, err := strconv.ParseFloat(token, 64)
		    if err != nil {
			    return 0, errors.New("invalid expression"), ttt
		    }
		    stack = append(stack, value)
	    }
    }

    if len(stack) != 1 {
	    return 0, errors.New("invalid expression"), ttt
    }

    return stack[0], nil, ttt
}

func infixToPostfix(expression string) ([]string, error, bool) { 
    var result []string 
    var operators []string
    var h bool
    //в tokenize ошибок нет 
	tokens := tokenize(expression)

    for _, token := range tokens {
	    if isNumber(token) {
		    result = append(result, token)
	    } else if token == "(" {
		    operators = append(operators, token)
	    } else if token == ")" {
		    for len(operators) > 0 && operators[len(operators)-1] != "(" {
			    result = append(result, operators[len(operators)-1])
			    operators = operators[:len(operators)-1]
		    }
		    if len(operators) == 0 {
				return nil, errors.New("unmatched parentheses"), h
		    }
		    operators = operators[:len(operators)-1]
	    } else if isOperator(token) {
		    for len(operators) > 0 && precedence(operators[len(operators)-1]) >= precedence(token) {
			    result = append(result, operators[len(operators)-1])
			    operators = operators[:len(operators)-1]
		    }
		    operators = append(operators, token)
	    } else {
		    h = false
			return nil, errors.New("invalid token"), h
			
	    }
    }

    for len(operators) > 0 {
	    if operators[len(operators)-1] == "(" {
			return nil, errors.New("unmatched parentheses"), h
	    }
	    result = append(result, operators[len(operators)-1])
	    operators = operators[:len(operators)-1]
    }

    return result, nil, h
}

func tokenize(expression string) []string { //            !!!
    var tokens []string 
    token := ""

    for _, char := range expression {
	    if isOperator(string(char)) || string(char) == "(" || string(char) == ")" {
		    if token != "" {
			    tokens = append(tokens, token)
			    token = ""
		    }
		    tokens = append(tokens, string(char))
	    } else if string(char) == " " {
		    continue
	    } else {
		    token += string(char)
	    }
    }

    if token != "" {
	    tokens = append(tokens, token)
    }

    return tokens
}

func isNumber(token string) bool { 
    _, err := strconv.ParseFloat(token, 64) 
    return err == nil 
}

func isOperator(token string) bool { 
    return token == "+" || token == "-" || token == "*" || token == "/" 
}

func precedence(operator string) int { 
    switch operator { 
    case "+", "-": 
        return 1 
    case "*", "/": 
        return 2 
    default: 
        return 0 } 
}

func performOperation(operand1, operand2 float64, operator string) float64 { 
    switch operator { 
    case "+": 
        return operand1 + operand2 
    case "-": 
        return operand1 - operand2 
    case "*": 
        return operand1 * operand2 
    case "/": 
        return operand1 / operand2 
    default: 
        return 0 } 
}
//********************************************************
func Resultic(w http.ResponseWriter, r *http.Request)  { 
    expression := r.URL.Query().Get("expression")
    res, err, h := Calc(expression)
    s := fmt.Sprint(res)
	if err != nil && h == true {
		fmt.Fprint(w, "error: ", err)
		http.ResponseWriter.WriteHeader(w, 422)
    } else if err != nil && h != true {
        fmt.Fprint(w, "error: ", err)
		http.ResponseWriter.WriteHeader(w, 500)
	} else {
        fmt.Fprint(w, "result: ", s)
        http.ResponseWriter.WriteHeader(w, 200)
	}
}


func main() {  // main
	http.HandleFunc("/", Resultic)
	fmt.Println("Сервер запущен на порту :8080")
	http.ListenAndServe(":8080", nil)
   }