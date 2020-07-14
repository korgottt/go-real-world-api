package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"realworldapi/model"
	"strings"
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

const jsonContentType = "application/json; charset=utf-8"

//ArticlesStore stores articles.
type ArticlesStore interface {
	GetArticle(slug string) (model.Article, error)
	CreateArticle(a model.SingleArticleWrap) (model.Article, error)
}

//GlobalServer general server struct
type GlobalServer struct {
	store ArticlesStore
	http.Handler
}

// NewGlobalServer return
func NewGlobalServer(a ArticlesStore) *GlobalServer {
	s := &GlobalServer{store: a}

	route := http.NewServeMux()
	route.Handle("/api/user", http.HandlerFunc(s.getUserInfo))
	route.Handle("/api/users/login", http.HandlerFunc(s.loginHandler))
	route.Handle("/api/users", http.HandlerFunc(s.regHandler))
	route.Handle("/api/articles/", http.HandlerFunc(s.getArticles))
	route.Handle("/api/articles", http.HandlerFunc(s.createArticles))

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
	w.Header().Set("content-type", jsonContentType)
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
	w.Header().Set("content-type", jsonContentType)
	err = json.NewEncoder(w).Encode(got)
	if err != nil {
		log.Fatalf("Failed encode error: %v", err)
	}
}

func (gs *GlobalServer) getArticles(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/api/articles/")
	article, err := gs.store.GetArticle(slug)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("content-type", jsonContentType)
		json.NewEncoder(w).Encode(model.SingleArticleWrap{Article: article})
	}
}

func (gs *GlobalServer) createArticles(w http.ResponseWriter, r *http.Request) {
	var got model.SingleArticleWrap
	body, _ := ioutil.ReadAll(r.Body)
	_ = json.NewDecoder(bytes.NewReader(body)).Decode(&got)
	article, err := gs.store.CreateArticle(got)
	if err != nil {
		log.Fatalf("Failed ctreation an article: %v", err)
	}
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(article)
}
