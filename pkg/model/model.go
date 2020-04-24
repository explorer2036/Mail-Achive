package model

import (
	"time"
)

// Email structure for achive
type Email struct {
	From      string    `json:"from"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

// User for login and refresh
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Manager  bool   `json:"manager"`
}
