package main

import (
	"encoding/json"
	"fmt"
	"log"
	// "log"
	"net/http"
	"os"
	// import mux
	"github.com/gorilla/mux"
	//import mongodb  lib
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// import app engine
)

// var (
// 	db_uri, db, collection string
// )

// func init() {
// 	db_uri := "mongodb://duynguyen0428:cuongduy0428@ds221228.mlab.com:21228/todoapp"
// 	db := "todoapp"
// 	collection := "user"
// }

type User struct {
	Email    string `json:"email`
	Password string `json:"password"`
}

func main() {
	// router := AppRouter()
	router := mux.NewRouter()
	// Assign router-handler
	router.HandleFunc("/user", GetUsersHander).Methods("GET")
	router.HandleFunc("/user", CreateUserHander).Methods("POST")

	// http.Handle("/", router)
	// appengine.Main()
	log.Fatal(http.ListenAndServe(":8000", router))
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

func CreateUserHander(w http.ResponseWriter, r *http.Request) {
	dburi := "mongodb://duynguyen0428:cuongduy0428@ds221228.mlab.com:21228/todoapp"
	// db := "todoapp"
	// collection := "user"

	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		// Return Bad Request Here
		// http.Error(w, err.Error(), http.StatusBadRequest)
		ErrorWithJSON(w, "error", http.StatusBadRequest)
		return
	}

	sess, err := mgo.Dial(dburi)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	coll := sess.DB("todoapp").C("user")

	err = coll.Insert(&user)
	if err != nil {
		// Return Bad Request Here
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		ErrorWithJSON(w, "error", http.StatusInternalServerError)
		return
	}

	ResponseWithJSON(w, nil, http.StatusCreated)

}

func GetUsersHander(w http.ResponseWriter, r *http.Request) {

	dburi := "mongodb://duynguyen0428:cuongduy0428@ds221228.mlab.com:21228/todoapp"
	// db := "todoapp"
	// collection := "user"
	var users []User

	sess, err := mgo.Dial(dburi)
	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})

	coll := sess.DB("todoapp").C("user")

	err = coll.Find(bson.M{}).All(&users)
	if err != nil {
		// Return Bad Request Here
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		ErrorWithJSON(w, "error", http.StatusInternalServerError)
		return
	}

	usersDTO, err := json.Marshal(users)

	ResponseWithJSON(w, usersDTO, http.StatusOK)

}
