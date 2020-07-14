package model

//User represent table model of users
type User struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

//SingleUserWrap is http request/response model for single user
type SingleUserWrap struct {
	User User `json:"user"`
}

//Article represent table model of articles
type Article struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

//SingleArticleWrap is http request/response model for single user
type SingleArticleWrap struct {
	Article Article `json:"article"`
}
