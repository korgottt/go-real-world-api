package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/korgottt/go-real-world-api/model"
	"github.com/stretchr/testify/assert"
)

type StubArticlesStore struct {
	data []model.Article
}

func (s *StubArticlesStore) GetArticle(slug string) (article model.Article, e error) {
	e = fmt.Errorf("Article with slug %s was not found", slug)
	for _, a := range s.data {
		if a.Slug == slug {
			article = a
			e = nil
			break
		}
	}
	return
}

func (s *StubArticlesStore) CreateArticle(a model.SingleArticleWrap) (article model.Article, e error) {
	return a.Article, nil
}

func TestGETUsers(t *testing.T) {
	server := NewGlobalServer(&StubArticlesStore{})

	t.Run("returns user struct", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/user", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Body.String()
		want := defaultUserInfo

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns 200 on /api/user", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/api/user", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}

func TestPOSTuser(t *testing.T) {
	server := NewGlobalServer(&StubArticlesStore{})

	t.Run("Authentication test", func(t *testing.T) {
		requestBody := `{
			"user":{
			  "email": "jake@jake.jake",
			  "password": "jakejake"
			}
		  }`

		request, _ := http.NewRequest(http.MethodPost, "/api/users/login", strings.NewReader(requestBody))
		request.Header.Set("content-type", jsonContentType)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		var got model.SingleUserWrap
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from %q into User struct, '%v'", response.Body, err)
		}

		assert.Equal(t, "jake@jake.jake", got.User.Email)
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("Registration test", func(t *testing.T) {
		requestBody := `{
			"user":{
			  "username": "Jacob",
			  "email": "jake@jake.jake",
			  "password": "jakejake"
			}
		  }`

		request, _ := http.NewRequest(http.MethodPost, "/api/users", strings.NewReader(requestBody))
		request.Header.Set("content-type", jsonContentType)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		var got model.SingleUserWrap
		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse response from %q into User struct, '%v'", response.Body, err)
		}

		assert.Equal(t, "jakejake", got.User.Password)
		assertStatus(t, response.Code, http.StatusOK)
	})
}

func TestArticles(t *testing.T) {
	tt := []model.Article{
		{
			1,
			"slug-1",
			"title-1",
		},
		{
			2,
			"slug-2",
			"title-2",
		},
	}

	server := NewGlobalServer(&StubArticlesStore{tt})
	for _, tc := range tt {

		t.Run("returns articles struct", func(t *testing.T) {

			request := createGetArticleRequest(tc.Slug)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			assertStatus(t, response.Code, http.StatusOK)
			assertJSONBody(t, response.Body.String(), model.SingleArticleWrap{tc})
			assertContentType(t, response, jsonContentType)

		})

		t.Run("create new article", func(t *testing.T) {

			request := createPostArticleRequest(tc)
			response := httptest.NewRecorder()
			server.ServeHTTP(response, request)

			assertStatus(t, response.Code, http.StatusOK)
			assertJSONBody(t, response.Body.String(), tc)
			assertContentType(t, response, jsonContentType)

		})
	}
}

func createPostArticleRequest(a model.Article) *http.Request {
	serializedArticle, _ := json.Marshal(model.SingleArticleWrap{Article: a})
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/articles"), bytes.NewBuffer(serializedArticle))
	return req
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}

func createGetArticleRequest(slug string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/articles/%s", slug), nil)
	return req
}

func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func assertJSONBody(t *testing.T, got string, want interface{}) {
	t.Helper()
	builder := strings.Builder{}
	err := json.NewEncoder(&builder).Encode(want)
	if err != nil {
		t.Errorf("There is encode error:%s", err)
	}
	if !reflect.DeepEqual(got, builder.String()) {
		t.Errorf("response %v not equal %v", got, builder.String())
	}
}
