package route

import (
	svc "OutGetWay/context"
	"OutGetWay/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init(ctx *svc.Context) {
	service.InitService(ctx)
	engine := gin.Default()
	engine.GET("/GetHost", GetHost)
	if err := engine.Run(":" + ctx.Config.Port); err != nil {
		panic(err.Error())
	}
}

func GetHost(c *gin.Context) {
	c.String(http.StatusOK, service.SelectService())
	return
}
