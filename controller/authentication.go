package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/luisveasc/geocoder/middleware"
	"github.com/luisveasc/geocoder/model"
	"github.com/luisveasc/geocoder/util"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

// Roles en el sistema
const (
	RolAdmin = "Admin"
	RolUser  = "User"
)

// AuthenticationController : Estructura controladora de las colecciones
type AuthenticationController struct {
}

// Routes : Define las rutas del controlador
func (authenticationController *AuthenticationController) Routes(base *gin.RouterGroup, authNormal *jwt.GinJWTMiddleware) {

	// Refresh time can be longer than token timeout
	base.GET("/refresh_token", middleware.SetRoles(RolAdmin, RolUser), authNormal.MiddlewareFunc(), authNormal.RefreshHandler)

	//funcion de login: recibe un objeto {email: , pass:}
	base.POST("/login", authNormal.LoginHandler)

	//creacion de usuarios
	//base.POST("/users", CreateUser)
	base.POST("/users", middleware.SetRoles(RolAdmin, RolUser), authNormal.MiddlewareFunc(), CreateUser)
	//Obtener usuarios
	base.GET("/users", middleware.SetRoles(RolAdmin, RolUser), authNormal.MiddlewareFunc(), GetAllUsers)
	//creacion de usuarios
	base.DELETE("/users/:id", middleware.SetRoles(RolAdmin, RolUser), authNormal.MiddlewareFunc(), DeleteUser)
	//creacion de usuarios
	base.PUT("/users/:id", middleware.SetRoles(RolAdmin, RolUser), authNormal.MiddlewareFunc(), UpdateUser)
}

var userModel model.User

// UpdateUser : Obtener todos los usuarios de una institucion
func UpdateUser(c *gin.Context) {
	//userOwner := userModel.LoadFromContext(c)

	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		c.JSON(http.StatusNotFound, util.GetError("El id ingresado no es valido", nil))
	}
	userId := bson.ObjectIdHex(id)
	_, err := userModel.GetUser(userId)

	if err != nil {
		c.JSON(http.StatusNotFound, util.GetError("No se pudo encontrar al usuario", err))
		return
	}
	//if user.Institution != userOwner.Institution {
	//	c.JSON(http.StatusNotFound, util.GetError("El usuario no pertenece a la institucion", nil))
	//	return
	//}

	var userJson *model.User
	err = c.ShouldBind(&userJson)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo convertir json del usuario", err))
		return
	}

	if userJson.Password != "" {
		userJson.Hash = middleware.GeneratePassword(userJson.Password)
		userJson.Password = ""
	}
	err = userModel.UpdateUser(userId, userJson)

	if err != nil {
		c.JSON(http.StatusNotFound, util.GetError("No se pudo actualizar al usuario", err))
	}
	c.String(http.StatusOK, "")
}

// DeleteUser : Obtener todos los usuarios de una institucion
func DeleteUser(c *gin.Context) {
	//userOwner := userModel.LoadFromContext(c)
	id := c.Param("id")
	if !bson.IsObjectIdHex(id) {
		c.JSON(http.StatusNotFound, util.GetError("El id ingresado no es valido", nil))
	}
	userId := bson.ObjectIdHex(id)
	_, err := userModel.GetUser(userId)

	if err != nil {
		c.JSON(http.StatusNotFound, util.GetError("No se pudo encontrar al usuario", err))
		return
	}
	//if user.Institution != userOwner.Institution {
	//	c.JSON(http.StatusNotFound, util.GetError("El usuario no pertenece a la institucion", nil))
	//	return
	//}

	err = userModel.DeleteUser(userId)

	if err != nil {
		c.JSON(http.StatusNotFound, util.GetError("No se pudo eliminar al usuario", err))
		return
	}
	c.String(http.StatusOK, "")
}

// GetAllUsers : Obtener todos los usuarios de una institucion
func GetAllUsers(c *gin.Context) {
	//userOwner := userModel.LoadFromContext(c)

	/* obtener parametros de paginacion*/
	pagination := PaginationParams{}
	err := c.ShouldBind(&pagination)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se puedieron encontrar los parametros limit, offset", err))
		return
	}
	//page, err := userModel.FindPaginate(bson.M{"institution": userOwner.Institution}, pagination.Limit, pagination.Offset)
	page, err := userModel.FindPaginate(bson.M{}, pagination.Limit, pagination.Offset)

	if err != nil {
		c.JSON(http.StatusNotFound, util.GetError("No se pudo obtener la lista de institucion", err))
	}
	// c.Header("",page.Metadata.)
	if len(page.Metadata) != 0 {
		c.Header("Pagination-Count", fmt.Sprintf("%d", page.Metadata[0]["total"]))
	}

	c.JSON(http.StatusOK, page.Data)
}

// CreateUser : Crear usuario para la institucion
func CreateUser(c *gin.Context) {
	/* Crear usuario para institucion */
	userOwner := userModel.LoadFromContext(c)
	log.Printf("owner -> %+v", userOwner)
	//log.Printf("OBJECTID -> %+v", userOwner.Institution)
	var user model.User
	e := c.BindJSON(&user)
	if e != nil {
		c.JSON(http.StatusBadRequest, util.GetError("No se pudo registrar al usuario", e))
		return
	}
	user.ID = bson.NewObjectId()
	user.Hash = middleware.GeneratePassword(user.Password)
	user.Password = ""
	if user.Rol != RolUser && user.Rol != RolAdmin {
		//Error
		c.JSON(http.StatusInternalServerError, util.GetError("El rol ingresado no es correcto", nil))
		return
	}

	//user.Institution = userOwner.Institution
	if err := userModel.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, util.GetError("Fallo al crear el usuario", err))
		return
	}
	c.String(http.StatusCreated, "")

}

// // CreateUser : Registrar usuario
// func CreateUser(c *gin.Context) {
// 	var user model.User
// 	e := c.BindJSON(&user)
// 	if e != nil {
// 		c.JSON(http.StatusBadRequest, util.GetError("No se pudo registrar al usuario", e))
// 		return
// 	}
// 	user.ID = bson.NewObjectId()
// 	user.Hash = middleware.GeneratePassword(user.Password)
// 	user.Password = ""
// 	user.Rol = RolUser
// 	if err := userModel.Create(&user); err != nil {
// 		c.JSON(http.StatusInternalServerError, util.GetError("Fallo al crear el usuario", err))
// 		return
// 	}
// 	c.String(http.StatusCreated, "")
// }
