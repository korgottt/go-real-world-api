package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"realworldapi/model"
)

var defaultUserInfo = `{
	"user": {
	  "email": "jake@jake.jake",
	  "token": "jwt.token.here",
	  "username": "jake",
	  "bio": "I work at statefarm",
	  "image": null
	}
  }`

// type PlayerStore interface {
// 	GetPlayerInfo() string
// }

//GlobalServer general server struct
type GlobalServer struct {
	//store PlayerStore
	http.Handler
}

func NewGlobalServer() *GlobalServer {
	s := &GlobalServer{}

	route := http.NewServeMux()
	route.Handle("/api/user", http.HandlerFunc(s.getUserInfo))
	route.Handle("/api/users/login", http.HandlerFunc(s.loginHandler))
	route.Handle("/api/users", http.HandlerFunc(s.regHandler))

	s.Handler = route

	return s
}

func (gs *GlobalServer) getUserInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, defaultUserInfo)
}

func (gs *GlobalServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	var got model.SingleUserWrap
	err := json.NewDecoder(r.Body).Decode(&got)
	if err != nil {
		log.Fatalf("Can't decode request body, error: %v", err)
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(got)
	if err != nil {
		log.Fatalf("Failed encode error: %v", err)
	}
}

func (gs *GlobalServer) regHandler(w http.ResponseWriter, r *http.Request) {
	var got model.SingleUserWrap
	err := json.NewDecoder(r.Body).Decode(&got)
	if err != nil {
		log.Fatalf("Can't decode request body, error: %v", err)
	}
	w.Header().Set("content-type", "application/json; charset=utf-8")
	err = json.NewEncoder(w).Encode(got)
	if err != nil {
		log.Fatalf("Failed encode error: %v", err)
	}
}
