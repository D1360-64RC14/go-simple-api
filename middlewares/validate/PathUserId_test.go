package validate

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPathUserId(t *testing.T) {
	engine := gin.New()

	engine.GET("/:id", PathUserId, func(ctx *gin.Context) {
		id := ctx.GetInt("id")
		ctx.String(http.StatusOK, fmt.Sprint(id))
	})

	t.Run("WithValidIdType", testPathUserId_WithValidIdType(engine))
	t.Run("WithInvalidIdType", testPathUserId_WithInvalidIdType(engine))
}

func testPathUserId_WithValidIdType(engine *gin.Engine) func(*testing.T) {
	return func(t *testing.T) {
		validIds := []string{"42", "0", "-5", "000", "99", "+100", "1337"}

		for _, userId := range validIds {
			intId, _ := strconv.ParseInt(userId, 10, 32)

			resultingCode := http.StatusOK
			resultingBody := fmt.Sprint(intId)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/"+userId, nil)

			engine.ServeHTTP(rec, req)
			body := rec.Body.String()

			if rec.Code != resultingCode {
				t.Errorf("Code should be '%d', got '%d'", resultingCode, rec.Code)
			}
			if body != resultingBody {
				t.Errorf("Returned id should be '%s', got '%s'", resultingBody, body)
			}
		}
	}
}

func testPathUserId_WithInvalidIdType(engine *gin.Engine) func(*testing.T) {
	return func(t *testing.T) {
		invalidIds := []string{"foo", "bar", ".", ",", ":", ":id", "id"}

		for _, userId := range invalidIds {
			resultingCode := http.StatusBadRequest
			resultingBody := fmt.Sprintf("The user ID in the path should be an integer, not '%s'", userId)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/"+userId, nil)

			engine.ServeHTTP(rec, req)
			body := rec.Body.String()

			if rec.Code != resultingCode {
				t.Errorf("Code should be '%d', got '%d'", resultingCode, rec.Code)
			}
			if !strings.Contains(body, resultingBody) {
				t.Errorf("Returned body should have '%s', got '%s'", resultingBody, body)
			}
		}
	}
}
