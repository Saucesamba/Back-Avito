package handlers

import (
	"Backend_trainee_assigment_2025/internal/config"
	"Backend_trainee_assigment_2025/internal/db"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type ProductHandler struct {
	*BaseHandler
}

func NewProductHandler(db *db.AvitoDB, config *config.AppConfig) *ProductHandler {
	return &ProductHandler{&BaseHandler{DB: db, Config: config}}
}

func (h *ProductHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	typ, ok := requestBody["type"]
	if !ok {
		RespondeWithError(w, http.StatusBadRequest, "Type is required")
		return
	}

	pvzIdStr, ok := requestBody["pvzId"]
	if !ok {
		RespondeWithError(w, http.StatusBadRequest, "PvzId is required")
		return
	}
	pvzId, err := uuid.Parse(pvzIdStr)
	if err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid pvzId format (UUID required)")
		return
	}
	product, err := h.DB.CreateProduct(context.Background(), typ, pvzId)
	if err != nil {
		log.Println(err)
		RespondeWithError(w, http.StatusBadRequest, "Failed to add product")
		return
	}

	ResponseWithJSON(w, http.StatusCreated, product)
}

func (h *ProductHandler) DeleteLastProductHandler(w http.ResponseWriter, r *http.Request) {
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

	err = h.DB.DeleteProduct(context.Background(), pvzId)
	if err != nil {
		log.Println(err)
		RespondeWithError(w, http.StatusBadRequest, "no products to delete in the current reception or reception closed")
		return
	}

	w.WriteHeader(http.StatusOK)
}
