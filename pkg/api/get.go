package api

import (
	"errors"
	"strconv"
	"tracker/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// @BasePath /api

// Get godoc
// @Summary get by name
// @Schemes
// @Param name path string true "Name"
// @Param max query int true "Max" default(0)
// @Description Get current index of `name`
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} utils.Unisponse.response
// @Router /get/{name} [get]
// @Router /g/{name} [get]
func Get(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")
	maxRaw := c.DefaultQuery("max", "0") // 0 - unlimited

	maxValue, err := strconv.ParseInt(maxRaw, 10, 64)
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
	if errors.Is(err, redis.Nil) {
		// key does not exist
		err = rdb.Set(ctx, name, 0, 0).Err()
		if err != nil {
			utils.Unirror(c, err)
			return
		}

		utils.Unisponse(c, 0)
		return
	}

	if err != nil {
		utils.Unirror(c, err)
		return
	}

	if maxValue > 0 && val >= maxValue {
		utils.Unisponse(c, maxValue)
		return
	}

	utils.Unisponse(c, val)
}
