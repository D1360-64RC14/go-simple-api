package mocks

import (
	"errors"
	"net/http"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/d1360-64rc14/simple-api/utils"
)

// MockUserRepository implements interfaces.UserRepository
var _ interfaces.UserRepository = (*MockedUserRepository)(nil)

type MockedUserRepository struct {
	IdCounter int
	Closed    bool
	Users     []*dtos.IdentifiedUserWithHash
}

func NewMockedUserRepository() *MockedUserRepository {
	return &MockedUserRepository{
		IdCounter: 0,
		Closed:    false,
		Users:     make([]*dtos.IdentifiedUserWithHash, 0, 5),
	}
}

func NewMockedUserRepositoryWith(idCounter int, users []*dtos.IdentifiedUserWithHash) *MockedUserRepository {
	return &MockedUserRepository{
		IdCounter: idCounter,
		Closed:    false,
		Users:     users,
	}
}

func (r *MockedUserRepository) Close() error {
	if r.Closed {
		return errors.New("already closed")
	}
	r.Closed = true

	return nil
}

func (r *MockedUserRepository) CreateUser(user *dtos.UserWithHash) (*dtos.IdentifiedUser, *utils.ErrorCode) {
	if r.Closed {
		return nil, utils.NewErrorCodeString(http.StatusInternalServerError, "repository closed")
	}

	for _, u := range r.Users {
		if u.Email == user.Email {
			return nil, utils.NewErrorCodeString(http.StatusConflict, "email already exist")
		}
	}

	newUser := &dtos.IdentifiedUserWithHash{
		Hash: user.Hash,
		IdentifiedUser: dtos.IdentifiedUser{
			ID:        r.IdCounter,
			UserModel: user.UserModel,
		},
	}

	r.IdCounter++

	r.Users = append(r.Users, newUser)

	return &newUser.IdentifiedUser, nil
}

func (r *MockedUserRepository) RemoveUser(id int) *utils.ErrorCode {
	if r.Closed {
		return utils.NewErrorCodeString(http.StatusInternalServerError, "repository closed")
	}

	indexToRemove := -1

	for i, user := range r.Users {
		if user.ID == id {
			indexToRemove = i
			break
		}
	}

	if indexToRemove == -1 {
		return utils.NewErrorCodeString(http.StatusBadRequest, "id not found")
	}

	if len(r.Users) == 1 {
		r.Users = r.Users[:0]
		return nil
	}

	// Order is not important: https://stackoverflow.com/a/37335777
	r.Users[indexToRemove] = r.Users[len(r.Users)-1]
	r.Users = r.Users[:len(r.Users)-1]

	return nil
}

func (r MockedUserRepository) SelectAllUsers() ([]*dtos.IdentifiedUser, *utils.ErrorCode) {
	if r.Closed {
		return nil, utils.NewErrorCodeString(http.StatusInternalServerError, "repository closed")
	}

	users := make([]*dtos.IdentifiedUser, len(r.Users))

	for i, user := range r.Users {
		users[i] = &user.IdentifiedUser
	}

	return users, nil
}

func (r MockedUserRepository) SelectCompleteUserFromId(id int) (*dtos.IdentifiedUserWithHash, *utils.ErrorCode) {
	if r.Closed {
		return nil, utils.NewErrorCodeString(http.StatusInternalServerError, "repository closed")
	}

	for _, user := range r.Users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, utils.NewErrorCodeString(http.StatusBadRequest, "id not found")
}

func (r MockedUserRepository) SelectUserFromEmail(email string) (*dtos.IdentifiedUser, *utils.ErrorCode) {
	if r.Closed {
		return nil, utils.NewErrorCodeString(http.StatusInternalServerError, "repository closed")
	}

	for _, user := range r.Users {
		if user.Email == email {
			return &user.IdentifiedUser, nil
		}
	}

	return nil, utils.NewErrorCodeString(http.StatusBadRequest, "id not found")
}

func (r MockedUserRepository) SelectUserFromId(id int) (*dtos.IdentifiedUser, *utils.ErrorCode) {
	if r.Closed {
		return nil, utils.NewErrorCodeString(http.StatusInternalServerError, "repository closed")
	}

	for _, user := range r.Users {
		if user.ID == id {
			return &user.IdentifiedUser, nil
		}
	}

	return nil, utils.NewErrorCodeString(http.StatusBadRequest, "id not found")
}

func (r MockedUserRepository) SelectUserHashFromId(id int) (string, *utils.ErrorCode) {
	if r.Closed {
		return "", utils.NewErrorCodeString(http.StatusInternalServerError, "repository closed")
	}

	for _, user := range r.Users {
		if user.ID == id {
			return user.Hash, nil
		}
	}

	return "", utils.NewErrorCodeString(http.StatusBadRequest, "id not found")
}

func (r *MockedUserRepository) UpdateUsername(id int, newUsername string) *utils.ErrorCode {
	if r.Closed {
		return utils.NewErrorCodeString(http.StatusInternalServerError, "repository closed")
	}

	for _, user := range r.Users {
		if user.ID == id {
			user.UserName = newUsername
			return nil
		}
	}

	return utils.NewErrorCodeString(http.StatusBadRequest, "id not found")
}

func (r MockedUserRepository) UserExist(id int) (bool, *utils.ErrorCode) {
	if r.Closed {
		return false, utils.NewErrorCodeString(http.StatusInternalServerError, "repository closed")
	}

	for _, user := range r.Users {
		if user.ID == id {
			return true, nil
		}
	}

	return false, nil
}
