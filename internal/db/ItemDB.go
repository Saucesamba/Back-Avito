package db

import (
	"Backend_trainee_assigment_2025/internal/schemas"
	"fmt"
	"strconv"
)

func CreateProduct(db AvitoDB, typ, pvzId string) (*schemas.Product, error) {
	query := "with open_rec as (select id from receptions where pvz_id = $1 and status = 'in_progress' order by date_time desc limit 1) insert into products (product_type, reception_id) values ($2, open_rec.id) returning *"
	var resp schemas.Product
	intId, _ := strconv.Atoi(pvzId)
	err := db.QueryRow(query, intId, typ).Scan(&resp.Id, &resp.DateTime, &resp.ReceptionId, &resp.Type)
	if err != nil {
		return &schemas.Product{}, fmt.Errorf("Error creating product: %v", err)
	}
	return &resp, nil
}

func DeleteProduct(db AvitoDB, pvzId string) error {
	intId, _ := strconv.Atoi(pvzId)
	query := "with last_rec as(select id from receptions where pvz_id = $1 order by date_time desc limit 1), last_pr as (select p.id from products p join last_rec lr on p.reception_id = lr.id order by p.date_time desc limit 1) delete from products where id in (select id from last_pr)"
	_, err := db.Exec(query, intId)
	if err != nil {
		return fmt.Errorf("Error deleting product: %v", err)
	}
	return nil
}
