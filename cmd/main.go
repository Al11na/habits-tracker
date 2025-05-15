package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"habit-tracker-api/internal/auth"
	"habit-tracker-api/internal/handler"
	"habit-tracker-api/internal/repository"
	"habit-tracker-api/internal/service"
)

func main() {
	// 1) Инициализируем BoltDB
	repository.InitDB()
	defer repository.DB.Close()

	// 2) Зависимости для авторизации
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// 3) Зависимости для CRUD привычек
	habitRepo := repository.NewHabitRepository()
	habitService := service.NewHabitService(habitRepo)
	habitHandler := handler.NewHabitHandler(habitService)

	// 4) Создаём Gin-роутер
	r := gin.Default()

	// 5) Регистрируем публичные маршруты: регистрация и логин
	userHandler.RegisterRoutes(r)

	// 6) Регистрируем CRUD для привычек с заглушкой middleware
	// authMiddleware пока только ставит userEmail в контекст
	authMiddleware := func(c *gin.Context) {
		c.Set("userEmail", "test@example.com")
		c.Next()
	}
	habitHandler.RegisterRoutes(r, authMiddleware)

	// Check‑in и Stats
	checkinRepo := repository.NewHabitCheckinRepository()
	checkinService := service.NewHabitCheckinService(habitRepo, checkinRepo)
	checkinHandler := handler.NewHabitCheckinHandler(checkinService)

	// Регистрируем внутри той же группы /habits
	grp := r.Group("/habits", authMiddleware)
	{
		grp.POST("/:id/checkin", checkinHandler.CheckIn)
		grp.GET("/:id/stats", checkinHandler.Stats)
		grp.GET("/:id/report", checkinHandler.Report)
	}

	// 7) Пример защищённого route /api/me с настоящим JWT middleware
	protected := r.Group("/api")
	protected.Use(auth.AuthMiddleware())
	{
		protected.GET("/me", func(c *gin.Context) {
			email := c.GetString("userEmail")
			c.JSON(200, gin.H{"email": email})
		})
	}

	// 8) Health‑check на корневом /
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Habit Tracker API running with BoltDB!"})
	})

	// 9) Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
