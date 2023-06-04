package validate

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestPathUserId(t *testing.T) {
	testCases := []struct {
		inputId  string
		respCode int
		respBody string
	}{
		{"42", http.StatusOK, "42"},
		{"0", http.StatusOK, "0"},
		{"-5", http.StatusOK, "-5"},
		{"000", http.StatusOK, "0"},
		{"99", http.StatusOK, "99"},
		{"+100", http.StatusOK, "100"},
		{"1337", http.StatusOK, "1337"},
		{"foo", http.StatusBadRequest, "{\"error\":\"The user ID in the path should be an integer, not 'foo'\"}"},
		{"bar", http.StatusBadRequest, "{\"error\":\"The user ID in the path should be an integer, not 'bar'\"}"},
		{".", http.StatusBadRequest, "{\"error\":\"The user ID in the path should be an integer, not '.'\"}"},
		{",", http.StatusBadRequest, "{\"error\":\"The user ID in the path should be an integer, not ','\"}"},
		{":", http.StatusBadRequest, "{\"error\":\"The user ID in the path should be an integer, not ':'\"}"},
		{":id", http.StatusBadRequest, "{\"error\":\"The user ID in the path should be an integer, not ':id'\"}"},
		{"id", http.StatusBadRequest, "{\"error\":\"The user ID in the path should be an integer, not 'id'\"}"},
	}

	engine := gin.New()

	engine.GET("/:id", PathUserId, func(ctx *gin.Context) {
		id := ctx.GetInt("id")
		ctx.String(http.StatusOK, fmt.Sprint(id))
	})

	for i, _case := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/"+_case.inputId, nil)

			engine.ServeHTTP(rec, req)
			body := rec.Body.String()

			if rec.Code != _case.respCode {
				t.Errorf("Code should be '%d', got '%d'", _case.respCode, rec.Code)
			}
			if body != _case.respBody {
				t.Errorf("Returned id should be '%s', got '%s'", _case.respBody, body)
			}
		})
	}
}
