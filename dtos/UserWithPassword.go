package dtos

import "github.com/d1360-64rc14/simple-api/models"

type UserWithPassword struct {
	models.UserModel
	Password string `json:"password" binding:"required,min=8,max=72,ascii"`
}
