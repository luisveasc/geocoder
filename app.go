package main

import (
	"os"
	"time"

	"github.com/luisveasc/geocoder/util"
)

func main() {

	time.Local = time.UTC

	util.LoadEnv()

	util.LoadLogFile(os.Getenv("GEOC_LOG_FILENAME"), os.Getenv("GEOC_LOG_NUM_BACKUPS"), os.Getenv("GEOC_LOG_AMOUNT_MB"))

	// Cargar Variables de entorno:
	//

	// // Log
	// log.Println("Start aysana-backend-auth")
	// log.Printf("serverUp, %s", os.Getenv("ADDR"))

	// // Cargar base de datos
	// model.LoadDatabase()

	// gin.SetMode(os.Getenv("ABAU_GIN_MODE"))
	// //Raiz
	// app := gin.Default()

	// //LOGS
	// gin.DefaultWriter = io.MultiWriter(log.Writer())
	// app.Use(gin.Logger())

	// // CORS
	// app.Use(middleware.CorsMiddleware())
	// // Url Base
	// base := app.Group("/api/v1/")

	// controller.Routes(base)

	// app.NoRoute(func(c *gin.Context) {
	// 	c.JSON(404, gin.H{"message": "Servicio no encontrado."})
	// })

	// http.ListenAndServe(os.Getenv("ABAU_ADDR"), app)

}
