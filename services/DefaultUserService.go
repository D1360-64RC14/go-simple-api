package services

import (
	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/repositories"
)

type DefaultUserService struct {
	repo repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &DefaultUserService{
		repo: userRepository,
	}
}

func (s DefaultUserService) CreateUser(username, email, password string) (*dtos.IdentifiedUser, error) {
	return s.repo.CreateUser(username, email, password)
}

func (s DefaultUserService) SelectUserFromId(id int) (*dtos.IdentifiedUser, error) {
	return s.repo.SelectUserFromId(id)
}

func (s DefaultUserService) SelectUserHashFromId(id int) (string, error) {
	return s.repo.SelectUserHashFromId(id)
}

func (s DefaultUserService) SelectCompleteUserFromId(id int) (*dtos.IdentifiedUserWithHash, error) {
	return s.repo.SelectCompleteUserFromId(id)
}

func (s DefaultUserService) SelectAllUsers() ([]*dtos.IdentifiedUser, error) {
	return s.repo.SelectAllUsers()
}
