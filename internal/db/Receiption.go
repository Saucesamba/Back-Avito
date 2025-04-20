package db

import (
	"Backend_trainee_assigment_2025/internal/schemas"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (a *AvitoDB) OpenRec(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error) {
	checkStatusQ := "select * from receptions order by date_time desc limit 1"
	recept := schemas.Reception{}
	row := a.QueryRow(checkStatusQ)
	err := row.Scan(&recept.Id, &recept.DateTime, &recept.PVZId, &recept.Status)
	var id uuid.UUID
	var createdTime time.Time
	if recept.Status == "in_progress" {
		return &schemas.Reception{id, createdTime, id, "unable", []schemas.Product{}}, nil
	}
	query := "INSERT INTO receptions (pvz_id,status) VALUES ($1, $2) returning id, date_time"

	err = a.QueryRowContext(ctx, query, pvzId, "in_progress").Scan(&id, &createdTime) // scan values from query

	if err != nil {
		return &schemas.Reception{}, fmt.Errorf("failed to open reception: %w", err)
	}
	resp := schemas.Reception{id, createdTime, pvzId, "in_progress", []schemas.Product{}}
	return &resp, nil
}

func (a *AvitoDB) CloseLastRec(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error) {
	query := "Update receptions set status='closed' where pvz_id = $1 and id = (select id from receptions where pvz_id=$1 and status='in_progress' order by date_time desc limit 1) returning id, date_time"
	var id uuid.UUID
	var createdTime time.Time
	err := a.QueryRowContext(ctx, query, pvzId).Scan(&id, &createdTime)
	if id == uuid.Nil {
		return &schemas.Reception{Status: "failed"}, nil
	}
	if err != nil {
		return &schemas.Reception{}, fmt.Errorf("failed to close last reception: %w", err)
	}
	resp := schemas.Reception{id, createdTime, pvzId, "closed", []schemas.Product{}}
	return &resp, nil
}
