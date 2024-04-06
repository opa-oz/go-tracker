package api

import (
	"errors"
	"strconv"
	"tracker/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// @BasePath /api

// Increase godoc
// @Summary increase by name
// @Schemes
// @Param name path string true "Name"
// @Param max query int true "Max" default(0)
// @Param step query int true "Step" default(1)
// @Description Increase current index of `name`
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} utils.UnisponseBA.response
// @Router /increase/{name} [get]
// @Router /i/{name} [get]
func Increase(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")
	maxRaw := c.DefaultQuery("max", "0") // 0 - unlimited
	step := c.DefaultQuery("step", "1")

	maxValue, err := strconv.ParseInt(maxRaw, 10, 64)
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

	nextValue := val + stepValue
	if maxValue > 0 && nextValue >= maxValue {
		nextValue = maxValue
	}

	err = rdb.Set(ctx, name, nextValue, 0).Err()
	if err != nil {
		utils.Unirror(c, err)
		return
	}

	utils.UnisponseBA(c, val, nextValue)
}
