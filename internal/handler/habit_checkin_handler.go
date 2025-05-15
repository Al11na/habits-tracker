package handler

import (
	"net/http"

	"habit-tracker-api/internal/service"

	"github.com/gin-gonic/gin"
)

type HabitCheckinHandler struct {
	service *service.HabitCheckinService
}

func NewHabitCheckinHandler(s *service.HabitCheckinService) *HabitCheckinHandler {
	return &HabitCheckinHandler{s}
}

// CheckinRequest — тело POST /habits/:id/checkin
type CheckinRequest struct {
	Comment string `json:"comment"`
}

// POST /habits/:id/checkin
func (h *HabitCheckinHandler) CheckIn(c *gin.Context) {
	habitID := c.Param("id")
	var req CheckinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CheckIn(habitID, req.Comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "checked in"})
}

// GET /habits/:id/stats
func (h *HabitCheckinHandler) Stats(c *gin.Context) {
	habitID := c.Param("id")
	streak, total, possible, rate, err := h.service.Stats(habitID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"streak":          streak,
		"total_checks":    total,
		"possible_checks": possible,
		"completion_rate": rate,
	})
}

// internal/handler/habit_checkin_handler.go

func (h *HabitCheckinHandler) Report(c *gin.Context) {
	habitID := c.Param("id")
	report, err := h.service.Report(habitID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, report)
}
