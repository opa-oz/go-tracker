package api

import (
	"strconv"
	"tracker/pkg/utils"

	"github.com/gin-gonic/gin"
)

// @BasePath /api

// Reset godoc
// @Summary reset by name
// @Schemes
// @Param name path string true "Name"
// @Param init query int true "Init" default(0)
// @Description Reset current index of `name`
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} utils.Unisponse.response
// @Router /reset/{name} [get]
// @Router /r/{name} [get]
func Reset(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")
	init := c.DefaultQuery("init", "0") // 0 - unlimited

	initValue, err := strconv.ParseInt(init, 10, 64)
	if err != nil {
		utils.Unirror(c, err)
		return
	}

	rdb, err := utils.GetRedis(c)
	if err != nil {
		utils.Unirror(c, err)
		return
	}

	err = rdb.Set(ctx, name, initValue, 0).Err()
	if err != nil {
		utils.Unirror(c, err)
		return
	}

	utils.Unisponse(c, initValue)
}
