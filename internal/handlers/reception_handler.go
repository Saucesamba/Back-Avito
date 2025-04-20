package handlers

import (
	"Backend_trainee_assigment_2025/internal/config"
	"Backend_trainee_assigment_2025/internal/db"
	"Backend_trainee_assigment_2025/internal/schemas"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ReceptionHandler struct {
	*BaseHandler
}

func NewReceptionHandler(db *db.AvitoDB, config *config.AppConfig) *ReceptionHandler {
	return &ReceptionHandler{&BaseHandler{DB: db, Config: config}}
}

func (h *ReceptionHandler) CreateReceptionHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	pvzIdStr, ok := requestBody["pvzId"]
	if !ok {
		RespondeWithError(w, http.StatusBadRequest, "Missing pvzId in request body")
		return
	}

	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid pvzId format (UUID required)")
		return
	}

	newReception, err := h.DB.OpenRec(context.Background(), pvzId)
	newReception.Products = []schemas.Product{}
	if err != nil {
		log.Println(err)
		RespondeWithError(w, http.StatusInternalServerError, "Failed to open reception")
		return
	}
	if newReception.Status == "unable" {
		RespondeWithError(w, http.StatusBadRequest, "You already have opened reception")
		return
	}

	ResponseWithJSON(w, http.StatusCreated, newReception)
}

func (h *ReceptionHandler) CloseLastReceptionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pvzIdStr, ok := vars["pvzId"]
	if !ok {
		RespondeWithError(w, http.StatusBadRequest, "Missing pvzId in path")
		return
	}

	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid pvzId format (UUID required)")
		return
	}

	updatedReception, err := h.DB.CloseLastRec(context.Background(), pvzId)
	if updatedReception.Status == "failed" {
		RespondeWithError(w, http.StatusBadRequest, "All receptions has been already closed")
		return
	}

	if err != nil {
		RespondeWithError(w, http.StatusInternalServerError, "")
		return
	}
	updatedReception.Products, err = h.DB.GetProduct(updatedReception.Id)
	ResponseWithJSON(w, http.StatusOK, updatedReception)

}
