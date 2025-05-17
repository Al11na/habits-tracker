package domain

import "time"

type Habit struct {
	ID        string    `json:"id"`
	UserEmail string    `json:"user_email"`
	Name      string    `json:"name"`
	Goal      string    `json:"goal"`
	CreatedAt time.Time `json:"created_at"`
}
