package server

import (
	"realworldapi/model"
)

//InMemoryStore collects data about articles in memory.
type InMemoryStore struct {
	Article []model.Article
}

//GetArticle return single article if exists
func (s *InMemoryStore) GetArticle(slug string) (article model.Article, e error) {
	return model.Article{Slug: "some-slug", Title: "some title"}, nil
}

//CreateArticle return article by http request model
func (s *InMemoryStore) CreateArticle(a model.SingleArticleWrap) (article model.Article, e error) {
	return a.Article, nil
}
