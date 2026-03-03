package http

import (
	"encoding/json"
	"net/http"
)

// TODO: почему функцию ты сделал экспортируемой? Она где-то будет вне этого пакета использоваться?
// TODO: клиенту не надо отдавать сообщение с ошибкой. Сделай нормальную обработку сообщений в зависимости от типа ошибки
// анврапаешь ошибку и в зависимости какая ошибка будет, такой статус ей присваиваешь и решаешь, логировать или нет
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// TODO: почему ошибку не обработал?
	json.NewEncoder(w).Encode(map[string]string{"error": message})

}
