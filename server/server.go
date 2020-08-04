package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"

	"github.com/korgottt/go-real-world-api/model"
)

//var claims jwt.MapClaims

const jsonContentType = "application/json; charset=utf-8"

//ArticlesStore stores articles.
type ArticlesStore interface {
	GetArticle(slug string) (model.Article, error)
	CreateArticle(a model.SingleArticleWrap) (model.Article, error)
	RegUser(data model.User) (model.User, error)
	GetUser(username string) (model.User, error)
	UpdateUser(username string, data model.SingleUserWrap) (model.SingleUserWrap, error)
}

//GlobalServer general server struct
type GlobalServer struct {
	store ArticlesStore
	http.Handler
}

func isAuthorized(endpoint http.Handler) http.Handler {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(mySecretKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
	return jwtMiddleware.Handler(endpoint)
}

// NewGlobalServer return
func NewGlobalServer(a ArticlesStore) *GlobalServer {
	s := &GlobalServer{store: a}

	route := http.NewServeMux()
	route.Handle("/api/user", isAuthorized(http.HandlerFunc(s.getUserInfo)))
	route.Handle("/api/users/login", http.HandlerFunc(s.loginHandler))
	route.Handle("/api/users", http.HandlerFunc(s.regHandler))
	route.Handle("/api/articles/", http.HandlerFunc(s.getArticles))
	route.Handle("/api/articles", http.HandlerFunc(s.createArticles))

	s.Handler = route

	return s
}

func (gs *GlobalServer) getUserInfo(w http.ResponseWriter, r *http.Request) {
	t, err := jwtmiddleware.FromAuthHeader(r)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err.Error())
	}
	userData, _ := ParseToken(t)
	u, e := gs.store.GetUser(userData.User.UserName)
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		writeJSONResponse(w, model.SingleUserWrap{
			model.User{
				Email:    u.Email,
				Token:    t,
				UserName: u.UserName,
				Bio:      u.Bio,
				Image:    u.Image,
			},
		})
	}
}

func (gs *GlobalServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	var got model.SingleUserWrap
	err := json.NewDecoder(r.Body).Decode(&got)
	if err != nil {
		log.Fatalf("Can't decode request body, error: %v", err)
	}
	w.Header().Set("content-type", jsonContentType)
	got.User.Token = GenerateJwtToken(got.User.UserName, 0)
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
	_, err = gs.store.RegUser(got.User)
	if err != nil {
		fmt.Fprint(w, err.Error())
	}
	w.Header().Set("content-type", jsonContentType)
	got.User.Token = GenerateJwtToken(got.User.UserName, 0)
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

func writeJSONResponse(w http.ResponseWriter, v interface{}) {
	writeJSONContentType(w)
	json.NewEncoder(w).Encode(v)
}

func writeJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", jsonContentType)
}
