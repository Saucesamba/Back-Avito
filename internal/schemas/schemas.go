package schemas

import "time"

type Token struct {
	Token string `json:"token"`
}

type User struct {
	Id    string `json:"uuid,omitempty"`
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required,oneof=employee moderator"`
}

type UserReg struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=employee moderator"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type PVZ struct {
	Id               string    `json:"uuid,omitempty"`
	RegistrationDate time.Time `json:"date-time,omitempty"`
	City             string    `json:"city," validate:"required, oneof = Москва Санкт-Петербург Казань"`
}
type Reception struct {
	Id       string    `json:"uuid,omitempty"`
	DateTime time.Time `json:"date-time" validate:"required"`
	PVZId    string    `json:"pvzuuid" validate:"required"`
	Status   string    `json:"status" validate:"required oneof = in_progress close"`
}

type Product struct {
	Id          string `json:"uuid,omitempty"`
	DateTime    string `json:"date-time,omitempty"`
	Type        string `json:"type" validate:"required,oneof = электроника одежда обувь"`
	ReceptionId string `json:"receptionid" validate:"required"`
}

type Error struct {
	Message string `json:"message" validate:"required"`
}

type MegaResponse struct {
	Pvz     PVZ
	Receipt Reception
	Product Product
}
