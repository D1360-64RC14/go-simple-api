package mocks

import (
	"fmt"
	"strings"

	"github.com/d1360-64rc14/simple-api/interfaces"
)

// MockedAuthenticator implements interfaces.Authenticator
var _ interfaces.Authenticator = (*MockedAuthenticator)(nil)

type MockedAuthenticator struct{}

func NewMockedAuthenticator() *MockedAuthenticator {
	return &MockedAuthenticator{}
}

func (*MockedAuthenticator) GenerateToken(id int, email string) (string, error) {
	token := fmt.Sprintf("valid-for(%d)[%s]", id, email)

	return token, nil
}

func (*MockedAuthenticator) IsTokenValid(inputToken string) bool {
	return strings.HasPrefix(inputToken, "valid")
}
