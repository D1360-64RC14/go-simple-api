package dtos

import "github.com/d1360-64rc14/simple-api/models"

type UserWithPassword struct {
	models.UserModel
	Password string `json:"password"`
}
