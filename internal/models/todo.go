package models

import "time"

// TODO: зачем тут три параметра используют указатель?
type Todo struct {
	// TODO: тут итак понятно чье это айди, не надо писать IdList
	// TODO: Переменные которые названы аббреиватурно, надо писать все с заглавной буквы. Пример: ID, URL
	IdList int `json:"id_list"`
	// TODO: итак понятно, что за лист, в названии не надо слово Todo
	TodoList    *string    `json:"todo_list"`
	Description *string    `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
}
