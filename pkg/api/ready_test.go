package api_test

import (
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"tracker/pkg/api"

	"github.com/gin-gonic/gin"
)

func TestReadyRoute(t *testing.T) {
	r := gin.Default()

	r.GET("/ready", api.Ready)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ready", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	snaps.MatchSnapshot(t, w.Body.String())
}
