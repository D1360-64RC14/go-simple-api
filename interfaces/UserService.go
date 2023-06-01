package interfaces

import (
	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/utils"
)

type UserService interface {
	CreateUser(user dtos.UserWithPassword) (*dtos.IdentifiedUser, *utils.ErrorCode)
	SelectUserFromId(id int) (*dtos.IdentifiedUser, *utils.ErrorCode)
	SelectUserHashFromId(id int) (string, *utils.ErrorCode)
	SelectCompleteUserFromId(id int) (*dtos.IdentifiedUserWithHash, *utils.ErrorCode)
	SelectAllUsers() ([]*dtos.IdentifiedUser, *utils.ErrorCode)
	RemoveUser(id int) *utils.ErrorCode
	UpdateUser(id int, newUserData *dtos.UserUpdate) *utils.ErrorCode
	AuthenticateUser(email string) (string, *utils.ErrorCode)
}
