package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/luisveasc/geocoder/controller"
	"github.com/luisveasc/geocoder/middleware"
	"github.com/luisveasc/geocoder/model"
	"github.com/luisveasc/geocoder/util"
)

func main() {

	time.Local = time.UTC

	// Cargar Variables de entorno:
	util.LoadEnv()

	//inicializa log
	util.LoadLogFile(os.Getenv("GEOC_LOG_FILENAME"), os.Getenv("GEOC_LOG_AMOUNT_MB"), os.Getenv("GEOC_LOG_NUM_BACKUPS"))

	log.Printf("Start system: %s, version: V%s, created by: %s", os.Getenv("GEOC_SYSTEM_NAME"), os.Getenv("GEOC_VERSION"), os.Getenv("GEOC_ORGANIZATION"))
	log.Printf("serverUp, %s", os.Getenv("GEOC_ADDR"))

	// Cargar base de datos
	model.LoadDatabase()

	gin.SetMode(os.Getenv("GEOC_GIN_MODE"))
	//Raiz
	app := gin.Default()

	//LOGS
	gin.DefaultWriter = io.MultiWriter(log.Writer())
	app.Use(gin.Logger())

	// CORS
	app.Use(middleware.CorsMiddleware())
	// Url Base
	base := app.Group(os.Getenv("GEOC_URLBASE"))
	controller.Routes(base)

	app.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Servicio no encontrado."})
	})

	http.ListenAndServe(os.Getenv("GEOC_ADDR"), app)

}
