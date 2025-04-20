package db

import (
	"Backend_trainee_assigment_2025/internal/schemas"
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (a *AvitoDB) CreateUser(ctx context.Context, user *schemas.UserReg) (*schemas.User, error) {
	query := "INSERT INTO users (email, password, role) values ($1, $2, $3) RETURNING id"
	var id uuid.UUID
	err := a.QueryRowContext(ctx, query, user.Email, user.Password, user.Role).Scan(&id)
	if err != nil {
		return &schemas.User{}, fmt.Errorf("Error creating user: %v", err)
	}
	userResp := schemas.User{Id: id, Email: user.Email, Role: user.Role}
	return &userResp, nil
}

func (a *AvitoDB) GetUser(ctx context.Context, user schemas.UserLogin) (*schemas.User, error) {
	var userResp schemas.User
	query := "SELECT id, email, role FROM users WHERE email=$1 AND password=$2" // id, email, role added
	row := a.QueryRowContext(ctx, query, user.Email, user.Password)
	err := row.Scan(&userResp.Id, &userResp.Email, &userResp.Role) // changed vars to point on
	if err != nil {
		return &schemas.User{}, fmt.Errorf("error getting user: %v", err)
	}
	return &userResp, nil
}
