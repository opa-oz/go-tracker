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

func TestGetRoute(t *testing.T) {
	r := gin.Default()

	rdb, mock := redismock.NewClientMock()

	r.Use(middlewares.RedisMiddleware(rdb))
	r.GET("/get/:name", api.Get)
	r.GET("/g/:name", api.Get)

	t.Run("if keys is not set, get returns 0", func(t *testing.T) {
		mock.ExpectGet("test").RedisNil()
		mock.ExpectSet("test", 0, 0).SetVal("0")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/get/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"value":0}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("if key is set, get returns value", func(t *testing.T) {
		value := "100"
		mock.ExpectGet("test").SetVal(value)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/g/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"value":100}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("if max is set, get returns less or equal", func(t *testing.T) {
		value := "100"
		mock.ExpectGet("test").SetVal(value)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/get/test?max=50", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"value":50}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
