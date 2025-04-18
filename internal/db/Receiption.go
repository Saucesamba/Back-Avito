package db

import (
	"Backend_trainee_assigment_2025/internal/schemas"
	"fmt"
	"strconv"
	"time"
)

func OpenRec(db AvitoDB, pvzId string) (*schemas.Reception, error) {

	query := "INSERT INTO receiptions (pvz_id) VALUES ($1) returning id, created_at"
	var id int
	var createdTime time.Time
	intId, _ := strconv.Atoi(pvzId)
	err := db.QueryRow(query, intId).Scan(&id, &createdTime)

	if err != nil {
		return &schemas.Reception{}, fmt.Errorf("failed to open reception: %w", err)
	}
	resp := schemas.Reception{string(id), createdTime, pvzId, "in_progress"}
	return &resp, nil
}

func CloseLastRec(db AvitoDB, pvzId string) (*schemas.Reception, error) {

	query := "Update receiptions set status='closed' where pvz_id = $1 and id = (select id from receptions where pvz_id=$1 order by date_time desc limit 1)) returning id, created_at"

	intId, _ := strconv.Atoi(pvzId)
	var id int
	var createdTime time.Time
	err := db.QueryRow(query, intId).Scan(&id, &createdTime)
	if err != nil {
		return &schemas.Reception{}, fmt.Errorf("failed to close last reception: %w", err)
	}
	resp := schemas.Reception{string(id), createdTime, pvzId, "closed"}
	return &resp, nil
}
