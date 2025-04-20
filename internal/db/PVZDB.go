package db

import (
	"Backend_trainee_assigment_2025/internal/schemas"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
)

func (a *AvitoDB) OpenPVZ(ctx context.Context, city string) (*schemas.PVZ, error) {
	query := "INSERT INTO pvzs (city) VALUES ($1) returning id, registration_date"
	var id uuid.UUID
	var createdTime time.Time
	err := a.QueryRowContext(ctx, query, city).Scan(&id, &createdTime)
	if err != nil {
		return &schemas.PVZ{}, fmt.Errorf("failed to open PVZ: %w", err)
	}

	if err != nil {
		fmt.Println(err)
	}
	resp := schemas.PVZ{Id: id, RegistrationDate: createdTime, City: city}
	return &resp, nil
}

func (a *AvitoDB) GetPVZ(ctx context.Context, startTime, endTime string, offset, limit int) ([]schemas.PVZWithReceptionsAndProducts, error) {

	first, err := time.Parse("2006-01-02 15:04:05.999999", startTime)
	second, _ := time.Parse("2006-01-02 15:04:05.999999", endTime)

	query := `
	select p.id, p.registration_date, p.city, r.id, 
	r.date_time, r.pvz_id, r.status, pr.id, pr.date_time,
	pr.reception_id, pr.type from pvzs p
	left join receptions r on p.id = r.pvz_id
	left join products pr on r.id = pr.reception_id 
	where pr.date_time 
	between $1 and $2
`
	rows, err := a.QueryContext(ctx, query, first, second)

	if err != nil {
		return []schemas.PVZWithReceptionsAndProducts{}, err
	}
	defer rows.Close()

	pvzMap := make(map[uuid.UUID]*schemas.PVZWithReceptionsAndProducts)
	for rows.Next() {
		var pvzID uuid.UUID
		var pvzRegistrationDate time.Time
		var city string
		var receptionID sql.NullString
		var receptionDateTime sql.NullTime
		var status sql.NullString
		var productID sql.NullString
		var productDateTime sql.NullTime
		var productType sql.NullString
		var blank sql.NullString
		err = rows.Scan(&pvzID, &pvzRegistrationDate, &city, &receptionID, &receptionDateTime, &blank, &status, &productID, &productDateTime, &blank, &productType)
		if err != nil {
			return nil, fmt.Errorf("failed at scan %v", err)
		}
		//Get PVZ if not exists
		if _, ok := pvzMap[pvzID]; !ok {
			pvzMap[pvzID] = &schemas.PVZWithReceptionsAndProducts{
				PVZ: schemas.PVZ{
					Id:               pvzID,
					RegistrationDate: pvzRegistrationDate,
					City:             city,
				},
				Receptions: []schemas.ReceptionWithProducts{},
			}
		}
		//validateReceptions
		if receptionID.Valid {
			receptionUUID, err := uuid.Parse(receptionID.String)
			if err != nil {
				return nil, fmt.Errorf("parsing reception uuid: %w", err)
			}
			var receptionInfo *schemas.ReceptionWithProducts
			receptionFound := false
			for i := range pvzMap[pvzID].Receptions {
				if pvzMap[pvzID].Receptions[i].Reception.Id == receptionUUID {
					receptionInfo = &pvzMap[pvzID].Receptions[i]
					receptionFound = true
					break
				}
			}
			if !receptionFound {
				receptionInfo = &schemas.ReceptionWithProducts{
					Reception: schemas.Reception{
						Id:       receptionUUID,
						DateTime: receptionDateTime.Time,
						Status:   status.String,
					},
					Products: []schemas.Product{},
				}
				pvzMap[pvzID].Receptions = append(pvzMap[pvzID].Receptions, *receptionInfo)
			}

			if productID.Valid {
				productUUID, err := uuid.Parse(productID.String)
				if err != nil {
					log.Printf("Error parsing product UUID: %v", err)
					continue
				}
				product := schemas.Product{
					Id:          productUUID,
					DateTime:    productDateTime.Time,
					ReceptionId: receptionUUID,
					Type:        productType.String,
				}
				receptionInfo.Products = append(receptionInfo.Products, product)
			}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed at rows %w", err)
	}
	var pvzsList []schemas.PVZWithReceptionsAndProducts
	for _, pvzInfo := range pvzMap {
		pvzsList = append(pvzsList, *pvzInfo)
	}

	return pvzsList, nil
}
