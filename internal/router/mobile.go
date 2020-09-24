package router

import (
	"github.com/gin-gonic/gin"
	"github.com/realwrtoff/go-http/internal/service"
)

func InitMobileRouter(Router *gin.Engine, svc *service.Service) {
	MobileRouter := Router.Group("mobile")
	{
		MobileRouter.GET("set", svc.SetMobileMd5)
		MobileRouter.GET("query", svc.QueryMobile)
		MobileRouter.POST("query", svc.QueryMobileMulti)
	}
}