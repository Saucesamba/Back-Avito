package db

import (
	"Backend_trainee_assigment_2025/internal/schemas"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (a *AvitoDB) CreateProduct(ctx context.Context, typ string, pvzId uuid.UUID) (*schemas.Product, error) {
	query := `
        WITH open_rec AS (
            SELECT id FROM receptions
            WHERE pvz_id = $1 AND status = 'in_progress'
            ORDER BY date_time DESC
            LIMIT 1
        )
        INSERT INTO products (type, reception_id)
        SELECT $2, id FROM open_rec
        RETURNING id, type, reception_id, date_time
    `

	var resp schemas.Product

	err := a.QueryRowContext(ctx, query, pvzId, typ).Scan(&resp.Id, &resp.DateTime, &resp.ReceptionId, &resp.Type)
	if err != nil {
		return &schemas.Product{}, fmt.Errorf("Error creating product: %v", err)
	}
	return &resp, nil
}

func (a *AvitoDB) GetProduct(id uuid.UUID) ([]schemas.Product, error) {
	query := "select * from products where reception_id=$1"
	var resp []schemas.Product
	rows, err := a.Query(query, id)
	if err != nil {
		return resp, err
	}
	defer rows.Close()
	for rows.Next() {
		var p schemas.Product
		err = rows.Scan(&p.Id, &p.DateTime, &p.Type, &p.ReceptionId)
		if err != nil {
			return resp, err
		}
		resp = append(resp, p)
	}
	return resp, nil
}

func (a *AvitoDB) DeleteProduct(ctx context.Context, pvzId uuid.UUID) error {
	query := `
    WITH last_rec AS (
      SELECT id
      FROM receptions
      WHERE pvz_id = $1 AND status = 'in_progress'
      ORDER BY date_time DESC
      LIMIT 1
    ),
    last_pr AS (
      SELECT p.id
      FROM products p
      JOIN last_rec lr ON p.reception_id = lr.id
      ORDER BY p.date_time DESC
      LIMIT 1
    )
    DELETE FROM products
    WHERE id IN (
      SELECT id
      FROM last_pr
    )
    AND EXISTS (SELECT 1 FROM last_pr);
`
	res, err := a.ExecContext(ctx, query, pvzId)

	if err != nil {
		return fmt.Errorf("Error deleting product: %v", err)
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error getting affected row: %v", err)
	}

	if rowCnt == 0 {
		return fmt.Errorf("no products to delete in the current reception")
	}

	return nil
}
