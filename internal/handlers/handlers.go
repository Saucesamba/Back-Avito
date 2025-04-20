package handlers

import (
	"Backend_trainee_assigment_2025/internal/config"
	"Backend_trainee_assigment_2025/internal/db"
	"Backend_trainee_assigment_2025/internal/schemas"
	"encoding/json"
	"net/http"
)

type BaseHandler struct {
	DB     db.Database
	Config *config.AppConfig
}

func NewBaseHandler(DB db.Database, Config *config.AppConfig) *BaseHandler {
	return &BaseHandler{
		DB:     DB,
		Config: Config,
	}
}

func RespondeWithError(w http.ResponseWriter, code int, message string) {
	ResponseWithJSON(w, code, schemas.Error{message})
}

func ResponseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
