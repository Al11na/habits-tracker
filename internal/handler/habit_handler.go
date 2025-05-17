package handler

import (
	"net/http"
	"strconv"
	"time"

	"habit-tracker-api/internal/domain"
	"habit-tracker-api/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// HabitHandler — HTTP-хендлер для CRUD привычек
type HabitHandler struct {
	service *service.HabitService
}

func NewHabitHandler(s *service.HabitService) *HabitHandler {
	return &HabitHandler{s}
}

// CreateHabitRequest — тело POST /habits
type CreateHabitRequest struct {
	Name string `json:"name" binding:"required"`
	Goal string `json:"goal"`
}

// CreateHabit — создаёт новую привычку
func (h *HabitHandler) CreateHabit(c *gin.Context) {
	var req CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userEmail := c.GetString("userEmail")
	habit := &domain.Habit{
		ID:        uuid.New().String(),
		UserEmail: userEmail,
		Name:      req.Name,
		Goal:      req.Goal,
		CreatedAt: time.Now(),
	}

	if err := h.service.Create(habit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, habit)
}

// GetHabits — возвращает отфильтрованный и постраничный список привычек
func (h *HabitHandler) GetHabits(c *gin.Context) {
	userEmail := c.GetString("userEmail")

	name := c.Query("name")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	habits, err := h.service.ListHabits(userEmail, name, dateFrom, dateTo, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      habits,
		"page":      page,
		"page_size": pageSize,
		"total":     len(habits), // можно улучшить, если нужен общий count
	})
}

// GetHabit — возвращает одну привычку по ID
func (h *HabitHandler) GetHabit(c *gin.Context) {
	id := c.Param("id")
	habit, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "habit not found"})
		return
	}
	c.JSON(http.StatusOK, habit)
}

// UpdateHabit — обновляет привычку по ID
func (h *HabitHandler) UpdateHabit(c *gin.Context) {
	id := c.Param("id")
	var req CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "habit not found"})
		return
	}
	existing.Name = req.Name
	existing.Goal = req.Goal

	if err := h.service.Update(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, existing)
}

// DeleteHabit — удаляет привычку по ID
func (h *HabitHandler) DeleteHabit(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "habit deleted"})
}

// RegisterRoutes — привязывает маршруты CRUD привычек к роутеру
func (h *HabitHandler) RegisterRoutes(r *gin.Engine, authMiddleware gin.HandlerFunc) {
	grp := r.Group("/habits", authMiddleware)
	{
		grp.POST("/", h.CreateHabit)
		grp.GET("/", h.GetHabits)
		grp.GET("/:id", h.GetHabit)
		grp.PUT("/:id", h.UpdateHabit)
		grp.DELETE("/:id", h.DeleteHabit)
	}
}
