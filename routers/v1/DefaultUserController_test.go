package v1

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/d1360-64rc14/simple-api/config"
	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/mocks"
	"github.com/d1360-64rc14/simple-api/models"
	"github.com/d1360-64rc14/simple-api/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		BCryptCost:      bcrypt.MinCost,
	},
}

var authenticator = mocks.NewMockedAuthenticator()

func TestListAllUsers(t *testing.T) {
	testCases := []struct {
		usersToAdd []*dtos.IdentifiedUserWithHash
		idCount    int
		respBody   string
	}{
		{
			[]*dtos.IdentifiedUserWithHash{},
			0,
			"[]",
		},
		{
			[]*dtos.IdentifiedUserWithHash{
				{IdentifiedUser: dtos.IdentifiedUser{ID: 0, UserModel: models.UserModel{UserName: "diego", Email: "diego@mail.com"}}, Hash: "fb78ed1e-a121-542f-a68d-fcd21ffe83c5"},
			},
			1,
			`[{"id":0,"username":"diego","email":"diego@mail.com"}]`,
		},
		{
			[]*dtos.IdentifiedUserWithHash{
				{IdentifiedUser: dtos.IdentifiedUser{ID: 0, UserModel: models.UserModel{UserName: "diego", Email: "diego@mail.com"}}, Hash: "fb78ed1e-a121-542f-a68d-fcd21ffe83c5"},
				{IdentifiedUser: dtos.IdentifiedUser{ID: 5, UserModel: models.UserModel{UserName: "alex", Email: "alex@mail.com"}}, Hash: "d296ee89-edba-58c6-8745-d45557cefb90"},
			},
			6,
			`[{"id":0,"username":"diego","email":"diego@mail.com"},{"id":5,"username":"alex","email":"alex@mail.com"}]`,
		},
		{
			[]*dtos.IdentifiedUserWithHash{
				{IdentifiedUser: dtos.IdentifiedUser{ID: 0, UserModel: models.UserModel{UserName: "diego", Email: "diego@mail.com"}}, Hash: "fb78ed1e-a121-542f-a68d-fcd21ffe83c5"},
				{IdentifiedUser: dtos.IdentifiedUser{ID: 5, UserModel: models.UserModel{UserName: "alex", Email: "alex@mail.com"}}, Hash: "d296ee89-edba-58c6-8745-d45557cefb90"},
				{IdentifiedUser: dtos.IdentifiedUser{ID: 2, UserModel: models.UserModel{UserName: "R2D2", Email: "r2d2@mail.com"}}, Hash: "621d7d27-1730-59fc-b989-8c4de8d34cd9"},
			},
			6,
			`[{"id":0,"username":"diego","email":"diego@mail.com"},{"id":5,"username":"alex","email":"alex@mail.com"},{"id":2,"username":"R2D2","email":"r2d2@mail.com"}]`,
		},
	}

	for i, _case := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			userRepo := mocks.NewMockedUserRepositoryWith(_case.idCount, _case.usersToAdd)
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

func TestCreateSingleUser(t *testing.T) {
	testCases := []struct {
		usersInDB []*dtos.IdentifiedUserWithHash
		idCount   int
		postBody  string
		respCode  int
		respBody  string
		location  string
	}{
		{
			[]*dtos.IdentifiedUserWithHash{},
			0,
			`{"username": "diego", "email": "diego@mail.com", "password": "d1360@m41l.c0m"}`,
			http.StatusCreated,
			`{"id":0,"username":"diego","email":"diego@mail.com"}`,
			fmt.Sprintf("%s://%s/user/0", settings.Api.Protocol, settings.Api.BaseUrl),
		},
		{
			[]*dtos.IdentifiedUserWithHash{
				{IdentifiedUser: dtos.IdentifiedUser{ID: 0, UserModel: models.UserModel{UserName: "diego", Email: "diego@mail.com"}}, Hash: "fb78ed1e-a121-542f-a68d-fcd21ffe83c5"},
			},
			1,
			`{"username": "alex", "email": "alex@mail.com", "password": "4l3x@m41l.c0m"}`,
			http.StatusCreated,
			`{"id":1,"username":"alex","email":"alex@mail.com"}`,
			fmt.Sprintf("%s://%s/user/1", settings.Api.Protocol, settings.Api.BaseUrl),
		},
		{
			[]*dtos.IdentifiedUserWithHash{
				{IdentifiedUser: dtos.IdentifiedUser{ID: 0, UserModel: models.UserModel{UserName: "diego", Email: "diego@mail.com"}}, Hash: "fb78ed1e-a121-542f-a68d-fcd21ffe83c5"},
				{IdentifiedUser: dtos.IdentifiedUser{ID: 5, UserModel: models.UserModel{UserName: "alex", Email: "alex@mail.com"}}, Hash: "d296ee89-edba-58c6-8745-d45557cefb90"},
			},
			6,
			`{"username": "R2D2", "email": "r2d2@mail.com", "password": "r2d2@m41l.c0m"}`,
			http.StatusCreated,
			`{"id":6,"username":"R2D2","email":"r2d2@mail.com"}`,
			fmt.Sprintf("%s://%s/user/6", settings.Api.Protocol, settings.Api.BaseUrl),
		},
		{
			[]*dtos.IdentifiedUserWithHash{
				{IdentifiedUser: dtos.IdentifiedUser{ID: 0, UserModel: models.UserModel{UserName: "diego", Email: "diego@mail.com"}}, Hash: "fb78ed1e-a121-542f-a68d-fcd21ffe83c5"},
				{IdentifiedUser: dtos.IdentifiedUser{ID: 5, UserModel: models.UserModel{UserName: "alex", Email: "alex@mail.com"}}, Hash: "d296ee89-edba-58c6-8745-d45557cefb90"},
				{IdentifiedUser: dtos.IdentifiedUser{ID: 2, UserModel: models.UserModel{UserName: "R2D2", Email: "r2d2@mail.com"}}, Hash: "621d7d27-1730-59fc-b989-8c4de8d34cd9"},
			},
			6,
			`{"username": "R2D2", "email": "r2d2@mail.com", "password": "r2d2@m41l.c0m"}`,
			http.StatusConflict,
			`{"error":"email already exist"}`,
			"",
		},
	}

	for i, _case := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			userRepo := mocks.NewMockedUserRepositoryWith(_case.idCount, _case.usersInDB)
			userService := services.NewDefaultUserService(userRepo, authenticator, &settings)

			controller := NewDefaultUserController(userService, userRepo, &settings)

			engine := gin.New()
			controller.AttachTo(&engine.RouterGroup)

			rec := httptest.NewRecorder()
			reqBody := bytes.NewBuffer([]byte(_case.postBody))
			req := httptest.NewRequest("POST", "/user", reqBody)

			engine.ServeHTTP(rec, req)
			body := rec.Body.String()
			location := rec.Header().Get("Location")

			if rec.Code != _case.respCode {
				t.Errorf("Returned code should be '%d', got '%d'", _case.respCode, rec.Code)
			}

			if body != _case.respBody {
				t.Errorf("Returned body should be '%s', got '%s'", _case.respBody, body)
			}

			if location != _case.location {
				t.Errorf("Returned header Location should be '%s', got '%s'", _case.location, location)
			}
		})
	}
}
