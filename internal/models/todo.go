package models

import "time"

type Todo struct {
	IdList      int        `json:"id_list"`
	TodoList    *string    `json:"todo_list"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
}
