# Calc Service

**Calc Service** — это простой сервис на Go, который принимает на вход математическое выражение, вычисляет его и возвращает результат.  
Сервис демонстрирует базовые приёмы работы с HTTP, JSON и простейшими арифметическими выражениями.

## Оглавление
1. [Особенности проекта](#особенности-проекта)
2. [Установка и запуск](#установка-и-запуск)
3. [Примеры использования](#примеры-использования)
    1. [Успешный вызов (200)](#успешный-вызов-200)
    2. [Ошибка валидации (422)](#ошибка-валидации-422)
    3. [Внутренняя ошибка (500)](#внутренняя-ошибка-500)
4. [Контакты](#контакты)

---

## Особенности проекта
- **HTTP API**: принимает POST-запросы на эндпоинт `/api/v1/calculate`.
- **JSON**: входные данные (поле `expression`), выходные данные (поле `result`).
- **Ошибка**: в случае невалидного ввода возвращает JSON с описанием ошибки и соответствующим HTTP-статусом.

Пример входных данных:
```json
{
  "expression": "2+2*2"
}
```
Пример выходных данных при успешном вычислении:
```json
{
  "result": 6
}
```
В случае ошибки возвращается JSON вида:
```json
{
  "error": "division by zero is not allowed"
}
```

---

## Установка и запуск

1. Склонируйте репозиторий:
   ```bash
   git clone https://github.com/RickDevQQQ/yandex_calculator.git
   ```
2. Перейдите в директорию проекта:
   ```bash
   cd calc_service
   ```
3. Запустите сервис (предполагается, что у вас установлен Go):
   ```bash
   go run main.go
   ```
4. По умолчанию сервис слушает на порту `8080`.

> **Примечание**: Порт и другие параметры могут быть изменены в исходном коде.

---

## Примеры использования

Ниже — примеры, как отправлять запросы к сервису с помощью `curl`.  
Обратите внимание на заголовок `Content-Type: application/json`, который обязателен.

### Успешный вызов (200)

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
     --header 'Content-Type: application/json' \
     --data '{
       "expression": "2+2*2"
     }'
```

**Ожидаемый ответ (HTTP 200 OK)**:
```json
{
  "result": 6
}
```

---

### Ошибка валидации (422)
Пример ситуации, когда мы вообще не передали поле `expression` или передали пустую строку.

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
     --header 'Content-Type: application/json' \
     --data '{
       "expression": ""
     }'
```
**Ожидаемый ответ (HTTP 422 Unprocessable Entity)**:
```json
{
  "error": "expression parameter is missing"
}
```