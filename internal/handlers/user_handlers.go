package handlers

import (
	"Backend_trainee_assigment_2025/internal/auth"
	"Backend_trainee_assigment_2025/internal/config"
	"Backend_trainee_assigment_2025/internal/db"
	"Backend_trainee_assigment_2025/internal/schemas"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type UserHandler struct {
	*BaseHandler
}

func NewUserHandler(db db.Database, config *config.AppConfig) *UserHandler {
	return &UserHandler{&BaseHandler{DB: db, Config: config}}
}

func (h *UserHandler) DummyLoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	role := requestBody["role"]
	if role != "moderator" && role != "employee" {
		RespondeWithError(w, http.StatusBadRequest, "Invalid role")
		return
	}
	//Для дамилогин создаем новый ууид
	userID := uuid.New()

	token, err := auth.GenerateToken(userID, role, h.Config.JWT.Secret, h.Config.JWT.TokenExpiryHours)
	if err != nil {
		RespondeWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	ResponseWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user schemas.UserReg
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if user.Email == "" || user.Password == "" || user.Role == "" {
		RespondeWithError(w, http.StatusBadRequest, "All fields are required")
		return
	}
	newUser, err := h.DB.CreateUser(context.Background(), &user)
	if err != nil {
		RespondeWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	ResponseWithJSON(w, http.StatusCreated, newUser)
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user schemas.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if user.Email == "" || user.Password == "" {
		RespondeWithError(w, http.StatusBadRequest, "All fields are required")
		return
	}
	existingUser, err := h.DB.GetUser(context.Background(), user)
	if err != nil {
		RespondeWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}
	userID := existingUser.Id // Get from your code where user ID is accessible
	role := existingUser.Role
	token, err := auth.GenerateToken(userID, role, h.Config.JWT.Secret, h.Config.JWT.TokenExpiryHours)

	if err != nil {
		RespondeWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	ResponseWithJSON(w, http.StatusOK, map[string]string{"token": token})
}
