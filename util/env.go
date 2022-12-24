package util

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func checkVars() []string {
	vars := []string{"GEOC_GO_REST_ENV",
		"GEOC_GIN_MODE",
		"GEOC_DB_USER",
		"GEOC_DB_PASS",
		"GEOC_DB_NAME",
		"GEOC_DB_URL",
		"GEOC_NOMINATIM_API",
		"GEOC_LOG_FILENAME",
		"GEOC_LOG_NUM_BACKUPS",
		"GEOC_LOG_AMOUNT_MB",
		"GEOC_NOMINATIM_API",
		"GEOC_ADDR",
		"GEOC_JWT_KEY"}
	missing := []string{}
	for _, v := range vars {
		_, set := os.LookupEnv(v)
		if !set {
			missing = append(missing, v)
		}
	}
	return missing
}

// LoadEnv : Se cargan variables de entorno
func LoadEnv() {
	env := os.Getenv("GEOC_GO_REST_ENV")
	if "" == env {
		env = "dev"
	}
	godotenv.Load(".env." + env)

	if vars := checkVars(); len(vars) != 0 {
		log.Printf("ERROR: Variables de entorno necesarias no definidas: %v", vars)
		panic(fmt.Sprintf("ERROR: Variables de entorno necesarias no definidas: %v", vars))
	}
}
