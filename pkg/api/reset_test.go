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

func TestResetRoute(t *testing.T) {
	r := gin.Default()

	rdb, mock := redismock.NewClientMock()

	r.Use(middlewares.RedisMiddleware(rdb))
	r.GET("/reset/:name", api.Reset)
	r.GET("/r/:name", api.Reset)

	t.Run("if keys is not set, reset to 0", func(t *testing.T) {
		mock.ExpectSet("test", int64(0), 0).SetVal("0")
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/reset/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"value":0}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("if key is set, reset to 0", func(t *testing.T) {
		mock.ExpectSet("test", int64(0), 0).SetVal("0")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/reset/test", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"value":0}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("if init is set, reset to init value", func(t *testing.T) {
		mock.ExpectSet("test", int64(50), 0).SetVal("0")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/r/test?init=50", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, `{"value":50}`, w.Body.String())
		snaps.MatchSnapshot(t, w.Body.String())
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
