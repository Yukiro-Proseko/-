Чудо-калькулятор работает на честном слове разработчика. Включает в себя “errors”, “strconv”, “net/http”, “fmt”.
Главная функция main обращается к функции Resultic, и выводит этим самым результат(при наличии), ошибку(при наличии) и код ошибки(200/422/500)

func main() {  // main
    http.HandleFunc("/", Resultic)
    fmt.Println("Сервер запущен на порту :8080")
    http.ListenAndServe(":8080", nil)
   }

Пакет называется main

Функция Resultic принимает введённую строку формата 
"expression": "выражение, которое ввёл пользователь" 
Потом вычисляет через вторую главную функцию Calc результат

Функция Calc занимается организацией работы всех остальных побочных функций. 
Calc:
    A-evaluateExpression:
        B-infixToPostfix:
            C-tokenize:
                D-isOperator
            C-isNumber
            C-isOperator
            C-precedence
        B-isOperator
        B-performOperation
Выше показано что во что входит.
