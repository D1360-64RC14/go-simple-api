package v1

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/mocks"
	"github.com/d1360-64rc14/simple-api/models"
	"github.com/d1360-64rc14/simple-api/services"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

var settings = config.Settings{
	Api: config.Api{
		Protocol: "http",
		BaseUrl:  "testing",
	},
	Auth: config.Auth{
		Base64TokenSeed: "IXRoZXF1aWNrZm94anVtcHNvdmVydGhlbGF6eWRvZyE",
		BCryptCost:      12,
	},
}

var authenticator = mocks.NewMockedAuthenticator()

func TestListAllUsers(t *testing.T) {
	testCases := []struct {
		usersToAdd []dtos.UserWithHash
		respBody   string
	}{
		{
			[]dtos.UserWithHash{},
			"[]",
		},
		{
			[]dtos.UserWithHash{
				{UserModel: models.UserModel{UserName: "diego", Email: "diego@mail.com"}, Hash: "fb78ed1e-a121-542f-a68d-fcd21ffe83c5"},
			},
			`[{"id":0,"username":"diego","email":"diego@mail.com"}]`,
		},
		{
			[]dtos.UserWithHash{
				{UserModel: models.UserModel{UserName: "diego", Email: "diego@mail.com"}, Hash: "fb78ed1e-a121-542f-a68d-fcd21ffe83c5"},
				{UserModel: models.UserModel{UserName: "alex", Email: "alex@mail.com"}, Hash: "d296ee89-edba-58c6-8745-d45557cefb90"},
			},
			`[{"id":0,"username":"diego","email":"diego@mail.com"},{"id":1,"username":"alex","email":"alex@mail.com"}]`,
		},
	}

	for i, _case := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			userRepo := mocks.NewMockedUserRepository()

			for _, user := range _case.usersToAdd {
				if _, err := userRepo.CreateUser(&user); err != nil {
					t.Fatalf("userRepo.CreateUser returned an error: (%d) %s", err.Code(), err.Error())
				}
			}

			userService := services.NewDefaultUserService(userRepo, authenticator, &settings)

			controller := NewDefaultUserController(userService, userRepo, &settings)

			engine := gin.New()
			controller.AttachTo(&engine.RouterGroup)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/users", nil)

			engine.ServeHTTP(rec, req)
			body := rec.Body.String()

			if body != _case.respBody {
				t.Errorf("Returned body should be '%s', got '%s'", _case.respBody, body)
			}
		})
	}
}
