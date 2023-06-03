package validate

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestQueryHave(t *testing.T) {
	engine := gin.New()

	engine.GET("/foo", QueryHave("foo"), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ctx.Query("foo"))
	})
	engine.GET("/foobar", QueryHave("foo", "bar"), func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ctx.Query("foo")+ctx.Query("bar"))
	})

	t.Run("WithFoo", testQueryHave_WithFoo(engine))
	t.Run("WithFooBar", testQueryHave_WithFooBar(engine))
	t.Run("WithoutFoo", testQueryHave_WithoutFoo(engine))
}

func testQueryHave_WithFoo(engine *gin.Engine) func(*testing.T) {
	return func(t *testing.T) {
		validQueries := []struct {
			query string
			body  string
		}{
			{query: "foo=bar", body: "bar"},
			{query: "foo=", body: ""},
			{query: "foo=baz&d=g", body: "baz"},
			{query: "foo=done&d=g&k=v", body: "done"},
		}

		for _, validQuery := range validQueries {
			resultingCode := http.StatusOK
			resultingBody := validQuery.body

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/foo?"+validQuery.query, nil)

			engine.ServeHTTP(rec, req)
			body := rec.Body.String()

			if rec.Code != resultingCode {
				t.Errorf("Code should be '%d', got '%d'", resultingCode, rec.Code)
			}
			if body != resultingBody {
				t.Errorf("Returned body should be '%s', got '%s'", resultingBody, body)
			}
		}
	}
}

func testQueryHave_WithFooBar(engine *gin.Engine) func(*testing.T) {
	return func(t *testing.T) {
		validQueries := []struct {
			query string
			body  string
		}{
			{query: "foo=bar&bar=foo", body: "barfoo"},
			{query: "foo=&bar=", body: ""},
			{query: "foo=baz&d=g&bar=none", body: "baznone"},
			{query: "foo=done&d=g&k=v&bar=0", body: "done0"},
		}

		for _, validQuery := range validQueries {
			resultingCode := http.StatusOK
			resultingBody := validQuery.body

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/foobar?"+validQuery.query, nil)

			engine.ServeHTTP(rec, req)
			body := rec.Body.String()

			if rec.Code != resultingCode {
				t.Errorf("Code should be '%d', got '%d'", resultingCode, rec.Code)
			}
			if body != resultingBody {
				t.Errorf("Returned body should be '%s', got '%s'", resultingBody, body)
			}
		}
	}
}

func testQueryHave_WithoutFoo(engine *gin.Engine) func(*testing.T) {
	return func(t *testing.T) {
		invalidQueries := []string{"d=g", "fo=o", "k=v", "k=v&fo=o", "baz=foo"}

		for _, invalidQuery := range invalidQueries {
			resultingCode := http.StatusBadRequest
			resultingBody := "Query should have the following elements: foo"

			rec := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/foo?"+invalidQuery, nil)

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
