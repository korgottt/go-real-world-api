package model

type User struct {
	UserName string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type SingleUserWrap struct {
	User User `json:"user"`
}