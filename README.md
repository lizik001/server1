# server 1

Проект REST API сервиса для учета заданий пользователей.

## Для исользования backend вам понадобится :

Установите MySQL на компьютер https://www.mysql.com/

Авторизуйтесь с помощью команды 

```sql
ALTER USER 'root'@'localhost' IDENTIFIED BY 'root';
```

Создайте у себя базу данных, назвав ее "test" с омощью команды 

```sql
CREATE DATABASE TEST;
```

Вставьте следующие команды: 

```sql
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    balance INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS quests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    cost INT DEFAULT 0
);

CREATE TABLE IF NOT EXISTS completed_quests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    quest_id INT NOT NULL,
    completion_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (quest_id) REFERENCES quests(id)
);
```

### После запуска backend используйте программу Postman или любой другой клиент API  

## Установка и запуск

1. Клонировать репозиторий:

git clone https://github.com/lizik001/server1.git


2. Перейти в каталог проекта:

cd myproject

3. Установить зависимости:

go mod tidy

4. Запустить приложение:

go run main.go


## Описание API методов

### Метод создания пользователя

POST http://localhost:8080/users

Параметры запроса:

- name (string): Имя пользователя.

Формат ответа:

```json
{
  "userID": 2
}
```

### Метод создания задания

POST http://localhost:8080/quests


git clone https:
- name (string): Название задания.
- cost (number): Стоимость задания.

Формат ответа:

```json
{
  "questID": 2
}
```


### Метод завершения задания

POST http://localhost:8080/complete:


- user_id (number): Идентификатор пользователя.
- quest_id (number): Идентификатор задания.

Формат ответа:

```json
{
"message": "Quest completed successfully"
}
```

### Метод получения истории выполненных заданий и баланса пользователя

GET http://localhost:8080/history/:userId

Параметры запроса:

- userId (number): Идентификатор пользователя.

Формат ответа:

```json
{
  "balance": 10,
  "history": [
    {
      "cost": 10,
      "id": 1,
      "name": "user0"
    }
  ]
}
```


## Примеры использования API

### Создание пользователя

curl -X POST -H "Content-Type: application/json" -d '{"name":"John"}' http://localhost:8080/users


### Создание задания

curl -X POST -H "Content-Type: application/json" -d '{"name":"Task 1","cost":10}' http://localhost:8080/quests


### Завершение задания

curl -X POST -H "Content-Type: application/json" -d '{"user_id":1,"quest_id":1}' http://localhost:8080/complete


### Получение истории выполненных заданий и баланса пользователя

curl http://localhost:8080/history/1

