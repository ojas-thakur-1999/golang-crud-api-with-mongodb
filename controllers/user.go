package controllers

import (
	"backend/mongo-crud-api/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{session: s}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	objectId := bson.ObjectIdHex(id)

	u := models.User{}

	if err := uc.session.DB("golang-backend-dev").C("users").FindId(objectId).One(&u); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsoned, err := json.Marshal(u)
	if err != nil {
		fmt.Println("failed in marshalling json output")
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsoned))
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user := models.User{}

	json.NewDecoder(r.Body).Decode(&user)

	user.Id = bson.NewObjectId()

	uc.session.DB("golang-backend-dev").C("users").Insert(user)

	jsoned, err := json.Marshal(user)
	if err != nil {
		fmt.Println("error in marshaling the user")
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(jsoned))
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	objectId := bson.ObjectIdHex(id)

	if err := uc.session.DB("golang-backend-dev").C("users").RemoveId(objectId); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, "Deleted user", objectId, "\n")
}
