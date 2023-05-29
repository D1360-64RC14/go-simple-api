package services

import (
	"fmt"
	"net/http"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/d1360-64rc14/simple-api/utils"
	"golang.org/x/crypto/bcrypt"
)

// DefaultUserService implements UserService
var _ interfaces.UserService = (*DefaultUserService)(nil)

type DefaultUserService struct {
	repo interfaces.UserRepository
}

func NewDefaultUserService(userRepository interfaces.UserRepository) interfaces.UserService {
	return &DefaultUserService{
		repo: userRepository,
	}
}

func (s DefaultUserService) CreateUser(user dtos.UserWithPassword) (*dtos.IdentifiedUser, *utils.ErrorCode) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return nil, utils.NewErrorCode(http.StatusBadRequest, err)
	}

	userHash := &dtos.UserWithHash{
		UserModel: user.UserModel,
		Hash:      string(hash),
	}

	return s.repo.CreateUser(userHash)
}

func (s DefaultUserService) SelectUserFromId(id int) (*dtos.IdentifiedUser, *utils.ErrorCode) {
	return s.repo.SelectUserFromId(id)
}

func (s DefaultUserService) SelectUserHashFromId(id int) (string, *utils.ErrorCode) {
	return s.repo.SelectUserHashFromId(id)
}

func (s DefaultUserService) SelectCompleteUserFromId(id int) (*dtos.IdentifiedUserWithHash, *utils.ErrorCode) {
	return s.repo.SelectCompleteUserFromId(id)
}

func (s DefaultUserService) SelectAllUsers() ([]*dtos.IdentifiedUser, *utils.ErrorCode) {
	return s.repo.SelectAllUsers()
}

func (s DefaultUserService) RemoveUser(id int) *utils.ErrorCode {
	userExist, err := s.repo.UserExist(id)
	if err != nil {
		return err
	}

	if !userExist {
		return utils.NewErrorCodeString(
			http.StatusNotFound,
			fmt.Sprintf("User ID %d doesn't exist", id),
		)
	}

	return s.repo.RemoveUser(id)
}
