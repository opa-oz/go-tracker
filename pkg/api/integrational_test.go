package api_test

import (
	"fmt"
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

func val(v int64) string {
	return fmt.Sprintf(`{"value":%d}`, v)
}

func val2(b, a int64) string {
	return fmt.Sprintf(`{"previous":%d,"value":%d}`, b, a)
}

type Case struct {
	A  string // action
	P  string // params
	O  int64  // output
	Pr int64  // previous
}

func TestAllEndpoints(t *testing.T) {
	name := "my_key"
	pipeline := []Case{
		{"get", "", 0, 0},
		{"increase", "?step=100", 100, 0},
		{"get", "", 100, 100},
		{"decrease", "?step=10", 90, 100},
		{"increase", "?step=100&max=160", 160, 90},
		{"reset", "", 0, 160},
		{"increase", "?step=10", 10, 0},
		{"decrease", "?step=100", 0, 10},
	}

	r := gin.Default()

	rdb, mock := redismock.NewClientMock()

	r.Use(middlewares.RedisMiddleware(rdb))
	r.GET("/decrease/:name", api.Decrease)
	r.GET("/increase/:name", api.Increase)
	r.GET("/get/:name", api.Get)
	r.GET("/reset/:name", api.Reset)

	mock.ExpectGet(name).RedisNil()
	mock.ExpectSet(name, 0, 0).SetVal("")
	mock.ExpectGet(name).SetVal("0")
	mock.ExpectSet(name, int64(100), 0).SetVal("")
	mock.ExpectGet(name).SetVal("100")
	mock.ExpectGet(name).SetVal("100")
	mock.ExpectSet(name, int64(90), 0).SetVal("")
	mock.ExpectGet(name).SetVal("90")
	mock.ExpectSet(name, int64(160), 0).SetVal("")
	mock.ExpectSet(name, int64(0), 0).SetVal("")
	mock.ExpectGet(name).SetVal("0")
	mock.ExpectSet(name, int64(10), 0).SetVal("")
	mock.ExpectGet(name).SetVal("10")
	mock.ExpectSet(name, int64(0), 0).SetVal("")

	for _, tcase := range pipeline {
		t.Run(fmt.Sprintf("%s%s -> %d", tcase.A, tcase.P, tcase.O), func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/%s/%s%s", tcase.A, name, tcase.P), nil)
			r.ServeHTTP(w, req)

			if tcase.A == "increase" || tcase.A == "decrease" {
				assert.Equal(t, val2(tcase.Pr, tcase.O), w.Body.String())
				snaps.MatchSnapshot(t, w.Body.String())
			} else {
				assert.Equal(t, val(tcase.O), w.Body.String())
				snaps.MatchSnapshot(t, w.Body.String())
			}

			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}
