package services

import (
	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/interfaces"
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
