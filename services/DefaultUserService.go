package services

import (
	"fmt"
	"net/http"

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/d1360-64rc14/simple-api/utils"
	"golang.org/x/crypto/bcrypt"
)

// DefaultUserService implements UserService
var _ interfaces.UserService = (*DefaultUserService)(nil)

type DefaultUserService struct {
	repo     interfaces.UserRepository
	auth     interfaces.Authenticator
	settings *config.Settings
}

func NewDefaultUserService(
	userRepository interfaces.UserRepository,
	authenticator interfaces.Authenticator,
	settings *config.Settings,
) interfaces.UserService {
	return &DefaultUserService{
		repo:     userRepository,
		auth:     authenticator,
		settings: settings,
	}
}

func (s DefaultUserService) CreateUser(user dtos.UserWithPassword) (*dtos.IdentifiedUser, *utils.ErrorCode) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), s.settings.Auth.BCryptCost)
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

func (s DefaultUserService) UpdateUser(id int, newUserData *dtos.UserUpdate) *utils.ErrorCode {
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

	if newUserData.UserName != "" {
		err := s.repo.UpdateUsername(id, newUserData.UserName)
		if err != nil {
			return err
		}
	}

	return nil
}

// AuthenticateUser returns the JWT token as result of the authentication
func (s DefaultUserService) AuthenticateUser(email string) (string, *utils.ErrorCode) {
	user, errC := s.repo.SelectUserFromEmail(email)
	if errC != nil {
		return "", errC
	}

	jwtToken, err := s.auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	return jwtToken, nil
}
