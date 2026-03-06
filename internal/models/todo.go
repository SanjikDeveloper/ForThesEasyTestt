package models

import "time"

type Todo struct {
	ID          int       `json:"id"`
	List        string    `json:"list"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
