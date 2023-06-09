package interfaces

import (
	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/utils"
)

type UserRepository interface {
	Close() error
	CreateUser(user *dtos.UserWithHash) (*dtos.IdentifiedUser, *utils.ErrorCode)
	SelectUserFromEmail(email string) (*dtos.IdentifiedUser, *utils.ErrorCode)
	SelectUserFromId(id int) (*dtos.IdentifiedUser, *utils.ErrorCode)
	SelectUserHashFromId(id int) (string, *utils.ErrorCode)
	SelectCompleteUserFromId(id int) (*dtos.IdentifiedUserWithHash, *utils.ErrorCode)
	SelectAllUsers() ([]*dtos.IdentifiedUser, *utils.ErrorCode)
	RemoveUser(id int) *utils.ErrorCode
	UserExist(id int) (bool, *utils.ErrorCode)
	UpdateUsername(id int, newUsername string) *utils.ErrorCode
}
