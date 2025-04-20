package schemas

import (
	"github.com/google/uuid"
	"time"
)

type Token struct {
	Token string `json:"token"`
}

type User struct {
	Id    uuid.UUID `json:"uuid,omitempty"`
	Email string    `json:"email" validate:"required"`
	Role  string    `json:"role" validate:"required,oneof=employee moderator"`
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
	Id               uuid.UUID `json:"uuid,omitempty"`
	RegistrationDate time.Time `json:"date-time,omitempty"`
	City             string    `json:"city," validate:"required, oneof = Москва Санкт-Петербург Казань"`
}
type Reception struct {
	Id       uuid.UUID `json:"uuid,omitempty"`
	DateTime time.Time `json:"date-time" validate:"required"`
	PVZId    uuid.UUID `json:"pvzuuid" validate:"required"`
	Status   string    `json:"status" validate:"required oneof = in_progress close"`
	Products []Product `json:"products"`
}

type Product struct {
	Id          uuid.UUID `json:"uuid,omitempty"`
	DateTime    time.Time `json:"date-time,omitempty"`
	Type        string    `json:"type" validate:"required,oneof = электроника одежда обувь"`
	ReceptionId uuid.UUID `json:"receptionid" validate:"required"`
}

type Error struct {
	Message string `json:"message" validate:"required"`
}

type PVZWithReceptionsAndProducts struct {
	PVZ        PVZ                     `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions"`
}

type ReceptionWithProducts struct {
	Reception Reception `json:"reception"`
	Products  []Product `json:"products"`
}
