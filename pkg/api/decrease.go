package api

import (
	"errors"
	"strconv"
	"tracker/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// @BasePath /api

// Decrease godoc
// @Summary decrease by name
// @Schemes
// @Param name path string true "Name"
// @Param min query int true "Min" default(0)
// @Param step query int true "Step" default(1)
// @Description Decrease current index of `name`
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} utils.UnisponseBA.response
// @Router /decrease/{name} [get]
// @Router /d/{name} [get]
func Decrease(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")
	minRaw := c.DefaultQuery("min", "0") // -1 - unlimited
	step := c.DefaultQuery("step", "1")

	minValue, err := strconv.ParseInt(minRaw, 10, 64)
	if err != nil {
		utils.Unirror(c, err)
		return
	}

	stepValue, err := strconv.ParseInt(step, 10, 64)
	if err != nil {
		utils.Unirror(c, err)
		return
	}

	rdb, err := utils.GetRedis(c)

	if err != nil {
		utils.Unirror(c, err)
		return
	}

	val, err := rdb.Get(ctx, name).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		utils.Unirror(c, err)
		return
	}
	if errors.Is(err, redis.Nil) {
		// key does not exist
		val = 0
	}

	nextValue := val - stepValue
	if minValue != -1 && nextValue < minValue {
		nextValue = minValue
	}

	err = rdb.Set(ctx, name, nextValue, 0).Err()
	if err != nil {
		utils.Unirror(c, err)
		return
	}

	utils.UnisponseBA(c, val, nextValue)
}
