package api_test

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"net/http"
	"net/http/httptest"
	"testing"
	"tracker/pkg/api"
	"tracker/pkg/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/go-redis/redismock/v9"
)

func TestDecreaseRoute(t *testing.T) {
	r := gin.Default()

	rdb, mock := redismock.NewClientMock()

	r.Use(middlewares.RedisMiddleware(rdb))
	r.GET("/decrease/:name", api.Decrease)
	r.GET("/d/:name", api.Decrease)

	t.Run("if keys is not set, decrease returns 0", func(t *testing.T) {
		mock.ExpectGet("test").RedisNil()
		mock.ExpectSet("test", int64(0), 0).SetVal("0")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/decrease/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"previous":0,"value":0}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)

		mock.ExpectGet("test").SetVal("2")
		mock.ExpectSet("test", int64(1), 0).SetVal("1")

		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, "/d/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"previous":2,"value":1}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("if key is set, get returns value", func(t *testing.T) {
		value := "100"
		mock.ExpectGet("test").SetVal(value)
		mock.ExpectSet("test", int64(92), 0).SetVal("92")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/d/test?step=8", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"previous":100,"value":92}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("if min is set, get returns less or equal", func(t *testing.T) {
		value := "100"
		mock.ExpectGet("test").SetVal(value)
		mock.ExpectSet("test", int64(50), 0).SetVal("50")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/decrease/test?min=50&step=100", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"previous":100,"value":50}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
