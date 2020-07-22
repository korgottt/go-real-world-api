package model

//User represent table model of users
type User struct {
	ID       int    `json:"id" db:"id"`
	UserName string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Token    string `json:"token" db:"token"`
	Password string `json:"password" db:"password"`
	Bio      string `json:"bio" db:"bio"`
}

//SingleUserWrap is http request/response model for single user
type SingleUserWrap struct {
	User User `json:"user"`
}

//Article represent table model of articles
type Article struct {
	ID    int    `json:"id" db:"id"`
	Slug  string `json:"slug" db:"slug"`
	Title string `json:"title" db:"title"`
}

//SingleArticleWrap is http request/response model for single user
type SingleArticleWrap struct {
	Article Article `json:"article"`
}
