package dtos

import "github.com/d1360-64rc14/simple-api/models"

type IdentifiedUser struct {
	ID int `json:"id"`
	models.UserModel
}

type UserWithPassword struct {
	models.UserModel
	Password string `json:"password"`
}

type UserWithHash struct {
	models.UserModel
	Hash string `json:"hash"`
}

type IdentifiedUserWithHash struct {
	IdentifiedUser
	Hash string `json:"hash"`
}
