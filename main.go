package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	// Устанавливаем соединение с базой данных
	var err error
	db, err = sql.Open("mysql", "root:root@tcp(localhost:3306)/test")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Инициализируем маршрутизатор Gin
	router := gin.Default()

	// Регистрируем обработчики маршрутов
	router.POST("/users", createUser)
	router.POST("/quests", createQuest)
	router.POST("/complete", completeQuest)
	router.GET("/history/:userId", getUserHistory)

	// Запускаем сервер
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// Обработчик для создания нового пользователя
func createUser(c *gin.Context) {
	// Реализация создания пользователя
	// Парсинг данных из JSON тела запроса
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка корректности входных данных
	if user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Выполнение SQL запроса для добавления пользователя
	result, err := db.Exec("INSERT INTO users (name, balance) VALUES (?, ?)", user.Name, user.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Получение ID нового пользователя
	userID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	// Возвращаем успешный ответ с ID созданного пользователя
	c.JSON(http.StatusCreated, gin.H{"userID": userID})

}

// Обработчик для создания нового задания
func createQuest(c *gin.Context) {
	// Реализация создания задания
	// Парсинг данных из JSON тела запроса
	var quest Quest
	if err := c.ShouldBindJSON(&quest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка корректности входных данных
	if quest.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	if quest.Cost <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cost must be greater than 0"})
		return
	}

	// Выполнение SQL запроса для добавления задания
	result, err := db.Exec("INSERT INTO quests (name, cost) VALUES (?, ?)", quest.Name, quest.Cost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quest"})
		return
	}

	// Получение ID нового задания
	questID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get quest ID"})
		return
	}

	// Возвращаем успешный ответ с ID созданного задания
	c.JSON(http.StatusCreated, gin.H{"questID": questID})
}

// Инициализация пяти заданий по информационной безопасности
func initSecurityQuests() {
	securityQuests := []Quest{
		{Name: "What is phishing?", Cost: 10},
		{Name: "Explain the concept of encryption", Cost: 15},
		{Name: "What is a firewall and how does it work?", Cost: 12},
		{Name: "Describe the difference between symmetric and asymmetric encryption", Cost: 20},
		{Name: "What is a DDoS attack?", Cost: 18},
	}

	for _, quest := range securityQuests {
		_, err := db.Exec("INSERT INTO quests (name, cost) VALUES (?, ?)", quest.Name, quest.Cost)
		if err != nil {
			log.Fatalf("Failed to initialize security quest: %s", err)
		}
	}

}

// Обработчик для завершения задания
func completeQuest(c *gin.Context) {
	// Реализация завершения задания
	// Парсинг данных из JSON тела запроса
	var completion struct {
		UserID  int `json:"user_id"`
		QuestID int `json:"quest_id"`
	}
	if err := c.ShouldBindJSON(&completion); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка корректности входных данных
	if completion.UserID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if completion.QuestID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quest ID"})
		return
	}

	// Проверка выполнения условия, что пользователь может выполнить задание только один раз
	var completed int
	if err := db.QueryRow("SELECT COUNT(*) FROM completed_quests WHERE user_id = ? AND quest_id = ?", completion.UserID, completion.QuestID).Scan(&completed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check completion status"})
		return
	}
	if completed > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User has already completed this quest"})
		return
	}

	// Получение стоимости задания
	var cost int
	if err := db.QueryRow("SELECT cost FROM quests WHERE id = ?", completion.QuestID).Scan(&cost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get quest cost"})
		return
	}

	// Выполнение SQL запроса для добавления выполненного задания
	_, err := db.Exec("INSERT INTO completed_quests (user_id, quest_id) VALUES (?, ?)", completion.UserID, completion.QuestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete quest"})
		return
	}

	// Обновление баланса пользователя
	_, err = db.Exec("UPDATE users SET balance = balance + ? WHERE id = ?", cost, completion.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user balance"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"message": "Quest completed successfully"})
}

// Обработчик для получения истории выполненных заданий и баланса пользователя
func getUserHistory(c *gin.Context) {
	// Реализация получения истории выполненных заданий и баланса пользователя
	// Парсинг данных из параметров запроса
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Проверка корректности входных данных (идентификатор пользователя)
	var userExists bool
	if err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&userExists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}
	if !userExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Выполнение SQL запроса для получения истории выполненных заданий и баланса пользователя
	rows, err := db.Query("SELECT q.id, q.name, q.cost FROM quests q JOIN completed_quests c ON q.id = c.quest_id WHERE c.user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user history"})
		return
	}
	defer rows.Close()

	var history []map[string]interface{}
	for rows.Next() {
		var questID, questCost int
		var questName string
		if err := rows.Scan(&questID, &questName, &questCost); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user history"})
			return
		}
		history = append(history, map[string]interface{}{
			"id":   questID,
			"name": questName,
			"cost": questCost,
		})
	}

	// Получение текущего баланса пользователя
	var balance int
	if err := db.QueryRow("SELECT balance FROM users WHERE id = ?", userID).Scan(&balance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user balance"})
		return
	}

	// Возвращаем историю выполненных заданий и баланс пользователя в виде JSON ответа
	c.JSON(http.StatusOK, gin.H{
		"history": history,
		"balance": balance,
	})

}
