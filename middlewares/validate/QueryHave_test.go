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

func TestQueryHave_WithOneQuery(t *testing.T) {
	testCases := []struct {
		respCode int
		query    string
		respBody string
	}{
		{http.StatusOK, "foo=bar", "bar"},
		{http.StatusOK, "foo=", ""},
		{http.StatusOK, "foo=baz&d=g", "baz"},
		{http.StatusOK, "foo=done&d=g&k=v", "done"},
		{http.StatusBadRequest, "d=g", "{\"error\":\"Query should have the following elements: foo\"}"},
		{http.StatusBadRequest, "fo=o", "{\"error\":\"Query should have the following elements: foo\"}"},
		{http.StatusBadRequest, "k=v", "{\"error\":\"Query should have the following elements: foo\"}"},
		{http.StatusBadRequest, "k=v&fo=o", "{\"error\":\"Query should have the following elements: foo\"}"},
		{http.StatusBadRequest, "baz=foo", "{\"error\":\"Query should have the following elements: foo\"}"},
	}

	engine := gin.New()

	engine.GET("/", QueryHave("foo"), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ctx.Query("foo"))
	})

	for i, _case := range testCases {
		t.Run(fmt.Sprintf("token_%d", i), func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/?"+_case.query, nil)

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

func TestQueryHave_WithTwoQueries(t *testing.T) {
	testCases := []struct {
		respCode int
		query    string
		respBody string
	}{
		{http.StatusOK, "foo=bar&bar=foo", "barfoo"},
		{http.StatusOK, "foo=&bar=", ""},
		{http.StatusOK, "foo=baz&d=g&bar=none", "baznone"},
		{http.StatusOK, "foo=done&d=g&k=v&bar=0", "done0"},
	}

	engine := gin.New()

	engine.GET("/", QueryHave("foo", "bar"), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ctx.Query("foo")+ctx.Query("bar"))
	})

	for i, _case := range testCases {
		t.Run(fmt.Sprintf("token_%d", i), func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/?"+_case.query, nil)

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
