package model

import (
	"errors"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/luisveasc/geocoder/util"
	"github.com/mitchellh/mapstructure"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

// User : Usuario del sistema
type User struct {
	ID       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email    string        `json:"email" bson:"email,omitempty"`
	Name     string        `json:"name" bson:"name,omitempty"`
	Rol      string        `json:"rol" bson:"rol,omitempty"`
	Hash     string        `json:"_hash" bson:"_hash,omitempty"`
	Password string        `json:"password,omitempty" bson:"password,omitempty"`
}

// DeleteUser : Se elimina el usuario de la institucion
func (userModel *User) DeleteUser(id bson.ObjectId) error {

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()

	err := colUser.Remove(bson.M{"_id": id})

	return err
}

// UpdateUser : Se actualizan los campos del usuario de la institucion
func (userModel *User) UpdateUser(id bson.ObjectId, user *User) error {

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()
	err := colUser.Update(bson.M{"_id": id}, bson.M{"$set": user})
	return err
}

// GetUser : Se obtiene al usuario de la institucion
func (userModel *User) GetUser(id bson.ObjectId) (User, error) {

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()
	var user User
	err := colUser.Find(bson.M{"_id": id}).One(&user)

	return user, err
}

// FindPaginate : Obtener usuario
func (userModel *User) FindPaginate(query bson.M, limit int, offset int) (Page, error) {

	col, session := GetCollection(CollectionNameUser)
	defer session.Close()
	pag := []bson.M{
		{"$project": bson.M{"_hash": 0}},
		{"$skip": offset},
	}
	if limit > 0 {
		pag = append(pag, bson.M{"$limit": limit})
	}
	pipeline := []bson.M{
		bson.M{"$match": query},
		bson.M{"$facet": bson.M{
			"metadata": []bson.M{{"$count": "total"}},
			"data":     pag, // add projection here wish you re-shape the docs
		}},
	}

	pageDoc := Page{}
	err := col.Pipe(pipeline).One(&pageDoc)

	return pageDoc, err
}

// LoadFromContext : Traer usuario desde contexto
func (userModel *User) LoadFromContext(c *gin.Context) *User {
	claims := jwt.ExtractClaims(c)
	var user User
	err := mapstructure.Decode(claims["user"], &user)
	if err != nil {
		panic(err)
	}
	user.ID = bson.ObjectIdHex(claims["user"].(map[string]interface{})["id"].(string))
	//user.Institution = bson.ObjectIdHex(claims["user"].(map[string]interface{})["institution"].(string))
	user.Hash = ""
	return &user
}

// Create : Traer usuario desde contexto
func (userModel *User) Create(user *User) error {

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()
	/* Verificar que el usuario no exista*/
	var userFound User
	err := colUser.Find(bson.M{"email": user.Email}).One(&userFound)
	if err == nil {
		return errors.New("El email ingresado ya existe")
	}

	err = colUser.Insert(&user)

	return err
}

// GetUser : Se obtiene el usuario
func GetUser(c *gin.Context) {
	id := c.Param("id")

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()
	var usuario User

	if err := colUser.FindId(bson.ObjectIdHex(id)).One(&usuario); err != nil {
		c.JSON(http.StatusNotFound, util.GetError("Usuario no encontrado", err))
	} else {
		c.JSON(http.StatusCreated, usuario)
	}
}

// CreateUsersBulk : Registrar usuario
func CreateUsersBulk(c *gin.Context) {
	var users []User
	e := c.BindJSON(&users)
	util.Check(e)

	colUser, session := GetCollection(CollectionNameUser)
	defer session.Close()
	type Par struct {
		Usuario User
		Result  bson.M
	}

	type Respuesta struct {
		NoCreados []Par
		Creados   []Par
	}

	var resp Respuesta
	resp.Creados = make([]Par, 0)
	resp.NoCreados = make([]Par, 0)

	for _, u := range users {
		if u.Email == "" {
			var aux Par
			aux.Usuario = u
			aux.Result = bson.M{"mensaje": "No se especifico un email."}
			resp.NoCreados = append(resp.NoCreados, aux)
			continue
		}
		var temp User
		log.Printf("Buscando %s\n", u.Email)
		if err := colUser.Find(bson.M{"email": u.Email}).One(&temp); err != nil {
			//no existe, por lo que puedo crearlo

			u.ID = bson.NewObjectId()
			if err := colUser.Insert(&u); err != nil {
				var aux Par
				aux.Usuario = u
				aux.Result = bson.M{"mensaje": err.Error()}
				resp.NoCreados = append(resp.NoCreados, aux)
			} else {
				var aux Par
				aux.Usuario = u
				aux.Result = bson.M{"mensaje": "Usuario creado OK."}
				resp.Creados = append(resp.Creados, aux)
			}
		} else {
			var aux Par
			aux.Usuario = u
			aux.Result = bson.M{"mensaje": "Ese email ya esta siendo usado."}
			resp.NoCreados = append(resp.NoCreados, aux)
		}
	}
	c.JSON(http.StatusOK, resp)
}
