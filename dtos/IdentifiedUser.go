package dtos

import "github.com/d1360-64rc14/simple-api/models"

type IdentifiedUser struct {
	ID int `json:"id"`
	models.UserModel
}
