package dtos

import "github.com/d1360-64rc14/simple-api/models"

type UserWithHash struct {
	models.UserModel
	Hash string `json:"hash"`
}
