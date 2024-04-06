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

func TestIncreaseRoute(t *testing.T) {
	r := gin.Default()

	rdb, mock := redismock.NewClientMock()

	r.Use(middlewares.RedisMiddleware(rdb))
	r.GET("/increase/:name", api.Increase)
	r.GET("/i/:name", api.Increase)

	t.Run("if keys is not set, increase returns 1", func(t *testing.T) {
		mock.ExpectGet("test").RedisNil()
		mock.ExpectSet("test", int64(1), 0).SetVal("1")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/increase/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"previous":0,"value":1}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)

		mock.ExpectGet("test").SetVal("1")
		mock.ExpectSet("test", int64(2), 0).SetVal("2")

		w = httptest.NewRecorder()
		req, _ = http.NewRequest(http.MethodGet, "/i/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"previous":1,"value":2}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("if key is set, get returns value", func(t *testing.T) {
		value := "100"
		mock.ExpectGet("test").SetVal(value)
		mock.ExpectSet("test", int64(112), 0).SetVal("112")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/i/test?step=12", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"previous":100,"value":112}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("if max is set, get returns less or equal", func(t *testing.T) {
		value := "30"
		mock.ExpectGet("test").SetVal(value)
		mock.ExpectSet("test", int64(50), 0).SetVal("50")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/increase/test?max=50&step=100", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"previous":30,"value":50}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
