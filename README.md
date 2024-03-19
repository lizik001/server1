# server 1

Проект REST API сервиса для учета заданий пользователей.

## Установка и запуск

1. Клонировать репозиторий:

git clone https://github.com/lizik001/server1.git


2. Перейти в каталог проекта:

cd myproject

3. Установить зависимости (если есть):

npm install

4. Запустить приложение:

go run main.go
shell

## Описание API методов

### Метод создания пользователя

POST /users

Параметры запроса:

- name (string): Имя пользователя.

Формат ответа:

{
"id": 1,
"name": "John",
"balance": 0
}
shell

### Метод создания задания

POST /quests
торий:

git clone https:
- name (string): Название задания.
- cost (number): Стоимость задания.

Формат ответа:

{
"id": 1,
"name": "Task 1",
"cost": 10
}
shell

### Метод завершения задания

POST /complete
й:

git clone https://gi
- user_id (number): Идентификатор пользователя.
- quest_id (number): Идентификатор задания.

Формат ответа:

{
"message": "Quest completed successfully"
}
shell

### Метод получения истории выполненных заданий и баланса пользователя

GET /users/:userId/history

Параметры запроса:

- userId (number): Идентификатор пользователя.

Формат ответа:

{
"history": [
{
"id": 1,
"name": "Task 1",
"cost": 10
},
{
"id": 2,
"name": "Task 2",
"cost": 15
}
],
"balance": 25
}
shell

## Примеры использования API

### Создание пользователя

curl -X POST -H "Content-Type: application/json" -d '{"name":"John"}' http://localhost:8000/users
shell

### Создание задания

curl -X POST -H "Content-Type: application/json" -d '{"name":"Task 1","cost":10}' http://localhost:8000/quests
shell

### Завершение задания

curl -X POST -H "Content-Type: application/json" -d '{"user_id":1,"quest_id":1}' http://localhost:8000/complete
shell

### Получение истории выполненных заданий и баланса пользователя

curl http://localhost:8000/users/1/history


Замените http://localhost:8000 на адрес вашего сервера, если он отличается.

Этот файл README.md содержит инструкции по установке и запуску вашего приложен