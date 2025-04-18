package db

import (
	"Backend_trainee_assigment_2025/internal/schemas"
	"fmt"
	"time"
)

func OpenPVZ(db AvitoDB, city string) (*schemas.PVZ, error) {
	query := "INSERT INTO pvzs (city) VALUES ($1) returning id, created_at"
	var id int
	var createdTime time.Time
	err := db.QueryRow(query, city).Scan(&id, &createdTime)
	if err != nil {
		return &schemas.PVZ{}, fmt.Errorf("failed to open PVZ: %w", err)
	}
	resp := schemas.PVZ{string(id), createdTime, city}
	return &resp, nil
}

func GetPVZ(db AvitoDB, startTime, endTime time.Time) ([]schemas.MegaResponse, error) {
	query := "select * from pvzs p left join receptions r on p.id = r.pvz_id left join products pr on r.id = pr.reception_id where pr.added_at between $1 and $2"
	var resp []schemas.MegaResponse
	rows, err := db.Query(query, startTime, endTime)
	if err != nil {
		return resp, err
	}
	defer rows.Close()
	for rows.Next() {
		var pvz schemas.PVZ
		var rec schemas.Reception
		var prod schemas.Product
		err := rows.Scan(&pvz.Id, &pvz.RegistrationDate, &pvz.City, &rec.Id, &rec.DateTime, &rec.PVZId, &rec.Status, &prod.Id, &prod.DateTime, &prod.ReceptionId)
		if err != nil {
			return resp, err
		}
		resp = append(resp, schemas.MegaResponse{pvz, rec, prod})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}
	return resp, nil
}
