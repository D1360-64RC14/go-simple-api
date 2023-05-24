package repositories

import "github.com/d1360-64rc14/simple-api/dtos"

type UserRepository interface {
	Close() error
	CreateUser(username, email, password string) (*dtos.IdentifiedUser, error)
	SelectUserFromId(id int) (*dtos.IdentifiedUser, error)
	SelectUserHashFromId(id int) (string, error)
	SelectCompleteUserFromId(id int) (*dtos.IdentifiedUserWithHash, error)
	SelectAllUsers() ([]*dtos.IdentifiedUser, error)
}
