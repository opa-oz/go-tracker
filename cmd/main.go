package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "tracker/cmd/docs"
	"tracker/pkg/api"
	"tracker/pkg/config"
	"tracker/pkg/middlewares"
	"tracker/pkg/redis"
)

func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	if cfg.Production {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	rdb := redis.GetClient(cfg)
	docs.SwaggerInfo.BasePath = "/api"

	r.Use(middlewares.RequestLogger())
	r.Use(middlewares.ResponseLogger())

	r.GET("/healz", api.Healz)
	r.GET("/ready", api.Ready)
	// {@link https://github.com/swaggo/gin-swagger?tab=readme-ov-file}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	rg := r.Group("/api")
	rg.Use(middlewares.CfgMiddleware(cfg))
	rg.Use(middlewares.RedisMiddleware(rdb))

	{
		rg.GET("/i/:name", api.Increase)
		rg.GET("/increase/:name", api.Increase)

		rg.GET("/d/:name", api.Decrease)
		rg.GET("/decrease/:name", api.Decrease)

		rg.GET("/g/:name", api.Get)
		rg.GET("/get/:name", api.Get)

		rg.GET("/r/:name", api.Reset)
		rg.GET("/reset/:name", api.Reset)
	}

	port := fmt.Sprintf(":%d", cfg.Port)
	err = r.Run(port)
	if err != nil {
		fmt.Println(err)
		return
	}
}
