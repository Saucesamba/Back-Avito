package db

import (
	"Backend_trainee_assigment_2025/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

// Интерфейс для того чтоб можно было легко заменить БД при необходимости
type Database interface {
	NewDB(*config.DBConfig) (Database, error)
}

// Конкретная БД - PostgreSQL
type AvitoDB struct {
	*sql.DB
}

// реализация метода создания подключения к БД
func (a *AvitoDB) NewDB(cfg *config.DBConfig) (Database, error) {
	source := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Name)
	dbConn, err := sql.Open("postgres", source)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	if err := dbConn.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}
	return &AvitoDB{dbConn}, err
}
