package db

import (
	"Backend_trainee_assigment_2025/internal/config"
	"Backend_trainee_assigment_2025/internal/schemas"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Интерфейс для того чтоб можно было легко заменить БД при необходимости
type Database interface {
	NewDB(*config.DBConfig) (Database, error)
	GetUser(ctx context.Context, user schemas.UserLogin) (*schemas.User, error)
	CreateUser(ctx context.Context, user *schemas.UserReg) (*schemas.User, error)
	OpenPVZ(ctx context.Context, city string) (*schemas.PVZ, error)
	GetPVZ(ctx context.Context, startTime, endTime string, offset, limit int) ([]schemas.PVZWithReceptionsAndProducts, error)
	OpenRec(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error)
	CloseLastRec(ctx context.Context, pvzId uuid.UUID) (*schemas.Reception, error)
	CreateProduct(ctx context.Context, typ string, pvzId uuid.UUID) (*schemas.Product, error)
	DeleteProduct(ctx context.Context, pvzId uuid.UUID) error
	GetProduct(id uuid.UUID) ([]schemas.Product, error)
}

// Конкретная БД - PostgreSQL
type AvitoDB struct {
	*sql.DB
}

// реализация метода создания подключения к БД
func (a *AvitoDB) NewDB(cfg *config.DBConfig) (Database, error) {
	source := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)
	dbConn, err := sql.Open("postgres", source)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}
	return &AvitoDB{dbConn}, err
}
