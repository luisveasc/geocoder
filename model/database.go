package model

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/globalsign/mgo"
)

const (
	CollectionNameUser = "users"
)

var (
	session *mgo.Session
)

// Page : Pagina de resultado
type Page struct {
	Metadata []map[string]int `json:"metadata" bson:"metadata,omitempty"`
	Data     []interface{}    `json:"data" bson:"data,omitempty"`
}

// GetCollection : Obtener coleccion desde la bd
func GetCollection(collection string) (*mgo.Collection, *mgo.Session) {
	s := session.Copy()
	return s.DB(os.Getenv("GEOC_DB_NAME")).C(collection), s
}

// LoadDatabase : Carga la base de datos y devuelve la session correspondiente
func LoadDatabase() {
	info := &mgo.DialInfo{
		Addrs:    []string{os.Getenv("GEOC_DB_URL")},
		Timeout:  30 * time.Second,
		Database: os.Getenv("GEOC_DB_NAME"),
		Username: os.Getenv("GEOC_DB_USER"),
		Password: os.Getenv("GEOC_DB_PASS"),
	}
	var err error
	session, err = mgo.DialWithInfo(info)
	if err != nil {
		log.Printf("Conection to DB,  %s/%s, (userName:%s) Error: %s", os.Getenv("GEOC_DB_URL"), os.Getenv("GEOC_DB_NAME"), os.Getenv("GEOC_DB_USER"), err)
		panic(fmt.Sprintf("Conection to DB,  %s/%s, (userName:%s) Error: %s", os.Getenv("GEOC_DB_URL"), os.Getenv("GEOC_DB_NAME"), os.Getenv("GEOC_DB_USER"), err))
	}
	log.Printf("Conected to DB, %s/%s", os.Getenv("GEOC_DB_URL"), os.Getenv("GEOC_DB_NAME"))

	session.SetMode(mgo.Monotonic, true)
}
