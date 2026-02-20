package models

type Todo struct {
	IdList      int    `json:"id_list"`
	TodoList    string `json:"todo_list"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}
