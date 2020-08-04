package server

type Error struct{
	Body []string `json:"body"`
}

type SingleErrorWrap struct {
	Error Error `json:"errors"`
}