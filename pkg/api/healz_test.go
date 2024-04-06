package api_test

import (
	"github.com/gin-gonic/gin"
	"github.com/gkampitakis/go-snaps/snaps"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"tracker/pkg/api"
)

func TestHealzRoute(t *testing.T) {
	r := gin.Default()

	r.GET("/healz", api.Healz)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/healz", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	snaps.MatchSnapshot(t, w.Body.String())
}
