package validate

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/mocks"
	"github.com/d1360-64rc14/simple-api/models"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func paramToKey(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, _ := strconv.Atoi(idStr)
	ctx.Set("id", id)

	ctx.Next()
}

func TestUserIdExist(t *testing.T) {
	testCases := []struct {
		userId   string
		respCode int
		respBody string
	}{
		{"-1", http.StatusNotFound, ""},
		{"0", http.StatusOK, "0"},
		{"1", http.StatusNotFound, ""},
		{"2", http.StatusOK, "2"},
		{"5", http.StatusNotFound, ""},
		{"10", http.StatusNotFound, ""},
	}

	userRepo := mocks.NewMockedUserRepository()

	userRepo.CreateUser(&dtos.UserWithHash{
		UserModel: models.UserModel{UserName: "Diego", Email: "diego@mail.com"},
		Hash:      "fb78ed1e-a121-542f-a68d-fcd21ffe83c5",
	})
	userRepo.CreateUser(&dtos.UserWithHash{
		UserModel: models.UserModel{UserName: "Alex", Email: "alex@mail.com"},
		Hash:      "d296ee89-edba-58c6-8745-d45557cefb90",
	})
	userRepo.RemoveUser(1)
	userRepo.CreateUser(&dtos.UserWithHash{
		UserModel: models.UserModel{UserName: "R2D2", Email: "r2d2@mail.com"},
		Hash:      "621d7d27-1730-59fc-b989-8c4de8d34cd9",
	})

	engine := gin.New()

	engine.GET("/:id", paramToKey, UserIdExist(userRepo), func(ctx *gin.Context) {
		id := ctx.GetInt("id")
		ctx.String(http.StatusOK, fmt.Sprint(id))
	})

	for i, _case := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/"+_case.userId, nil)

			engine.ServeHTTP(rec, req)
			body := rec.Body.String()

			if rec.Code != _case.respCode {
				t.Errorf("Code should be '%d', got '%d'", _case.respCode, rec.Code)
			}
			if body != _case.respBody {
				t.Errorf("Returned body should be '%s', got '%s'", _case.respBody, body)
			}
		})
	}
}
