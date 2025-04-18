package db

import (
	"Backend_trainee_assigment_2025/internal/schemas"
	"fmt"
)

func CreateUser(db AvitoDB, user *schemas.UserReg) (*schemas.User, error) {
	query := "INSERT INTO users (email, password, role) values ($1, $2, $3) RETURNING id"
	var id int
	err := db.QueryRow(query, user.Email, user.Password, user.Role).Scan(&id)
	if err != nil {
		return &schemas.User{}, fmt.Errorf("Error creating user: %v", err)
	}
	userResp := schemas.User{Id: string(id), Email: user.Email, Role: user.Role}
	return &userResp, nil
}

func GetUser(db AvitoDB, user schemas.UserLogin) (*schemas.User, error) {
	var userResp schemas.User
	query := "SELECT * FROM users WHERE email=$1 AND password=$2"
	row := db.QueryRow(query, user.Email, user.Password)
	err := row.Scan(&user.Email, &user.Password)
	if err != nil {
		return &schemas.User{}, fmt.Errorf("Error getting user: %v", err)
	}
	return &userResp, nil
}
