//handlers/pvz_handler.go

package handlers

import (
	"Backend_trainee_assigment_2025/internal/config"
	"Backend_trainee_assigment_2025/internal/db"
	"Backend_trainee_assigment_2025/internal/schemas"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type PVZHandler struct {
	*BaseHandler
}

func NewPVZHandler(db *db.AvitoDB, config *config.AppConfig) *PVZHandler {
	return &PVZHandler{&BaseHandler{DB: db, Config: config}}
}

func (h *PVZHandler) CreatePVZHandler(w http.ResponseWriter, r *http.Request) {
	var pvz schemas.PVZ
	if err := json.NewDecoder(r.Body).Decode(&pvz); err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if pvz.City == "" || (pvz.City != "Москва" && pvz.City != "Санкт-Петербург" && pvz.City != "Казань") {
		RespondeWithError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	newPVZ, err := h.DB.OpenPVZ(context.Background(), pvz.City)
	if err != nil {
		RespondeWithError(w, http.StatusInternalServerError, "Failed to open PVZ")

		return
	}

	ResponseWithJSON(w, http.StatusCreated, newPVZ)
}

func (h *PVZHandler) GetPVZsHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		RespondeWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	page := 1
	limit := 10
	var err error

	if requestBody["page"] != "" {
		page, err = strconv.Atoi(requestBody["page"])
		if err != nil || page < 1 {
			RespondeWithError(w, http.StatusBadRequest, "Invalid page number")
			return
		}
	}

	if requestBody["limit"] != "" {
		limit, err = strconv.Atoi(requestBody["limit"])
		if err != nil || limit < 1 || limit > 30 {
			RespondeWithError(w, http.StatusBadRequest, "Invalid limit number")
			return
		}
	}

	offset := (page - 1) * limit

	megaResponses, err := h.DB.GetPVZ(context.Background(), requestBody["startDate"], requestBody["endDate"], offset, limit)

	if err != nil {
		log.Println(err)
		RespondeWithError(w, http.StatusInternalServerError, "Failed to get PVZs")
		return
	}

	ResponseWithJSON(w, http.StatusOK, megaResponses)
}
