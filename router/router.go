package router

import (
	"github.com/gin-gonic/gin"
	"jrsfun-server-go/controller/api"
)

func StartRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Any("//api/home", api.Home)
	router.Any("//api/live", api.Live)
	router.Any("//api/channel", api.Channel)

	return router
}
