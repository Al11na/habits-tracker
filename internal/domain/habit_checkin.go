package domain

import "time"

// HabitCheckin — запись о выполнении привычки на определённую дату
type HabitCheckin struct {
	ID      string    `json:"id"`
	HabitID string    `json:"habit_id"`
	Date    time.Time `json:"date"`              // обычно без времени
	Comment string    `json:"comment,omitempty"` // опционально
}
