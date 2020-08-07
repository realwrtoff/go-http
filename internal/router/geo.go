package router

import (
	"github.com/gin-gonic/gin"
	"github.com/realwrtoff/go-http/internal/service"
)

func InitGeoRouter(Router *gin.Engine, svc *service.Service) {
	GeoRouter := Router.Group("geo")
	{
		GeoRouter.GET("address", svc.Address)
		GeoRouter.GET("distance", svc.Distance)
	}
}
