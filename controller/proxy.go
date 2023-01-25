package controller

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/luisveasc/geocoder/util"
)

// ProxyController : Controlador del proxy
type ProxyController struct {
}

// Routes : Define las rutas del controlador
func (proxyController *ProxyController) Routes(base *gin.RouterGroup, authNormal *jwt.GinJWTMiddleware) *gin.RouterGroup {

	// Rutas
	// middleware.SetRoles(RolAdmin, RolUser)
	// authNormal.MiddlewareFunc())
	proxyRouter := base.Group("", authNormal.MiddlewareFunc())
	{
		//debemos ir separando por grupos de funcionalidades (servicios ofrecidos) para que este más ordenado,
		//ejemplo: "search" => búsquedas generalizadas
		proxyRouter.Any("/search/*d", proxyController.Any(os.Getenv("GEOC_NOMINATIM_API")))
		//proxyRouter.Any("/instancesapps/*d", proxyController.Any(os.Getenv("ABAU_APP_MANAGER")))

	}
	return proxyRouter
}

// GetAll : Proxy a servicios
func (dogController *ProxyController) Any(route string) func(c *gin.Context) {
	return func(c *gin.Context) {
		user := userModel.LoadFromContext(c)
		client := &http.Client{}
		//uri := route + c.Request.RequestURI
		uri := route + "/search.php"
		resp := "" + route + "" + uri

		r, err := http.NewRequest(c.Request.Method, uri, c.Request.Body)
		paramPairs := c.Request.URL.Query()

		if paramPairs != nil {
			r.URL.RawQuery = paramPairs.Encode()
		}

		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo crear consulta", err))
			return
		}
		r.Header.Set("user", user.ID.Hex())
		//r.Header.Set("institution", user.Institution.Hex())

		response, err := client.Do(r)
		if err != nil {
			c.JSON(http.StatusBadRequest, util.GetError("No se pudo crear consulta", err))
			return
		}
		if response != nil && response.Body != nil {
			reader := response.Body
			contentLength := response.ContentLength
			contentType := response.Header.Get("Content-Type")

			extraHeaders := map[string]string{

				//"Content-Disposition": `attachment; filename="gopher.png"`,
			}
			for name, value := range response.Header {
				if len(value) > 0 {
					extraHeaders[name] = value[0]
				}
			}
			log.Printf("extra: %+v", r.Header)

			c.DataFromReader(response.StatusCode, contentLength, contentType, reader, extraHeaders)

			defer func() {
				io.Copy(ioutil.Discard, response.Body)
				response.Body.Close()
			}()
		} else {
			c.String(http.StatusOK, resp)
		}
	}
}
