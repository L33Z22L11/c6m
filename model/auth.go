package model

type Auth struct {
	Uid      string `json:"uid"`
	Username string `json:"username"`
	Hash     string `json:"hash"`
	Salt     string `json:"salt"`
}
