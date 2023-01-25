package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/luisveasc/geocoder/middleware"
)

// PaginationParams : Parametros de paginacion
type PaginationParams struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

// Controllers
var authenticationController AuthenticationController
var proxyController ProxyController

//Models
//debemos ir agregando ejemplo SERVICES
//var institutionModel model.Institution

func Routes(base *gin.RouterGroup) {
	// Middleware
	authNormal := middleware.LoadJWTAuth(middleware.LoginFunc)

	authenticationController.Routes(base, authNormal)
	proxyController.Routes(base, authNormal)

}
